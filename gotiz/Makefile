.PHONY: build lint clean

build:
	go build -o gotiz .

lint:
	golangci-lint run ./...

clean:
	rm -f gotiz
