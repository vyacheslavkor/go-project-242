build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

help:
	go run ./cmd/hexlet-path-size --help

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./tests/...