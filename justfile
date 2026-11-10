BUILD_DIR := ".build"

alias b := build
alias l := lint

_default:
    @just --list

_mkdir:
    @mkdir -p {{ BUILD_DIR }}

# Lint the repository
lint:
    @gofmt -w .
    @golangci-lint run

# Build the binary
build:
    @go build -o {{ BUILD_DIR }}/go-tree main.go
