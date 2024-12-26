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
build: generate
	go build ${LDFLAGS} -o bin/openk ./cmd/openk

.PHONY: install
install: generate
	go install ${LDFLAGS} ./cmd/server

.PHONY: test
test: generate
	go test ./...

.PHONY: test-verbose
test-verbose: generate
	go test -v ./...

.PHONY: proto-vendor
proto-vendor:
	mkdir -p proto/vendor/google/protobuf
	curl -o proto/vendor/google/protobuf/timestamp.proto https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/timestamp.proto
	curl -o proto/vendor/google/protobuf/empty.proto https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/empty.proto

.PHONY: proto-lint
proto-lint:
	cd proto && buf lint

.PHONY: proto-breaking
proto-breaking:
	cd proto && buf breaking --against 'https://github.com/neox5/openk.git#branch=main'

.PHONY: proto
proto:
	cd proto && buf generate

.PHONY: flatf
flatf:
	./scripts/flat_folder.sh
