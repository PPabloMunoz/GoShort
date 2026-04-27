# go-shorten - URL shortener service

mainFile := './cmd/main.go'
executable := './bin/go-shorten'

default:
  @just --list

build:
  go build -o {{executable}} {{mainFile}}

run: build
  {{executable}}

release: build
  GIN_MODE=release {{executable}}

dev:
  air -c .air.toml

test:
  go test -v ./...

check:
  go vet ./...
  golangci-lint run ./...

fmt:
  go fmt ./...
  gofumpt -w .

deps:
  go mod tidy
  go mod vendor
