# Makefile for generate-auto-commit-message

# Variables
BINARY_NAME=generate-auto-commit-message
MCP_SERVER_NAME=gcm-mcp-server
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
	@rm -f $(BINARY_NAME) $(MCP_SERVER_NAME)
	@go clean

# Install the application to $GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install .

# Build the MCP server
.PHONY: build-mcp
build-mcp:
	@echo "Building $(MCP_SERVER_NAME)..."
	@go build $(GOFLAGS) -o $(MCP_SERVER_NAME) ./cmd/mcp-server

# Install the MCP server to $GOPATH/bin
.PHONY: install-mcp
install-mcp: build-mcp
	@echo "Installing $(MCP_SERVER_NAME)..."
	@go install ./cmd/mcp-server

build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=linux GOARCH=amd64 go build -o dist/$(MCP_SERVER_NAME)-linux-amd64 ./cmd/mcp-server
	GOOS=linux GOARCH=arm64 go build -o dist/$(MCP_SERVER_NAME)-linux-arm64 ./cmd/mcp-server
	GOOS=darwin GOARCH=amd64 go build -o dist/$(MCP_SERVER_NAME)-darwin-amd64 ./cmd/mcp-server
	GOOS=darwin GOARCH=arm64 go build -o dist/$(MCP_SERVER_NAME)-darwin-arm64 ./cmd/mcp-server
	GOOS=windows GOARCH=amd64 go build -o dist/$(MCP_SERVER_NAME)-windows-amd64.exe ./cmd/mcp-server

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         : Build the application (default)"
	@echo "  build       : Build the application"
	@echo "  build-mcp   : Build the MCP server"
	@echo "  test        : Run tests"
	@echo "  clean       : Clean build artifacts"
	@echo "  install     : Install the application to GOPATH/bin"
	@echo "  install-mcp : Install the MCP server to GOPATH/bin"
	@echo "  build-all   : Build for multiple platforms"
	@echo "  help        : Show this help message"
