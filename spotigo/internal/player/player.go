// Package player handles Spotify auth and streams audio to PCM files.
package player

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode"

	librespot "github.com/devgianlu/go-librespot"
	"github.com/devgianlu/go-librespot/ap"
	lsplayer "github.com/devgianlu/go-librespot/player"
	devicespb "github.com/devgianlu/go-librespot/proto/spotify/connectstate/devices"
	metadatapb "github.com/devgianlu/go-librespot/proto/spotify/metadata"
	"github.com/devgianlu/go-librespot/session"
	"github.com/devgianlu/go-librespot/spclient"

	"github.com/jmanser/spotigo/internal/creds"
)

const (
	bitrate   = 320
	oauthPort = 5173
)

// Session holds an authenticated Spotify session.
type Session struct {
	sess        *session.Session
	countryCode *string
}

// New creates an authenticated Session. On first run (no stored credentials),
// it triggers an interactive OAuth2 login at http://localhost:{oauthPort}/login.
func New(c *creds.Creds, credsFile string) (*Session, error) {
	ctx := context.Background()

	blob, err := c.StoredBytes()
	if err != nil {
		return nil, err
	}

	var authCreds any
	interactive := false
	if len(blob) > 0 {
		authCreds = session.StoredCredentials{Username: c.Username, Data: blob}
	} else {
		fmt.Printf("No credentials — open http://localhost:%d/login in your browser.\n", oauthPort)
		authCreds = session.InteractiveCredentials{CallbackPort: oauthPort}
		interactive = true
	}

	sess, err := session.NewSessionFromOptions(ctx, &session.Options{
		Log:         newLogger(),
		DeviceType:  devicespb.DeviceType_COMPUTER,
		DeviceId:    randomHex(20),
		Credentials: authCreds,
	})
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	if interactive {
		c.Username = sess.Username()
		c.SetStored(sess.StoredCredentials())
		if err := creds.Save(credsFile, c); err != nil {
			fmt.Fprintf(os.Stderr, "warning: save credentials: %v\n", err)
		} else {
			fmt.Println("Credentials saved to", credsFile)
		}
	}

	cc := new(string)
	apRecv := sess.Accesspoint().Receive(ap.PacketTypeCountryCode)
	go func() {
		for pkt := range apRecv {
			if pkt.Type == ap.PacketTypeCountryCode {
				*cc = string(pkt.Payload)
			}
		}
	}()

	return &Session{sess: sess, countryCode: cc}, nil
}

// Run streams the Spotify URL or URI, writing f32le PCM files into outDir.
// Albums and playlists create a subdirectory named after the album/playlist (or
// albumOverride if non-empty). Tracks write directly into outDir.
func (s *Session) Run(ctx context.Context, urlOrURI, outDir, albumOverride string) error {
	uri := toURI(urlOrURI)

	parts := strings.SplitN(uri, ":", 3)
	if len(parts) != 3 || parts[0] != "spotify" {
		return fmt.Errorf("unrecognised Spotify URL or URI: %s", urlOrURI)
	}
	typ, id62 := parts[1], parts[2]

	switch typ {
	case "track", "episode":
		id, err := librespot.SpotifyIdFromBase62(librespot.SpotifyIdType(typ), id62)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(outDir, 0o750); err != nil {
			return err
		}
		return s.streamTrack(ctx, *id, outDir, 0)

	default: // album, playlist, show, …
		return s.streamContext(ctx, uri, outDir, albumOverride)
	}
}

