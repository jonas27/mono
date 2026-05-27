// Package player handles Spotify auth and streams audio to Opus files.
package player

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	librespot "github.com/devgianlu/go-librespot"
	"github.com/devgianlu/go-librespot/ap"
	lsplayer "github.com/devgianlu/go-librespot/player"
	devicespb "github.com/devgianlu/go-librespot/proto/spotify/connectstate/devices"
	metadatapb "github.com/devgianlu/go-librespot/proto/spotify/metadata"
	"github.com/devgianlu/go-librespot/session"
	"github.com/devgianlu/go-librespot/spclient"
	"golang.org/x/sync/errgroup"

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
		slog.Info("no credentials stored", "login_url", fmt.Sprintf("http://localhost:%d/login", oauthPort))
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
			slog.Warn("save credentials", "error", err)
		} else {
			slog.Info("credentials saved", "path", credsFile)
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

// Close releases the session.
func (s *Session) Close() {
	s.sess.Close()
}

// Run streams the Spotify URL or URI, writing Opus files into outDir.
// Albums and playlists create a subdirectory; tracks write directly into outDir.
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
		rawPath, meta, err := s.downloadTrackRaw(ctx, *id, outDir)
		if err != nil {
			return err
		}
		defer os.Remove(rawPath)
		finalPath := filepath.Join(outDir, trackFilename(meta))
		slog.Info("encoding track", "file", filepath.Base(finalPath))
		return encodeOpusFromFile(ctx, rawPath, finalPath, meta)
	default:
		return s.streamContext(ctx, uri, outDir, albumOverride)
	}
}

// streamContext resolves an album/playlist and downloads every track sequentially.
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

	slog.Info("found tracks", "count", len(trackURIs), "dir", dir)

	g, gctx := errgroup.WithContext(ctx)

	for i, trackURI := range trackURIs {
		if gctx.Err() != nil {
			break
		}
		id, err := parseSpotifyID(trackURI)
		if err != nil {
			slog.Warn("skip track", "uri", trackURI, "error", err)
			continue
		}
		rawPath, meta, err := s.downloadTrackRaw(gctx, id, dir)
		if err != nil {
			slog.Error("download track failed", "index", i+1, "error", err)
			continue
		}
		finalPath := filepath.Join(dir, trackFilename(meta))
		slog.Info("encoding track", "file", filepath.Base(finalPath))
		g.Go(func() error {
			defer os.Remove(rawPath)
			return encodeOpusFromFile(gctx, rawPath, finalPath, meta)
		})
	}

	return g.Wait()
}

// downloadTrackRaw streams decoded PCM to a temp file in dir, retrying on audio key rate-limit errors.
// Returns the temp file path and track metadata; caller is responsible for removing the file.
func (s *Session) downloadTrackRaw(ctx context.Context, id librespot.SpotifyId, dir string) (string, trackMeta, error) {
	rawPath, meta, err := s.downloadRaw(ctx, id, dir)
	for i, d := range []time.Duration{20 * time.Second, 40 * time.Second, 60 * time.Second} {
		if err == nil || ctx.Err() != nil || !strings.Contains(err.Error(), "aes key") {
			break
		}
		slog.Warn("audio key rate-limited, retrying", "delay", d, "attempt", i+2)
		select {
		case <-ctx.Done():
			return "", trackMeta{}, ctx.Err()
		case <-time.After(d):
		}
		rawPath, meta, err = s.downloadRaw(ctx, id, dir)
	}
	return rawPath, meta, err
}

