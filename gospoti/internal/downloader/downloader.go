package downloader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/jmanser/gospoti/internal/lang"
)

const downloadWorkers = 8

// Downloader orchestrates the download pipeline for Spotify.
type Downloader struct {
	ytdlp          *YTDLPRunner
	spotify        *spotify.Client
	genius         *GeniusClient
	outDir         string
	overrideArtist string
	overrideAlbum  string
	lang           lang.Code
}

// New creates a new Downloader. The Spotify client is only initialised when spotifyID and
// spotifySecret are both non-empty. ytdlpPath and ffmpegPath are the resolved paths to the
// yt-dlp and ffmpeg binaries.
func New(
	spotifyID, spotifySecret, geniusToken, outDir, quality, overrideArtist, overrideAlbum string,
	language lang.Code,
	ytdlpPath, ffmpegPath string,
) (*Downloader, error) {
	d := &Downloader{
		ytdlp:          NewYTDLPRunner(quality, ytdlpPath, ffmpegPath),
		outDir:         outDir,
		overrideArtist: overrideArtist,
		overrideAlbum:  overrideAlbum,
		lang:           language,
	}

	if spotifyID != "" && spotifySecret != "" {
		cfg := &clientcredentials.Config{
			ClientID:     spotifyID,
			ClientSecret: spotifySecret,
			TokenURL:     spotifyauth.TokenURL,
		}
		httpClient := cfg.Client(context.Background())
		d.spotify = spotify.New(httpClient)
	}

	if geniusToken != "" {
		d.genius = NewGeniusClient(geniusToken)
	}

	return d, nil
}

// DownloadSpotify downloads all tracks from a Spotify playlist, album, or single track URL.
func (d *Downloader) DownloadSpotify(ctx context.Context, rawURL string) error {
	if d.spotify == nil {
		return errors.New("spotify API keys are not configured")
	}

	fmt.Println(lang.Get("connecting_spotify", d.lang))

	tracks, subDir, err := d.fetchSpotifyTracks(ctx, rawURL)
	if err != nil {
		return err
	}

	outDir := filepath.Join(d.outDir, dirName(subDir))
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return fmt.Errorf("create output directory %s: %w", outDir, err)
	}

	fmt.Printf(lang.Get("found_tracks", d.lang)+"\n", len(tracks))

	sem := make(chan struct{}, downloadWorkers)
	var wg sync.WaitGroup

	for _, t := range tracks {
		if ctx.Err() != nil {
			break
		}

		sem <- struct{}{}
		wg.Add(1)

		go func(track TrackMeta) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := d.downloadTrack(ctx, outDir, track); err != nil {
				fmt.Fprintf(os.Stderr, "!! DOWNLOAD ERROR for %q: %v\n", track.Title, err)
			}
		}(t)
	}

	wg.Wait()

	fmt.Println(lang.Get("process_finished", d.lang))

	return nil
}

// fetchSpotifyTracks retrieves tracks from a Spotify playlist, album, or single track URL.
// The returned string is the subdirectory name to place downloads in.
func (d *Downloader) fetchSpotifyTracks(ctx context.Context, rawURL string) ([]TrackMeta, string, error) {
	var tracks []TrackMeta

	switch {
	case strings.Contains(rawURL, "/playlist/"):
		id := spotifyIDFromURL(rawURL)
		if id == "" {
			return nil, "", fmt.Errorf("could not parse Spotify playlist ID from URL: %s", rawURL)
		}

		playlist, err := d.spotify.GetPlaylist(ctx, spotify.ID(id))
		if err != nil {
			return nil, "", fmt.Errorf("get Spotify playlist: %w", err)
		}

		page := &playlist.Tracks
		for {
			for i := range page.Tracks {
				tracks = append(tracks, *spotifyTrackToMeta(&page.Tracks[i].Track))
			}
			if err := d.spotify.NextPage(ctx, page); err == spotify.ErrNoMorePages {
				break
			} else if err != nil {
				return nil, "", fmt.Errorf("paginate Spotify playlist: %w", err)
			}
		}

		return tracks, playlist.Name, nil

	case strings.Contains(rawURL, "/album/"):
		id := spotifyIDFromURL(rawURL)
		if id == "" {
			return nil, "", fmt.Errorf("could not parse Spotify album ID from URL: %s", rawURL)
		}

		album, err := d.spotify.GetAlbum(ctx, spotify.ID(id))
		if err != nil {
			return nil, "", fmt.Errorf("get Spotify album: %w", err)
		}

		page := &album.Tracks
		pos := 1
		for {
			for i := range page.Tracks {
				meta := simpleTrackToMeta(&page.Tracks[i], &album.SimpleAlbum)
				meta.TrackNumber = pos
				pos++
				tracks = append(tracks, *meta)
			}
			if err := d.spotify.NextPage(ctx, page); err == spotify.ErrNoMorePages {
				break
			} else if err != nil {
				return nil, "", fmt.Errorf("paginate Spotify album: %w", err)
			}
		}

		return tracks, album.Name, nil

	case strings.Contains(rawURL, "/track/"):
		id := spotifyIDFromURL(rawURL)
		if id == "" {
			return nil, "", fmt.Errorf("could not parse Spotify track ID from URL: %s", rawURL)
		}

		t, err := d.spotify.GetTrack(ctx, spotify.ID(id))
		if err != nil {
			return nil, "", fmt.Errorf("get Spotify track: %w", err)
		}

		return append(tracks, *spotifyTrackToMeta(t)), "songs", nil

	case strings.Contains(rawURL, "/show/"):
		id := spotifyIDFromURL(rawURL)
		if id == "" {
			return nil, "", fmt.Errorf("could not parse Spotify show ID from URL: %s", rawURL)
		}

		show, err := d.spotify.GetShow(ctx, spotify.ID(id))
		if err != nil {
			return nil, "", fmt.Errorf("get Spotify show: %w", err)
		}

		page := &show.Episodes
		pos := 1
		for {
			for i := range page.Episodes {
				meta := episodeToMeta(&page.Episodes[i], &show.SimpleShow)
				meta.TrackNumber = pos
				pos++
				tracks = append(tracks, *meta)
			}
			if err := d.spotify.NextPage(ctx, page); err == spotify.ErrNoMorePages {
				break
			} else if err != nil {
				return nil, "", fmt.Errorf("paginate Spotify show: %w", err)
			}
		}

		return tracks, show.Name, nil

	default:
		return nil, "", fmt.Errorf("unrecognised Spotify URL: %s", rawURL)
	}
}

