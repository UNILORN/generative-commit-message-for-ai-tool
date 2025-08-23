package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/bedrock"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/claude"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/geminicli"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/git"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/message"
)

func main() {
	// Parse command line flags
	modelID := flag.String("model", "", "Model ID (default depends on provider)")
	region := flag.String("region", "us-east-1", "AWS region (for bedrock provider)")
	provider := flag.String("provider", "", "AI provider: 'bedrock', 'claude', or 'geminicli' (auto-detected if not specified)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("generate-auto-commit-message - Generate commit messages using AI")
		fmt.Println("\nUsage:")
		fmt.Println("  generate-auto-commit-message [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nProviders:")
		fmt.Println("  bedrock   - AWS Bedrock (requires AWS credentials)")
		fmt.Println("  claude    - Claude API (requires ANTHROPIC_API_KEY environment variable)")
		fmt.Println("  geminicli - Local Gemini CLI (requires 'gemini' command in PATH)")
		fmt.Println("\nAuto-detection:")
		fmt.Println("  If provider is not specified, it will be auto-detected based on available tools/credentials")
		os.Exit(0)
	}

	// Auto-detect provider if not specified
	if *provider == "" {
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			*provider = "claude"
		} else if _, err := exec.LookPath("gemini"); err == nil {
			*provider = "geminicli"
		} else {
			*provider = "bedrock"
		}
	}

	// Validate provider
	*provider = strings.ToLower(*provider)
	if *provider != "bedrock" && *provider != "claude" && *provider != "geminicli" {
		fmt.Fprintf(os.Stderr, "Error: Invalid provider '%s'. Must be 'bedrock', 'claude', or 'geminicli'\n", *provider)
		os.Exit(1)
	}

	// Set default model ID based on provider
	if *modelID == "" {
		if *provider == "bedrock" {
			*modelID = "anthropic.claude-3-sonnet-20240229-v1:0"
		} else if *provider == "claude" {
			*modelID = "claude-3-5-sonnet-20241022"
		} else if *provider == "geminicli" {
			*modelID = "gemini-2.5-pro"
		}
	}

	// Configure logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(nil)
	}

	// Get git diff
	diff, err := git.GetStagedDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git diff: %v\n", err)
		os.Exit(1)
	}

	if diff == "" {
		fmt.Println("No staged changes found. Please stage your changes with 'git add' first.")
		os.Exit(0)
	}

	// Get git diff
	branch, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git branch: %v\n", err)
		os.Exit(1)
	}

	if branch == "" {
		fmt.Println("No staged changes found. Please stage your changes with 'git add' before generating a commit message.")
		os.Exit(0)
	}

	// Initialize AI client based on provider
	var aiClient client.AIClient

	if *provider == "bedrock" {
		aiClient, err = bedrock.NewClient(*region, *modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing AWS Bedrock client: %v\n", err)
			os.Exit(1)
		}
	} else if *provider == "claude" {
		aiClient, err = claude.NewClient(*modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Claude API client: %v\n", err)
			os.Exit(1)
		}
	} else if *provider == "geminicli" {
		aiClient, err = geminicli.NewClient(*modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Gemini CLI client: %v\n", err)
			os.Exit(1)
		}
	}

	// Generate commit message
	commitMsg, err := message.Generate(aiClient, diff, branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// デバッグ情報の出力
	if *verbose {
		fmt.Println("=== Debug Information ===")
		fmt.Printf("Provider: %s\n", *provider)
		fmt.Printf("Model ID: %s\n", *modelID)
		if *provider == "bedrock" {
			fmt.Printf("Region: %s\n", *region)
		}
		fmt.Printf("Diff size: %d bytes\n", len(diff))
		fmt.Println("========================")
	}

	// Print the generated commit message
	fmt.Println(commitMsg)
}
