.PHONY: build lint clean docker-build docker-push

IMAGE    ?= jonas27/spotigo
VERSION  ?= latest
# VERSION  ?= $(shell git describe --tags --always --dirty)
PLATFORM := linux/amd64

build:
	CGO_ENABLED=0 go build -o spotigo .

lint:
	golangci-lint run ./...

clean:
	rm -f spotigo

docker-build:
	docker build --platform $(PLATFORM) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

docker-push: docker-build
	docker push $(IMAGE):$(VERSION) --platform $(PLATFORM)