// downloadRaw performs one attempt: streams decoded PCM via a FIFO into a temp file in dir.
// Returns the temp file path and track metadata; caller is responsible for removing the file.
func (s *Session) downloadRaw(ctx context.Context, id librespot.SpotifyId, dir string) (rawPath string, meta trackMeta, err error) {
	fifoPath := filepath.Join(dir, ".fifo-"+randomHex(8))
	if err = syscall.Mkfifo(fifoPath, 0o600); err != nil {
		return "", trackMeta{}, fmt.Errorf("create fifo: %w", err)
	}
	defer os.Remove(fifoPath)

	// Open read end non-blocking so the player's O_WRONLY|O_NONBLOCK open succeeds.
	rfd, ferr := syscall.Open(fifoPath, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if ferr != nil {
		return "", trackMeta{}, fmt.Errorf("open fifo: %w", ferr)
	}
	if ferr = syscall.SetNonblock(rfd, false); ferr != nil {
		_ = syscall.Close(rfd)
		return "", trackMeta{}, fmt.Errorf("set blocking: %w", ferr)
	}
	reader := os.NewFile(uintptr(rfd), fifoPath)

	pl, ferr := lsplayer.NewPlayer(&lsplayer.Options{
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
	if ferr != nil {
		_ = reader.Close()
		return "", trackMeta{}, fmt.Errorf("create player: %w", ferr)
	}

	events := pl.Receive()

	stream, ferr := pl.NewStream(ctx, http.DefaultClient, id, bitrate, 0)
	if ferr != nil {
		_ = reader.Close()
		pl.Close()
		return "", trackMeta{}, fmt.Errorf("load stream: %w", ferr)
	}

	s.downloadCover(ctx, dir, stream)
	meta = streamTrackMeta(stream)

	rawPath = filepath.Join(dir, ".raw-"+randomHex(8))
	copyDone := make(chan error, 1)
	go func() {
		f, cerr := os.Create(rawPath)
		if cerr != nil {
			_ = reader.Close()
			copyDone <- cerr
			return
		}
		_, cerr = io.Copy(f, reader)
		_ = f.Close()
		_ = reader.Close()
		if cerr != nil {
			_ = os.Remove(rawPath)
		}
		copyDone <- cerr
	}()

	if ferr = pl.SetPrimaryStream(stream.Source, false, false); ferr != nil {
		pl.Close()
		<-copyDone
		_ = os.Remove(rawPath)
		return "", trackMeta{}, fmt.Errorf("set stream: %w", ferr)
	}
	if ferr = pl.Play(); ferr != nil {
		pl.Close()
		<-copyDone
		_ = os.Remove(rawPath)
		return "", trackMeta{}, fmt.Errorf("play: %w", ferr)
	}

	for {
		select {
		case <-ctx.Done():
			pl.Close()
			<-copyDone
			_ = os.Remove(rawPath)
			return "", trackMeta{}, ctx.Err()
		case ev, ok := <-events:
			if !ok || ev.Type == lsplayer.EventTypeStop || ev.Type == lsplayer.EventTypeNotPlaying {
				pl.Close()
				if ferr = <-copyDone; ferr != nil {
					return "", trackMeta{}, ferr
				}
				return rawPath, meta, nil
			}
		}
	}
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

// trackFilename builds "Title.opus".
func trackFilename(meta trackMeta) string {
	return safeFilename(meta.Title) + ".opus"
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

type trackMeta struct {
	Title       string
	Artist      string
	Album       string
	TrackNumber int
	DiscNumber  int
}

func streamTrackMeta(stream *lsplayer.Stream) trackMeta {
	m := trackMeta{Title: stream.Media.Name()}
	if stream.Media.IsTrack() {
		t := stream.Media.Track()
		if len(t.GetArtist()) > 0 {
			m.Artist = t.GetArtist()[0].GetName()
		}
		if alb := t.GetAlbum(); alb != nil {
			m.Album = alb.GetName()
		}
		if n := int(t.GetNumber()); n > 0 {
			m.TrackNumber = n
		}
		if d := int(t.GetDiscNumber()); d > 0 {
			m.DiscNumber = d
		}
	}
	return m
}

// encodeOpus encodes f32le PCM from r into an Opus file at path using ffmpeg.
func encodeOpus(ctx context.Context, path string, r io.Reader, meta trackMeta) error {
	args := []string{
		"-hide_banner", "-loglevel", "error",
		"-f", "f32le", "-ar", "44100", "-ac", "2", "-i", "pipe:0",
		"-c:a", "libopus", "-b:a", "64k",
	}
	if meta.Title != "" {
		args = append(args, "-metadata", "title="+meta.Title)
	}
	if meta.Artist != "" {
		args = append(args, "-metadata", "artist="+meta.Artist)
	}
	if meta.Album != "" {
		args = append(args, "-metadata", "album="+meta.Album)
	}
	if meta.TrackNumber > 0 {
		args = append(args, "-metadata", fmt.Sprintf("track=%d", meta.TrackNumber))
	}
	if meta.DiscNumber > 0 {
		args = append(args, "-metadata", fmt.Sprintf("disc=%d", meta.DiscNumber))
	}
	args = append(args, "-f", "ogg", "-y", path)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// encodeOpusFromFile opens rawPath (f32le PCM) and encodes it to an Opus file at outPath.
func encodeOpusFromFile(ctx context.Context, rawPath, outPath string, meta trackMeta) error {
	f, err := os.Open(rawPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return encodeOpus(ctx, outPath, f, meta)
}

// parseSpotifyID parses a spotify: URI into a SpotifyId.
func parseSpotifyID(uri string) (librespot.SpotifyId, error) {
	p := strings.SplitN(uri, ":", 3)
	if len(p) != 3 {
		return librespot.SpotifyId{}, fmt.Errorf("invalid URI: %s", uri)
	}
	id, err := librespot.SpotifyIdFromBase62(librespot.SpotifyIdType(p[1]), p[2])
	if err != nil {
		return librespot.SpotifyId{}, err
	}
	return *id, nil
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

// safeFilename strips filesystem-unsafe characters, preserving original casing and spaces.
func safeFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == 0 || strings.ContainsRune(`/\:*?"<>|`, r) {
			continue
		}
		b.WriteRune(r)
	}
	return strings.Trim(b.String(), " .")
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
