GROUP := github.com/TheoBrigitte
NAME := kimsufi-notifier

PKG := ${GROUP}/${NAME}

BIN := ${NAME}

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

VERSION := $(shell git describe --always --long --dirty || date)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

build:
	GOOS=${OSTYPE} go build -v -o ${BIN} -ldflags=" \
	-s -w \
	-X github.com/prometheus/common/version.Version=$(shell git describe --tags) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/prometheus/common/version.BuildDate=$(shell date --utc +%FT%T)" \
	${PKG}

install:
	GOOS=${OSTYPE} go install -v -ldflags="-s -w -X main.Version=${VERSION}" ${PKG}

arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build -v -o ${BIN} -ldflags=" \
	-extldflags=-static -s -w \
	-X github.com/prometheus/common/version.Version=$(shell git describe --tags) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/prometheus/common/version.BuildDate=$(shell date --utc +%FT%T)" \
	${PKG}

test:
	@go test -v ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

clean:
	-@rm ${BIN}

.PHONY: build install test vet lint clean
