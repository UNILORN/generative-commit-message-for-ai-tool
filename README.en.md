# generative-commit-message

English | [日本語](README.md)

AI-powered commit message generator that analyzes your Git staged changes and generates meaningful commit messages.

## Features

- 🤖 Multiple AI provider support (AWS Bedrock, Claude API, Gemini CLI, Copilot CLI, Claude Code)
- 🔍 Automatic provider detection based on environment
- 📝 Generates concise and meaningful commit messages
- ⚡ Cross-platform support (Linux, macOS, Windows)
- 🎯 Evaluates commit granularity

## Installation

### Using go install (Recommended)

```sh
# Install latest version
go install github.com/UNILORN/generative-commit-message-for-ai-tool@latest

# Install specific version (e.g., v1.0.0)
go install github.com/UNILORN/generative-commit-message-for-ai-tool@v1.0.0
```

The binary will be installed to `$GOPATH/bin`. Make sure this directory is in your `PATH`.

#### Rename to Shorter Name (Optional)

If the command name feels too long, you can rename it after installation:

```sh
# Rename to gcm
mv $(go env GOPATH)/bin/generate-auto-commit-message $(go env GOPATH)/bin/gcm

# Usage example
git add .
gcm
```

### Download Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/UNILORN/generative-commit-message-for-ai-tool/releases) for your platform (Linux, macOS, Windows).

### Check Version

```sh
generate-auto-commit-message version
# or
generate-auto-commit-message --version
# or
generate-auto-commit-message -v
```

## Quick Start

The tool automatically detects the best available AI provider. Simply stage your changes and run:

```sh
git add .
generate-auto-commit-message
```

## Usage

### Automatic Provider Detection

The tool automatically selects an AI provider in the following priority order:

1. **Claude API** - if `ANTHROPIC_API_KEY` is set
2. **Claude Code** - if `claude` command is available
3. **Gemini CLI** - if `gemini` command is available
4. **Copilot CLI** - if `copilot` command is available
5. **AWS Bedrock** - if AWS credentials are configured

### Manual Provider Selection

#### Gemini CLI (Easiest)

```sh
# Requires 'gemini' command in PATH
git add .
generate-auto-commit-message --provider geminicli --model "gemini-2.5-pro"
```

#### Claude Code

```sh
# Requires 'claude' command in PATH
git add .
generate-auto-commit-message --provider claudecode --model "claude-sonnet-4.5"
```

#### Copilot CLI

```sh
# Requires 'copilot' command in PATH
git add .
generate-auto-commit-message --provider copilotcli --model "gpt-5"
```

#### Claude API

```sh
# Set API key
export ANTHROPIC_API_KEY="your-api-key"

git add .
generate-auto-commit-message --provider claude --model "claude-3-5-sonnet-20241022"
```

#### AWS Bedrock

```sh
# Configure AWS credentials
aws sso login --profile="bedrock"
export AWS_PROFILE="bedrock"

git add .
generate-auto-commit-message --provider bedrock --model "us.anthropic.claude-3-5-sonnet-20241022-v2:0"
```

### Example Output

```sh
$ git add .
$ generate-auto-commit-message
feat: :sparkles: Add Gemini CLI provider support

Implement multi-provider architecture with local gemini command integration and enhanced auto-detection

---
Commit granularity is appropriate. The Gemini CLI provider feature addition is highly related and suitable for a single commit.
```

## Configuration

### Environment Variables

- `ANTHROPIC_API_KEY` - Claude API key for direct API access
- `AWS_PROFILE` - AWS profile for Bedrock access
- `AWS_REGION` - AWS region for Bedrock (default: us-east-1)

### Command-line Options

```sh
generate-auto-commit-message [options]

Options:
  --provider string    AI provider (bedrock, claude, geminicli, copilotcli, claudecode)
  --model string       Model ID to use
  --region string      AWS region (for Bedrock)
  --verbose            Enable verbose output
  -v, --version        Show version
  version              Show version
```

## Requirements

Must be run from within a Git repository with staged changes.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

See [LICENSE](LICENSE) for details.
