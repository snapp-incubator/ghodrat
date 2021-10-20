.DEFAULT_GOAL := build

.PHONY: lint unit-test build clean

lint:
	golangci-lint run --enable-all

unit-test:
	go test ./...

build:
	go build -o ghodrat ./cmd/ghodrat

clean:
	rm ghodrat
