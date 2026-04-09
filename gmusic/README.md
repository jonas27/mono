# gmusic

Terminal music player built with [Bubbletea](https://github.com/charmbracelet/bubbletea).

Supports MP3, WAV, FLAC, OGG.

## Usage

```
make compile
./bin/gmusic
```

Or run directly:

```
make run
```

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate |
| `Enter` | Play file / enter directory |
| `Space` | Pause / resume |
| `+` / `-` | Volume up / down |
| `/` | Filter |
| `q` | Quit |

## Requirements

Requires CGO and ALSA dev headers (Linux):

```
apt install libasound2-dev
```
