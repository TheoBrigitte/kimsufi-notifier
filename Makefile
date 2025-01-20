GROUP := github.com/TheoBrigitte
NAME := kimsufi-notifier

PKG := ${GROUP}/${NAME}

BIN := ${NAME}

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

LD_FLAGS := -s -w \
	-extldflags=-static \
	-X github.com/prometheus/common/version.Version=$(shell git describe --tags) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/prometheus/common/version.BuildDate=$(shell date --utc +%FT%T)

all: build

build:
	CGO_ENABLED=1 \
			 go build -v -o ${BIN} -ldflags="${LD_FLAGS}" ${PKG}

linux-arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ \
			 go build -v -o ${BIN} -ldflags="-extldflags=-static ${LD_FLAGS}" ${PKG}

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
	rm -rf ${BIN}

.PHONY: build linux-arm install test vet lint nancy clean
