ASSETS := internal/binaries/assets
OS     := $(shell uname -s)
ARCH   := $(shell uname -m)

.PHONY: build build-lean fetch-deps clean

## build: download yt-dlp + ffmpeg and compile a self-contained binary
build: fetch-deps
	go build -tags bundled -o gospoti .

## build-lean: compile without bundling (requires yt-dlp and ffmpeg in PATH)
build-lean:
	go build -o gospoti .

## fetch-deps: download yt-dlp and ffmpeg static binaries for the current platform
fetch-deps: $(ASSETS)/yt-dlp $(ASSETS)/ffmpeg

$(ASSETS):
	mkdir -p $(ASSETS)

$(ASSETS)/yt-dlp: | $(ASSETS)
ifeq ($(OS),Linux)
  ifeq ($(ARCH),x86_64)
	curl -fsSL -o $@ \
		https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
  else ifeq ($(ARCH),aarch64)
	curl -fsSL -o $@ \
		https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64
  else
	$(error Unsupported Linux arch: $(ARCH))
  endif
else ifeq ($(OS),Darwin)
	curl -fsSL -o $@ \
		https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos
else
	$(error Unsupported OS: $(OS). Place yt-dlp at $(ASSETS)/yt-dlp manually.)
endif
	chmod +x $@

# Note: BtbN static ffmpeg builds are GPL-licensed. See https://github.com/BtbN/FFmpeg-Builds.
$(ASSETS)/ffmpeg: | $(ASSETS)
ifeq ($(OS),Linux)
  ifeq ($(ARCH),x86_64)
	curl -fsSL \
		https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz \
		| tar -xJf - --strip-components=2 --wildcards -C $(ASSETS) '*/bin/ffmpeg'
  else ifeq ($(ARCH),aarch64)
	curl -fsSL \
		https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linuxarm64-gpl.tar.xz \
		| tar -xJf - --strip-components=2 --wildcards -C $(ASSETS) '*/bin/ffmpeg'
  else
	$(error Unsupported Linux arch: $(ARCH))
  endif
else ifeq ($(OS),Darwin)
	@echo "ERROR: Automatic ffmpeg download not supported on macOS."
	@echo "Place a static ffmpeg binary at $(ASSETS)/ffmpeg, then re-run 'make build'."
	@echo "One source: https://evermeet.cx/ffmpeg/"
	@exit 1
else
	$(error Unsupported OS: $(OS). Place ffmpeg at $(ASSETS)/ffmpeg manually.)
endif
	chmod +x $@

## clean: remove the built binary and downloaded assets
clean:
	rm -f gospoti
	rm -rf $(ASSETS)
