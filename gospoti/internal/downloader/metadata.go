package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bogem/id3v2/v2"
)

// TrackMeta holds all metadata used when writing ID3 tags to an MP3 file.
type TrackMeta struct {
	// Title is the track title.
	Title string
	// Artist is the primary artist name.
	Artist string
	// Album is the album name.
	Album string
	// TrackNumber is the position within the album.
	TrackNumber int
	// ReleaseDate is the release date string (e.g. "2021-06-25").
	ReleaseDate string
	// CoverURL is a URL pointing to the cover art image.
	CoverURL string
	// Lyrics is the plain-text lyrics string.
	Lyrics string
}

// WriteMetadata writes ID3v2 tags to the MP3 file at path using the provided metadata.
func WriteMetadata(ctx context.Context, path string, meta TrackMeta) error {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("open MP3 for tagging %s: %w", path, err)
	}

	defer func() { _ = tag.Close() }()

	tag.SetTitle(meta.Title)
	tag.SetArtist(meta.Artist)
	tag.SetAlbum(meta.Album)
	tag.SetYear(releaseYear(meta.ReleaseDate))

	if meta.TrackNumber > 0 {
		tag.AddTextFrame(tag.CommonID("Track number/Position in set"), id3v2.EncodingUTF8, strconv.Itoa(meta.TrackNumber))
	}

	if meta.ReleaseDate != "" {
		tag.AddTextFrame("TDRC", id3v2.EncodingUTF8, meta.ReleaseDate)
	}

	if meta.CoverURL != "" {
		pic, mimeType, coverErr := fetchCover(ctx, meta.CoverURL)
		if coverErr == nil && len(pic) > 0 {
			tag.AddAttachedPicture(id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    mimeType,
				PictureType: id3v2.PTFrontCover,
				Description: "Cover",
				Picture:     pic,
			})
		}
	}

	if meta.Lyrics != "" {
		tag.AddUnsynchronisedLyricsFrame(id3v2.UnsynchronisedLyricsFrame{
			Encoding:          id3v2.EncodingUTF8,
			Language:          "eng",
			ContentDescriptor: "Lyrics",
			Lyrics:            meta.Lyrics,
		})
	}

	if err := tag.Save(); err != nil {
		return fmt.Errorf("save ID3 tags to %s: %w", path, err)
	}

	return nil
}

// WriteBasicMetadata writes minimal ID3v2 tags (title and artist) to the MP3 at path.
func WriteBasicMetadata(path, title, artist string) error {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("open MP3 for basic tagging %s: %w", path, err)
	}

	defer func() { _ = tag.Close() }()

	tag.SetTitle(title)
	tag.SetArtist(artist)

	if err := tag.Save(); err != nil {
		return fmt.Errorf("save basic ID3 tags to %s: %w", path, err)
	}

	return nil
}

func fetchCover(ctx context.Context, coverURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, coverURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create cover art request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, "", fmt.Errorf("fetch cover art: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read cover art body: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	return data, mimeType, nil
}

// releaseYear extracts the four-digit year from a release date string.
// Returns the original string if it is already 4 characters or fewer.
func releaseYear(date string) string {
	if len(date) >= 4 {
		return date[:4]
	}

	return date
}
