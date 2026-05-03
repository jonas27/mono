.PHONY: build lint clean

build:
	CGO_ENABLED=0 go build -mod=vendor -o spotigo .

lint:
	golangci-lint run ./...

clean:
	rm -f spotigo
