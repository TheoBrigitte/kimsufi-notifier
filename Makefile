GROUP := github.com/TheoBrigitte
NAME := kimsufi-notifier

PKG := ${GROUP}/${NAME}

BUILD_DIR := build
BIN := ${BUILD_DIR}/${NAME}

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

LD_FLAGS := -s -w \
	-X github.com/prometheus/common/version.Version=$(shell git describe --tags) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/prometheus/common/version.BuildDate=$(shell date --utc +%FT%T)

all: linux linux-arm darwin darwin-arm

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
			 go build -v -o ${BIN}_linux_amd64 -ldflags="${LD_FLAGS}" ${PKG}

linux-arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ \
			 go build -v -o ${BIN}_linux_arm64 -ldflags="-extldflags=-static ${LD_FLAGS}" ${PKG}

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
			 go build -v -o ${BIN}_darwin_amd64 -ldflags="${LD_FLAGS}" ${PKG}

darwin-arm:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
			 go build -v -o ${BIN}_darwin_arm64 -ldflags="${LD_FLAGS}" ${PKG}

install:
	GOOS=${OSTYPE} go install -v -ldflags="${LD_FLAGS}" ${PKG}

test:
	go test -v -race ./...

vet:
	go vet ${PKG_LIST}

lint:
	golint ./...
	golangci-lint run -E gosec -E goconst --timeout 10m --max-same-issues 0 --max-issues-per-linter 0

nancy:
	go list -json -m all | nancy sleuth

clean:
	rm -rf ${BUILD_DIR}

.PHONY: build install test vet lint clean
