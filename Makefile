GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=docker-stats-collector

CC=x86_64-linux-musl-gcc
CXX=x86_64-linux-musl-g++

all: build

build: clean-templ generate-templ
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CC=$(CC) CXX=$(CXX) \
	$(GOBUILD) -o $(BINARY_NAME) -v \
	-ldflags='-s -w -linkmode external -extldflags "-static"' \
	-tags netgo,osusergo

build-debug: clean-templ generate-templ
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CC=$(CC) CXX=$(CXX) \
	$(GOBUILD) -o $(BINARY_NAME) -v \
	-ldflags='-linkmode external -extldflags "-static" -X main.BuildTime=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")' \
	-gcflags="all=-N -l" \
	-tags netgo,osusergo

clean-templ:
	find . -name "*_templ.go" -type f -delete

generate-templ:
	templ generate

generate-templ-watch:
	templ generate --watch

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

build-local:
	$(GOBUILD) -o $(BINARY_NAME) -v

check-binary:
	file $(BINARY_NAME)

deps:
	$(GOGET) ./...

run:
	./$(BINARY_NAME)

run-dev:
	go run main.go -db ./stats.db

test:
	go test ./...

.PHONY: all build build-debug clean clean-templ generate-templ generate-templ-watch build-local check-binary deps run run-dev test
