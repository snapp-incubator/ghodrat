.DEFAULT_GOAL := build

.PHONY: lint unit-test build clean update

lint:
	golangci-lint run --enable-all

unit-test:
	go test ./...

build:
	go build -o ghodrat ./cmd/ghodrat

update:
	go get -u ./cmd/ghodrat

clean:
	rm ghodrat
