# Makefile for generate-auto-commit-message

# Variables
BINARY_NAME=generate-auto-commit-message
GOFLAGS=-ldflags="-s -w" # Strip debug information to reduce binary size
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(GOFLAGS) -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean

# Install the application to $GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install .

# Build for multiple platforms
.PHONY: release
release:
	@echo "Building for multiple platforms..."
	@mkdir -p ./bin
	
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o ./bin/$(BINARY_NAME)-linux-amd64 .
	
	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o ./bin/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o ./bin/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o ./bin/$(BINARY_NAME)-windows-amd64.exe .

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all      : Build the application (default)"
	@echo "  build    : Build the application"
	@echo "  test     : Run tests"
	@echo "  clean    : Clean build artifacts"
	@echo "  install  : Install the application to GOPATH/bin"
	@echo "  release  : Build for multiple platforms"
	@echo "  help     : Show this help message"
