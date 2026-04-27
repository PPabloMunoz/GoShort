# GoShort - URL shortener service

# Variables
mainFile := './cmd/main.go'
executable := './bin/GoShort'

default:
  @just --list

build:
  go build -o {{executable}} {{mainFile}}

run: build
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
