# syntax=docker/dockerfile:1

# --- build stage ---
FROM golang:1.26-alpine AS build

# CGO is required: go-librespot links libFLAC, libogg/libvorbis, and ALSA.
RUN apk add --no-cache \
    gcc musl-dev pkgconf \
    flac-dev libogg-dev libvorbis-dev alsa-lib-dev

WORKDIR /src

# Cache module downloads separately from the source build.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/spotigo .

# --- runtime stage ---
FROM alpine:3.21

# ffmpeg encodes the Opus output; the rest are the shared libs the cgo
# binary links against; ca-certificates for TLS to Spotify.
RUN apk add --no-cache \
    ffmpeg ca-certificates \
    flac-libs libogg libvorbis alsa-lib \
    && adduser -D -u 10001 spotigo

COPY --from=build /out/spotigo /usr/local/bin/spotigo

# Credentials (.creds.yaml) and downloads live under the working dir, which is
# mounted as a volume so they persist across runs and survive image rebuilds.
WORKDIR /data
RUN chown spotigo:spotigo /data
USER spotigo

# OAuth2 callback server for first-run interactive login.
EXPOSE 5173

ENTRYPOINT ["spotigo"]
CMD ["-help"]
