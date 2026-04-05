# gospoti

A command-line tool that downloads music as MP3 files from Spotify, YouTube, Apple Music, or Deezer.

## Prerequisites

- Spotify API credentials (required for Spotify source; optional but recommended for YouTube/Apple/Deezer to enrich metadata)
- Genius API token (optional, for embedding lyrics)

## Installation

### Self-contained binary (recommended)

Builds a single binary with yt-dlp and ffmpeg bundled inside — no external tools required at runtime:

https://open.spotify.com/album/6ICcWbsZxt5fVG02IbZuZg?si=RYyiiztZSjCVbGBmBGOrOg
```sh
make build
```

The binary is ~80 MB. On first run the bundled tools are extracted to `~/.cache/gospoti/` and reused on subsequent runs.

> **Note:** The bundled ffmpeg is from [BtbN's static GPL builds](https://github.com/BtbN/FFmpeg-Builds). macOS users need to provide their own static ffmpeg (see `make fetch-deps` error message).

### Lean binary

Requires `yt-dlp` and `ffmpeg` in PATH:

```sh
make build-lean
# or: go build -o gospoti .
```

### Via go install

```sh
go install github.com/jmanser/gospoti@latest
```

Requires Go 1.23+, and yt-dlp + ffmpeg in PATH at runtime.

## Configuration

Run the interactive setup to store API keys in `config.json`:

```sh
gospoti --setup
```

If you pass a Spotify URL without keys configured, setup runs automatically.

`config.json` is written to the working directory with the following fields:

| Field | Default | Description |
|---|---|---|
| `api_keys.spotify_id` | `""` | Spotify application Client ID |
| `api_keys.spotify_secret` | `""` | Spotify application Client Secret |
| `api_keys.genius_token` | `""` | Genius API access token (optional) |
| `download_path` | `"downloads"` | Directory where MP3 files are saved |
| `quality` | `"192"` | Audio bitrate in kbps |
| `source_type` | `"Spotify"` | Default source (informational only) |

Obtain Spotify credentials at [developer.spotify.com](https://developer.spotify.com/dashboard).
Obtain a Genius token at [genius.com/api-clients](https://genius.com/api-clients).

## Usage

```
gospoti [flags] <URL>
```

| Flag | Default | Description |
|---|---|---|
| `--source` | auto-detected | Force source type: `spotify`, `youtube`, `apple`, `deezer` |
| `--output` | config value or `downloads` | Output directory |
| `--quality` | config value or `192` | Audio bitrate: `128`, `192`, `256`, `320` |
| `--lang` | `en` | UI language: `en`, `bg`, `es` |
| `--artist` | — | Override artist name for all downloads (overrides Genius/Spotify metadata) |
| `--setup` | — | Configure API keys interactively, then exit |

The source type is auto-detected from the URL when `--source` is omitted.

### Examples

```sh
# Spotify playlist
gospoti https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M

# Spotify track
gospoti https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC

# YouTube video or playlist
gospoti https://www.youtube.com/watch?v=dQw4w9WgXcQ
gospoti https://www.youtube.com/playlist?list=PLxxx

# Apple Music album or song
gospoti https://music.apple.com/us/album/some-album/123456789
gospoti "https://music.apple.com/us/album/some-album/123456789?i=987654321"

# Apple Music playlist
gospoti https://music.apple.com/us/playlist/some-playlist/pl.abc123

# Deezer track, album, or playlist
gospoti https://www.deezer.com/track/123456
gospoti https://www.deezer.com/album/123456
gospoti https://www.deezer.com/playlist/123456

# Custom output directory and quality
gospoti --output ~/Music --quality 320 https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC
```

Place a `youtube-cookies.txt` (Netscape format) in the working directory to pass cookies to yt-dlp when needed.

## Source types

| Source | Track list resolution | Audio source |
|---|---|---|
| Spotify | Spotify API (requires credentials) | YouTube search via yt-dlp |
| YouTube | yt-dlp (video or playlist URL) | Direct YouTube download |
| Apple Music | iTunes Lookup API (albums/songs) or HTML scrape (playlists) | YouTube search via yt-dlp |
| Deezer | HTML scrape (track/album/playlist pages) | YouTube search via yt-dlp |

## How it works

- For Spotify and non-YouTube sources, the tool resolves a track list (title + artist) from the source, then searches YouTube for each track using `ytsearch1:artist - title audio` and downloads the first result via yt-dlp.
- For YouTube URLs, yt-dlp fetches metadata directly; if Spotify credentials are available, the tool fuzzy-matches the video title against Spotify search results to enrich ID3 tags.
- yt-dlp extracts audio and ffmpeg re-encodes it to MP3 at the requested bitrate.
- ID3v2 tags (title, artist, album, track number, release date, cover art, and optionally lyrics from Genius) are written to each file after download. Files that already exist are skipped.
