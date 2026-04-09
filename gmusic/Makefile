VERSION ?= $(shell git describe --always)

GO_VERSION = $(shell go version 2>/dev/null | cut -d' ' -f3)
SOURCES := $(shell find . -name '*.go' -o -name 'go.*')

GOLANGCI_LINT := go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)
export CGO_ENABLED = 1
BUILD_DATE ?= $(shell date -Iseconds)
GITHUB_SHA ?= $(shell git rev-parse HEAD)
GIT_TREE_STATE=$(shell test -n "`git status --porcelain`" && echo 'dirty' || echo 'clean')

override EXECUTABLE := bin/gmusic

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: compile
compile: $(EXECUTABLE) ## Compile binary

$(EXECUTABLE): $(SOURCES)
	go build -trimpath -o $(EXECUTABLE) .

.PHONY: run
run: ## Run the music player
	go run .

.PHONY: lint
lint: ## Run golangci-lint against code
	$(GOLANGCI_LINT) run -v ./...

.PHONY: lint-fix
lint-fix: ## Fix golangci-lint issues in code
	$(GOLANGCI_LINT) run -v ./... --fix

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: verify
verify: verify-tidy

.PHONY: verify-tidy
verify-tidy: tidy ## Verify that the dependencies are tidied.
	@if !(git diff --quiet HEAD); then \
		echo "run go mod tidy"; exit 1; \
	fi

.PHONY: test
test: ## Run tests
	go test -v ./... -race -shuffle on

##@ Install

INSTALL_BIN := $(HOME)/.local/bin
INSTALL_APP := $(HOME)/.local/share/applications
MUSIC_DIR   ?=

INSTALL_ICONS := $(HOME)/.local/share/icons/hicolor/scalable/apps

.PHONY: install
install: compile ## Install binary, icon and desktop launcher (optional: MUSIC_DIR=/path/to/music)
	install -Dm755 $(EXECUTABLE) $(INSTALL_BIN)/gmusic
	install -Dm644 gmusic.svg $(INSTALL_ICONS)/gmusic.svg
	sed 's|GMUSIC_BIN|$(INSTALL_BIN)/gmusic$(if $(MUSIC_DIR), -dir $(MUSIC_DIR),)|g' gmusic.desktop \
		> $(INSTALL_APP)/gmusic.desktop

.PHONY: uninstall
uninstall: ## Remove installed binary, icon and desktop launcher
	rm -f $(INSTALL_BIN)/gmusic $(INSTALL_APP)/gmusic.desktop $(INSTALL_ICONS)/gmusic.svg

##@ Build

.PHONY: clean
clean: ## Remove built executables
	rm -rf bin/
