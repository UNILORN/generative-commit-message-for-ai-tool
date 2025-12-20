# Contributing to generative-commit-message

Thank you for your interest in contributing to this project!

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Building from Source

Clone the repository:

```sh
git clone https://github.com/UNILORN/generative-commit-message-for-bedrock.git
cd generative-commit-message-for-bedrock
```

Build the project:

```sh
make build
# or
go build -o generate-auto-commit-message .
```

### Development Commands

The project uses a Makefile for common development tasks:

```sh
make help          # Show all available targets
make build         # Build the binary
make test          # Run all tests with verbose output
make install       # Build and install to $GOPATH/bin
make clean         # Remove build artifacts
make build-all     # Cross-compile for multiple platforms
```

You can also use standard Go commands:

```sh
go run .           # Run the application directly
go test ./...      # Run tests for all packages
go build .         # Build the binary
```

## Project Architecture

The codebase follows a clean modular architecture:

### Core Packages

- `client/` - Abstract interface for AI providers
- `bedrock/` - AWS Bedrock client implementation
- `claude/` - Claude API direct client implementation
- `geminicli/` - Local Gemini CLI client implementation
- `copilotcli/` - GitHub Copilot CLI client implementation
- `claudecode/` - Claude Code CLI client implementation
- `git/` - Git operations (staged diffs, branch detection, file status)
- `message/` - Commit message generation logic
- `main.go` - CLI entry point with flag parsing

### Data Flow

1. CLI parses flags (provider, model ID, region, verbose mode)
2. Auto-detects provider based on environment variables and available tools
3. `git` package extracts staged changes and current branch
4. Appropriate AI client is initialized through the `client` interface
5. `message` package combines git context and calls AI client
6. Generated message is printed to stdout

## Testing

The project includes unit tests in the `git/` package. Run tests with:

```sh
make test
# or
go test -v ./...
```

When adding new features, please include appropriate tests.

## Code Style

- Follow standard Go conventions and formatting (`gofmt`, `go vet`)
- Write clear, concise commit messages
- Keep functions focused and modular
- Add comments for complex logic

## Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Make your changes and commit them with clear messages
4. Push to your fork
5. Open a Pull Request

## Questions?

If you have questions about contributing, feel free to open an issue for discussion.
