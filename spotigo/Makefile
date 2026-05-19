.PHONY: build lint clean

build:
	CGO_ENABLED=0 go build -o spotigo .

lint:
	golangci-lint run ./...

clean:
	rm -f spotigo
