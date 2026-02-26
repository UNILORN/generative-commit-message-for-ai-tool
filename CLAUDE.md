# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI tool that generates commit messages using AI providers. The tool analyzes Git staged changes and generates appropriate commit messages through AWS Bedrock, Claude API directly, or local Gemini CLI.

## Build and Development Commands

Use the Makefile for all build operations:

- `make build` - Build the binary 
- `make test` - Run all tests with verbose output
- `make install` - Build and install to $GOPATH/bin
- `make clean` - Remove build artifacts
- `make release` - Cross-compile for multiple platforms (Linux, macOS, Windows)
- `make help` - Show all available targets

For development, use standard Go commands:
- `go run .` - Run the application directly
- `go test ./...` - Run tests for all packages

## Architecture

The codebase follows a clean modular architecture with three main packages:

### Core Packages
- `client/` - Abstract interface for AI providers
- `bedrock/` - AWS Bedrock client implementation
- `claude/` - Claude API direct client implementation
- `geminicli/` - Local Gemini CLI client implementation (executes `gemini` command)
- `git/` - Git operations (staged diffs, branch detection, file status)  
- `message/` - Commit message generation logic that orchestrates AI clients and git packages
- `main.go` - CLI entry point with flag parsing and orchestration

### Data Flow
1. CLI parses flags (provider, model ID, region, verbose mode)
2. Auto-detects provider based on environment variables and available tools if not specified
3. `git` package extracts staged changes and current branch
4. Appropriate AI client (`bedrock`, `claude`, or `geminicli`) is initialized
5. `message` package combines git context and calls AI client through interface
6. Generated message is printed to stdout

## Usage Requirements

### For Gemini CLI (Simplest)
- Requires `gemini` command available in PATH
- Default model: `gemini-2.5-pro`

### For Claude API (Recommended for API access)
- Requires `ANTHROPIC_API_KEY` environment variable
- Default model: `claude-sonnet-4-6`

### For AWS Bedrock
- Requires AWS credentials configured (via AWS Profile, environment variables, or IAM roles)
- Default model: `anthropic.claude-sonnet-4-5-20250929-v1:0` in `us-east-1` region

Must be run from within a Git repository with staged changes.

## Testing

The project includes unit tests in the `git/` package. Run tests with `make test` for verbose output including test names and results.