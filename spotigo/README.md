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

## Docker

Build the image (bundles `ffmpeg` and the required FLAC/Vorbis/ALSA libraries):

```
docker build -t spotigo:latest .
```

Run with the credentials file and downloads directory bind-mounted so they
persist on the host. Use `docker compose` for convenience:

```
UID=$(id -u) GID=$(id -g) docker compose run --rm spotigo spotify:track:xxxx
```

Or with plain `docker run`:

```
touch .creds.yaml            # first run only, so the bind mount is a file
docker run --rm -it \
  --user "$(id -u):$(id -g)" \
  -p 5173:5173 \
  -v "$PWD/.creds.yaml:/data/.creds.yaml" \
  -v "$PWD/downloads:/data/downloads" \
  spotigo:latest spotify:track:xxxx
```

The `-p 5173:5173` mapping is only needed for the first-run OAuth login (see
below); the container binds the callback server on all interfaces so
`http://localhost:5173/login` on the host reaches it. Running as your host user
keeps the `.creds.yaml` (mode `0600`) readable and downloaded files owned by you.

## Authentication

On first run, spotigo starts a local OAuth2 server and prints:

```
No credentials — open http://localhost:5173/login in your browser.
```

Open that URL, log in with your Spotify account, and credentials are saved automatically to `.creds.yaml`. Subsequent runs use the saved credentials without prompting.

### Credentials file structure

The file at `-creds` (default `.creds.yaml`) is YAML with two fields, written
automatically after login with mode `0600`:

```yaml
username: your_spotify_username
stored_credentials: <base64-encoded token>
```

| Field | Type | Description |
|-------|------|-------------|
| `username` | string | Spotify account username returned by the OAuth login. |
| `stored_credentials` | string | Base64-encoded reusable auth token (standard encoding). Used to re-authenticate on subsequent runs without logging in again. |

Both fields are optional: an empty or missing file triggers the interactive
OAuth login, which then populates both. The token in `stored_credentials`
grants full access to your Spotify account — **do not commit this file**
(it is already in `.gitignore`).
