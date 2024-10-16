GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=docker-stats-collector

CC=x86_64-linux-musl-gcc
CXX=x86_64-linux-musl-g++

all: build

build:
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CC=$(CC) CXX=$(CXX) \
	$(GOBUILD) -o $(BINARY_NAME) -v \
	-ldflags='-s -w -linkmode external -extldflags "-static"' \
	-tags netgo,osusergo

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

build-local:
	$(GOBUILD) -o $(BINARY_NAME)-local -v

check-binary:
	file $(BINARY_NAME)

deps:
	$(GOGET) ./...

run:
	./$(BINARY_NAME)

run-dev:
	go run main.go -interval 5 -db ./database.db

.PHONY: all build clean build-local check-binary deps run run-dev
