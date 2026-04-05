//go:build bundled

package binaries

import _ "embed"

//go:embed assets/yt-dlp
var ytdlpBin []byte

//go:embed assets/ffmpeg
var ffmpegBin []byte