// streamContext resolves an album/playlist and streams every track.
func (s *Session) streamContext(ctx context.Context, uri, outDir, albumOverride string) error {
	spotCtx, err := s.sess.Spclient().ContextResolve(ctx, uri)
	if err != nil {
		return fmt.Errorf("resolve context: %w", err)
	}

	resolver, err := spclient.NewContextResolver(ctx, newLogger(), s.sess.Spclient(), spotCtx)
	if err != nil {
		return fmt.Errorf("context resolver: %w", err)
	}

	var trackURIs []string
	for i := 0; ; i++ {
		tracks, err := resolver.Page(ctx, i)
		if err != nil {
			break
		}
		for _, t := range tracks {
			if u := t.GetUri(); u != "" {
				trackURIs = append(trackURIs, u)
			} else if gid := t.GetGid(); len(gid) > 0 {
				id := librespot.SpotifyIdFromGid(librespot.SpotifyIdTypeTrack, gid)
				trackURIs = append(trackURIs, id.Uri())
			}
		}
	}

	if len(trackURIs) == 0 {
		return fmt.Errorf("no tracks found for %s", uri)
	}

	dirName := contextDirName(resolver)
	if albumOverride != "" {
		dirName = safeFilename(albumOverride)
	}
	dir := filepath.Join(outDir, dirName)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	fmt.Printf("Found %d tracks → %s\n", len(trackURIs), dir)

	for i, trackURI := range trackURIs {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		p := strings.SplitN(trackURI, ":", 3)
		if len(p) != 3 {
			continue
		}
		id, err := librespot.SpotifyIdFromBase62(librespot.SpotifyIdType(p[1]), p[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", trackURI, err)
			continue
		}
		if err := s.streamTrackWithRetry(ctx, *id, dir, i+1); err != nil {
			fmt.Fprintf(os.Stderr, "error: track %d: %v\n", i+1, err)
		}
	}

	return nil
}

func (s *Session) streamTrackWithRetry(ctx context.Context, id librespot.SpotifyId, dir string, pos int) error {
	delays := []time.Duration{20 * time.Second, 40 * time.Second, 60 * time.Second}
	err := s.streamTrack(ctx, id, dir, pos)
	for i, d := range delays {
		if err == nil || ctx.Err() != nil {
			break
		}
		if !strings.Contains(err.Error(), "aes key") {
			break
		}
		fmt.Fprintf(os.Stderr, "audio key rate-limited, retrying track %d in %s (attempt %d/4)\n", pos, d, i+2)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(d):
		}
		err = s.streamTrack(ctx, id, dir, pos)
	}
	return err
}

// streamTrack streams one track to a PCM file inside dir.
// pos > 0 prefixes the filename with a zero-padded track number.
func (s *Session) streamTrack(ctx context.Context, id librespot.SpotifyId, dir string, pos int) error {
	fifoPath := filepath.Join(dir, ".fifo-"+randomHex(8))
	if err := syscall.Mkfifo(fifoPath, 0o600); err != nil {
		return fmt.Errorf("create fifo: %w", err)
	}
	defer os.Remove(fifoPath)

	// Open read end non-blocking so the player's O_WRONLY|O_NONBLOCK open succeeds.
	rfd, err := syscall.Open(fifoPath, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return fmt.Errorf("open fifo: %w", err)
	}
	if err := syscall.SetNonblock(rfd, false); err != nil {
		_ = syscall.Close(rfd)
		return fmt.Errorf("set blocking: %w", err)
	}
	reader := os.NewFile(uintptr(rfd), fifoPath)

	pl, err := lsplayer.NewPlayer(&lsplayer.Options{
		Spclient:              s.sess.Spclient(),
		AudioKey:              s.sess.AudioKey(),
		Events:                s.sess.Events(),
		Log:                   newLogger(),
		CountryCode:           s.countryCode,
		NormalisationEnabled:  true,
		AudioBackend:          "pipe",
		AudioOutputPipe:       fifoPath,
		AudioOutputPipeFormat: "f32le",
	})
	if err != nil {
		_ = reader.Close()
		return fmt.Errorf("create player: %w", err)
	}

	events := pl.Receive()

	stream, err := pl.NewStream(ctx, http.DefaultClient, id, bitrate, 0)
	if err != nil {
		_ = reader.Close()
		pl.Close()
		return fmt.Errorf("load stream: %w", err)
	}

	s.downloadCover(ctx, dir, stream)
	finalPath := filepath.Join(dir, trackFilename(stream, pos))
	tmpWAV := finalPath + ".tmp"

	copyDone := make(chan error, 1)
	go func() {
		f, err := os.Create(tmpWAV)
		if err != nil {
			_ = reader.Close()
			copyDone <- err
			return
		}
		err = writeWAV(f, reader)
		_ = f.Close()
		_ = reader.Close()
		if err != nil {
			_ = os.Remove(tmpWAV)
		}
		copyDone <- err
	}()

	if err := pl.SetPrimaryStream(stream.Source, false, false); err != nil {
		pl.Close()
		<-copyDone
		_ = os.Remove(tmpWAV)
		return fmt.Errorf("set stream: %w", err)
	}
	if err := pl.Play(); err != nil {
		pl.Close()
		<-copyDone
		_ = os.Remove(tmpWAV)
		return fmt.Errorf("play: %w", err)
	}

	fmt.Printf("  → %s\n", filepath.Base(finalPath))

	for {
		select {
		case <-ctx.Done():
			pl.Close()
			<-copyDone
			_ = os.Remove(tmpWAV)
			return ctx.Err()
		case ev, ok := <-events:
			if !ok || ev.Type == lsplayer.EventTypeStop || ev.Type == lsplayer.EventTypeNotPlaying {
				pl.Close()
				if err := <-copyDone; err != nil {
					_ = os.Remove(tmpWAV)
					return err
				}
				return os.Rename(tmpWAV, finalPath)
			}
		}
	}
}

