NAME    := git-hours
PACKAGE := github.com/trinhminhtriet/$(NAME)
DATE    :=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT     := $(shell [ -d .git ] && git rev-parse --short HEAD)
VERSION := $(shell git describe --tags)

default: build

tidy:
	go get -u && go mod tidy

build:
	CGO_ENABLED=0 go build \
	-a -tags netgo -o dist/${NAME} main.go

build-link:
	CGO_ENABLED=0 go build \
		-a -tags netgo -o dist/${NAME} main.go
	ln -sf ${PWD}/dist/${NAME} /usr/local/bin/${NAME}

release:
	goreleaser build --clean --snapshot --single-target

release-all:
	goreleaser build --clean --snapshot

link:
	ln -sf ${PWD}/dist/${NAME} /usr/local/bin/${NAME}
	which ${NAME}

clean:
	$(RM) -rf dist

.PHONY: default tidy build build-link release release-all
