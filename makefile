# Build variables
VERSION ?= $(shell git describe --tags --always --dirty)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_USER=$(shell whoami)

# Go build flags
LDFLAGS=-ldflags "\
	-X github.com/neox5/openk/internal/buildinfo.Version=${VERSION} \
	-X github.com/neox5/openk/internal/buildinfo.GitCommit=${GIT_COMMIT} \
	-X github.com/neox5/openk/internal/buildinfo.BuildTime=${BUILD_TIME} \
	-X github.com/neox5/openk/internal/buildinfo.BuildUser=${BUILD_USER}"

.PHONY: build
build:
	go build ${LDFLAGS} -o bin/openk ./cmd/openk

.PHONY: install
install:
	go install ${LDFLAGS} ./cmd/server

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: proto
proto:
	buf generate

.PHONY: proto-lint
proto-lint:
	buf lint

.PHONY: proto-breaking
proto-breaking:
	buf breaking --against 'https://github.com/neox5/openk.git#branch=main'