// Close releases the session.
func (s *Session) Close() {
	s.sess.Close()
}

// contextDirName derives a safe directory name from the resolver's metadata or URI.
func contextDirName(r *spclient.ContextResolver) string {
	for _, k := range []string{"context_description", "name", "playlist.title", "title"} {
		if v := r.Metadata()[k]; v != "" {
			return safeFilename(v)
		}
	}
	parts := strings.Split(r.Uri(), ":")
	return safeFilename(parts[len(parts)-1])
}

// trackFilename builds "01 - Artist - Title.pcm" (or "Artist - Title.pcm" when pos==0).
func trackFilename(stream *lsplayer.Stream, pos int) string {
	title := stream.Media.Name()
	artist := ""
	if t := stream.Media.Track(); t != nil && len(t.GetArtist()) > 0 {
		artist = t.GetArtist()[0].GetName()
	}
	name := title
	if artist != "" {
		name = artist + " - " + title
	}
	if pos > 0 {
		name = fmt.Sprintf("%02d - %s", pos, name)
	}
	return safeFilename(name) + ".wav"
}

// downloadCover saves cover.jpg in dir from the track's album art (no-op if already exists).
func (s *Session) downloadCover(ctx context.Context, dir string, stream *lsplayer.Stream) {
	if !stream.Media.IsTrack() {
		return
	}
	coverPath := filepath.Join(dir, "cover.jpg")
	if _, err := os.Stat(coverPath); err == nil {
		return
	}

	album := stream.Media.Track().GetAlbum()
	if album == nil {
		return
	}
	var images []*metadatapb.Image
	if cg := album.GetCoverGroup(); cg != nil {
		images = cg.GetImage()
	} else {
		images = album.GetCover()
	}
	if len(images) == 0 {
		return
	}

	img := images[0]
	for _, i := range images {
		if i.GetSize() > img.GetSize() {
			img = i
		}
	}
	url := "https://i.scdn.co/image/" + hex.EncodeToString(img.GetFileId())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(coverPath)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = io.Copy(f, resp.Body)
}

// writeWAV writes a WAV file (IEEE float 32-bit, 44100 Hz, stereo) from r into w.
func writeWAV(w *os.File, r io.Reader) error {
	const (
		sampleRate    = 44100
		channels      = 2
		bitsPerSample = 32
		byteRate      = sampleRate * channels * bitsPerSample / 8
		blockAlign    = channels * bitsPerSample / 8
		hdrSize       = 44
	)

	// Reserve space for the header; fill it in after we know the data size.
	if _, err := w.Write(make([]byte, hdrSize)); err != nil {
		return err
	}

	n, err := io.Copy(w, r)
	if err != nil {
		return err
	}

	if _, err := w.Seek(0, io.SeekStart); err != nil {
		return err
	}

	var buf [hdrSize]byte
	copy(buf[0:], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:], uint32(36+n))
	copy(buf[8:], "WAVE")
	copy(buf[12:], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:], 16)
	binary.LittleEndian.PutUint16(buf[20:], 3) // WAVE_FORMAT_IEEE_FLOAT
	binary.LittleEndian.PutUint16(buf[22:], channels)
	binary.LittleEndian.PutUint32(buf[24:], sampleRate)
	binary.LittleEndian.PutUint32(buf[28:], byteRate)
	binary.LittleEndian.PutUint16(buf[32:], blockAlign)
	binary.LittleEndian.PutUint16(buf[34:], bitsPerSample)
	copy(buf[36:], "data")
	binary.LittleEndian.PutUint32(buf[40:], uint32(n))
	_, err = w.Write(buf[:])
	return err
}

// toURI converts an open.spotify.com URL to a spotify: URI.
func toURI(s string) string {
	if strings.HasPrefix(s, "spotify:") {
		return s
	}
	parts := strings.Split(s, "/")
	for i, p := range parts {
		switch p {
		case "track", "album", "playlist", "episode", "show":
			if i+1 < len(parts) {
				id := strings.SplitN(parts[i+1], "?", 2)[0]
				return "spotify:" + p + ":" + id
			}
		}
	}
	return s
}

// safeFilename strips characters that are unsafe in filenames.
func safeFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
