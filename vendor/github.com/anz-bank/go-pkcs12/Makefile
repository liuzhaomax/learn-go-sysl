.PHONY: all
all: build lint test

build:
	go build ./...

lint:
	golangci-lint run

test:
	go test ./...
