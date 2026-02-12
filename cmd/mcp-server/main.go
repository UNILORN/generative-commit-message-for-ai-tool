package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/UNILORN/generative-commit-message-for-ai-tool/mcp"
)

func main() {
	provider := flag.String("provider", "", "Default AI provider (bedrock, claude, geminicli, copilotcli, claudecode)")
	modelID := flag.String("model", "", "Default model ID")
	region := flag.String("region", "us-east-1", "AWS region (for bedrock provider)")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	server, err := mcp.NewServer(*provider, *modelID, *region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating MCP server: %v\n", err)
		os.Exit(1)
	}

	if err := server.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("gcm-mcp-server - MCP server for generating commit messages using AI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gcm-mcp-server [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -provider string")
	fmt.Println("        Default AI provider (bedrock, claude, geminicli, copilotcli, claudecode)")
	fmt.Println("  -model string")
	fmt.Println("        Default model ID")
	fmt.Println("  -region string")
	fmt.Println("        AWS region (for bedrock provider) (default \"us-east-1\")")
	fmt.Println("  -help")
	fmt.Println("        Show help")
	fmt.Println()
	fmt.Println("Available Tools:")
	fmt.Println("  get_staged_diff        - Get the diff of all staged changes")
	fmt.Println("  get_staged_files       - Get the list of staged files with status")
	fmt.Println("  generate_commit_message - Generate a commit message using AI")
	fmt.Println("  commit                 - Create a git commit with a specified message")
	fmt.Println("  generate_and_commit    - Generate a message and create a commit")
	fmt.Println()
	fmt.Println("Example Claude Code configuration (~/.claude/claude_desktop_config.json):")
	fmt.Println(`  {
    "mcpServers": {
      "gcm": {
        "command": "gcm-mcp-server",
        "args": ["-provider", "claude"]
      }
    }
  }`)
}
