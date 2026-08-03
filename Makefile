PROJECT_NAME := kimsufi-notifier

# Directories
BUILD_DIR := build

# Build informations
BUILD_USER ?= $(shell whoami)@$(shell hostname)
BUILD_DATE ?= $(shell date --utc +%FT%T)
CC ?= $(shell go env CC)
CXX ?= $(shell go env CXX)
GOARCH ?= $(shell go env GOARCH)
GOOS ?= $(shell go env GOOS)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
REVISION := $(shell git rev-parse HEAD)
VERSION := $(shell git describe --always --tags)

# Build arguments
GO_BIN := $(BUILD_DIR)/$(PROJECT_NAME).$(GOARCH)
LD_FLAGS := -s -w \
	-extldflags=-static \
	-X github.com/prometheus/common/version.Version=$(VERSION) \
	-X github.com/prometheus/common/version.Revision=$(REVISION) \
	-X github.com/prometheus/common/version.Branch=$(BRANCH) \
	-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) \
	-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)

# Supported architectures and OSes
ARCHS = amd64 arm64
OSES = linux

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# Default target
.DEFAULT_GOAL := build

##@ Development environment

.PHONY: setup
setup: ## Setup the development environment
	pre-commit install

##@ Build

.PHONY: build
build: ## Build the Go binary
	@printf "$(CYAN)Building Go binary...$(RESET)\n"
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 \
			 CC=$(CC) CXX=$(CXX) \
			 go build -v -o $(GO_BIN) -ldflags="${LD_FLAGS}"	.
	@printf "$(GREEN)Build completed. Output is in $(GO_BIN)$(RESET)\n"

.PHONY: build-amd64
build-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

.PHONY: build-arm64
build-arm64: ## Build the Go binary for ARM64 architecture
	@$(MAKE) build GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++

.PHONY: build-all
build-all: ## Build the Go binary for all supported OSes and architectures
build-all: build-amd64 build-arm64

.PHONY: install
install: build ## Install the binary to ~/.local/bin
	@printf "$(CYAN)Installing binary to ~/.local/bin...$(RESET)\n"
	@mkdir -p ~/.local/bin
	cp $(GO_BIN) ~/.local/bin/$(PROJECT_NAME)
	chmod +x ~/.local/bin/$(PROJECT_NAME)
	@printf "$(GREEN)Binary installed to ~/.local/bin/$(PROJECT_NAME)$(RESET)\n"

.PHONY: clean
clean: ## Clean build artifacts
	@printf "$(CYAN)Cleaning build artifacts...$(RESET)\n"
	rm -rf $(BUILD_DIR)
	@printf "$(GREEN)Cleanup completed$(RESET)\n"

##@ Testing

.PHONY: test
test: ## Run the complete test suite, optionally filtered by run_pattern or bench_pattern
	@printf "$(CYAN)Running tests...$(RESET)\n"
	go test -v -race -run="$(run_pattern)" -bench="$(bench_pattern)" -benchmem ./...
	@printf "$(GREEN)Tests completed successfully$(RESET)\n"

##@ Code Quality

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint for comprehensive code analysis (requires CGO environment)
	@printf "$(CYAN)Running golangci-lint...$(RESET)\n"
	golangci-lint run -E gosec -E goconst --timeout 10m --max-same-issues 0 --max-issues-per-linter 0 ./...
	@printf "$(GREEN)Linting completed$(RESET)\n"

.PHONY: golint
golint: ## Run golint for basic linting
	@printf "$(CYAN)Running golint...$(RESET)\n"
	golint ./...
	@printf "$(GREEN)Golint completed$(RESET)\n"

.PHONY: vet
vet: ## Run go vet for static analysis
	@printf "$(CYAN)Running go vet...$(RESET)\n"
	go vet ./...
	@printf "$(GREEN)Static analysis completed$(RESET)\n"

.PHONY: fmt
fmt: ## Check code formatting
	@printf "$(CYAN)Checking code formatting...$(RESET)\n"
	gofmt -d .

.PHONY: lint
lint: fmt vet golint golangci-lint ## Run all linters

##@ Security

.PHONY: nancy
nancy: ## Run Nancy vulnerability scan
	@printf "$(CYAN)Running nancy vulnerability scan...$(RESET)\n"
	go list -json -m all | nancy sleuth
	@printf "$(GREEN)Nancy scan completed$(RESET)\n"

.PHONY: security
security: nancy ## Run all security scans

##@ Help

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(CYAN)Usage:$(RESET)\n  make $(YELLOW)<target>$(RESET)\n"} /^[a-zA-Z_0-9-]+.*?##/ { printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(CYAN)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