// downloadTrack downloads a single track via yt-dlp YouTube search and writes metadata.
func (d *Downloader) downloadTrack(ctx context.Context, outDir string, meta TrackMeta) error {
	if d.overrideArtist != "" {
		meta.Artist = d.overrideArtist
	}

	if d.overrideAlbum != "" {
		meta.Album = d.overrideAlbum
	}

	logName := meta.Artist + " - " + meta.Title
	safeName := cleanFilename(logName)
	if meta.TrackNumber > 0 {
		safeName = fmt.Sprintf("%02d - %s", meta.TrackNumber, safeName)
	}
	outPath := filepath.Join(outDir, safeName+".mp3")

	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("%s %s\n", lang.Get("song_skipped", d.lang), logName)
		return nil
	}

	fmt.Printf("%s %s\n", lang.Get("downloading_song", d.lang), logName)

	searchQuery := meta.Artist + " - " + meta.Title + " audio"

	// Download to a temp path so yt-dlp's own filename sanitization cannot cause
	// a mismatch between the path we pass and the file it actually writes.
	tmpFile, err := os.CreateTemp(outDir, "dl-*.mp3")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	_ = os.Remove(tmpPath)

	if err := d.ytdlp.DownloadSearch(ctx, searchQuery, tmpPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, outPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename download to %s: %w", outPath, err)
	}

	if err := d.writeFullMeta(ctx, outPath, meta); err != nil {
		fmt.Fprintf(os.Stderr, "   -> metadata error: %v\n", err)
	}

	fmt.Printf("%s %s\n", lang.Get("song_done", d.lang), logName)

	return nil
}

// writeFullMeta writes complete ID3 metadata including optional lyrics from Genius.
func (d *Downloader) writeFullMeta(ctx context.Context, path string, meta TrackMeta) error {
	if d.genius != nil {
		lyrics, err := d.genius.GetLyrics(ctx, meta.Artist, meta.Title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "   -> INFO: Error fetching lyrics: %v\n", err)
		} else {
			meta.Lyrics = lyrics
		}
	}

	return WriteMetadata(ctx, path, meta)
}

// spotifyTrackToMeta converts a Spotify FullTrack into a TrackMeta.
func spotifyTrackToMeta(t *spotify.FullTrack) *TrackMeta {
	artist := ""
	if len(t.Artists) > 0 {
		artist = t.Artists[0].Name
	}

	coverURL := ""
	if len(t.Album.Images) > 0 {
		coverURL = t.Album.Images[0].URL
	}

	return &TrackMeta{
		Title:       t.Name,
		Artist:      artist,
		Album:       t.Album.Name,
		TrackNumber: int(t.TrackNumber),
		ReleaseDate: t.Album.ReleaseDate,
		CoverURL:    coverURL,
	}
}

// simpleTrackToMeta converts a Spotify SimpleTrack and its album into a TrackMeta.
func simpleTrackToMeta(t *spotify.SimpleTrack, album *spotify.SimpleAlbum) *TrackMeta {
	artist := ""
	if len(t.Artists) > 0 {
		artist = t.Artists[0].Name
	}

	coverURL := ""
	if len(album.Images) > 0 {
		coverURL = album.Images[0].URL
	}

	return &TrackMeta{
		Title:       t.Name,
		Artist:      artist,
		Album:       album.Name,
		TrackNumber: int(t.TrackNumber),
		ReleaseDate: album.ReleaseDate,
		CoverURL:    coverURL,
	}
}

// episodeToMeta converts a Spotify EpisodePage and its show into a TrackMeta.
func episodeToMeta(e *spotify.EpisodePage, show *spotify.SimpleShow) *TrackMeta {
	coverURL := ""
	if len(e.Images) > 0 {
		coverURL = e.Images[0].URL
	} else if len(show.Images) > 0 {
		coverURL = show.Images[0].URL
	}

	return &TrackMeta{
		Title:       e.Name,
		Artist:      show.Publisher,
		Album:       show.Name,
		ReleaseDate: e.ReleaseDate,
		CoverURL:    coverURL,
	}
}

// spotifyIDFromURL extracts the Spotify resource ID from a Spotify URL.
func spotifyIDFromURL(rawURL string) string {
	parts := strings.Split(rawURL, "/")
	for i, p := range parts {
		if (p == "track" || p == "playlist" || p == "album" || p == "show") && i+1 < len(parts) {
			id := parts[i+1]
			// Remove any query string.
			id = strings.SplitN(id, "?", 2)[0]
			return id
		}
	}

	return ""
}

// dirName converts a name to a lowercase, dash-separated directory name.
func dirName(name string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if b.Len() > 0 {
			b.WriteRune('-')
		}
	}

	return strings.TrimRight(b.String(), "-")
}

// cleanFilename removes characters that are unsafe for file names.
func cleanFilename(name string) string {
	var b strings.Builder

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			b.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}
