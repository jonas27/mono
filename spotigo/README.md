# spotigo

Streams tracks from Spotify and writes them as MP3 files (320 kbps). Requires `ffmpeg` on `$PATH`.

## Usage

```
spotigo [flags] <spotify-url-or-uri>
```

Accepts Spotify URLs (`https://open.spotify.com/track/...`) or URIs (`spotify:track:...`).
Albums and playlists are downloaded in full; tracks are downloaded individually.

**Flags**

| Flag | Default | Description |
|------|---------|-------------|
| `-creds` | `.creds.yaml` | Credentials file |
| `-output` | `downloads` | Output directory |
| `-album` | _(none)_ | Override subdirectory name for albums/playlists |

## Authentication

On first run, spotigo starts a local OAuth2 server and prints:

```
No credentials — open http://localhost:5173/login in your browser.
```

Open that URL, log in with your Spotify account, and credentials are saved automatically to `.creds.yaml`. Subsequent runs use the saved credentials without prompting.

The credentials file looks like this after login:

```yaml
username: your_spotify_username
stored_credentials: <base64-encoded token>
```

Do not commit this file — it grants full access to your Spotify account.
