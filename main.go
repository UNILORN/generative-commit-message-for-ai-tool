package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-ai-tool/bedrock"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/claude"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/claudecode"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/client"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/config"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/copilotcli"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/copilotsdk"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/geminicli"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/git"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/message"
)

func main() {
	// Check for subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			runInit(os.Args[2:])
			return
		case "version", "--version", "-v":
			printVersion()
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	// No subcommand, run default behavior (generate commit message)
	runGenerate(os.Args[1:])
}

func printHelp() {
	fmt.Println("generate-auto-commit-message - Generate commit messages using AI")
	fmt.Println("\nUsage:")
	fmt.Println("  generate-auto-commit-message [options]          Generate commit message")
	fmt.Println("  generate-auto-commit-message init [options]     Initialize config file")
	fmt.Println("  generate-auto-commit-message version            Show version")
	fmt.Println("  generate-auto-commit-message help               Show this help")
	fmt.Println("\nGenerate Options:")
	generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	generateFlags.String("model", "", "Model ID (default depends on provider)")
	generateFlags.String("region", "us-east-1", "AWS region (for bedrock provider)")
	generateFlags.String("provider", "", "AI provider: 'bedrock', 'claude', 'geminicli', 'copilotcli', 'copilotsdk', or 'claudecode' (auto-detected if not specified)")
	generateFlags.String("config", "", "Path to config file (uses embedded default if not specified)")
	generateFlags.Bool("verbose", false, "Enable verbose output")
	generateFlags.PrintDefaults()
	fmt.Println("\nInit Options:")
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	initFlags.String("f", "./prompt.yaml", "Output file path (short form)")
	initFlags.String("file", "./prompt.yaml", "Output file path (long form)")
	initFlags.Bool("force", false, "Overwrite existing file")
	initFlags.PrintDefaults()
	fmt.Println("\nProviders:")
	fmt.Println("  bedrock    - AWS Bedrock (requires AWS credentials)")
	fmt.Println("  claude     - Claude API (requires ANTHROPIC_API_KEY environment variable)")
	fmt.Println("  geminicli  - Local Gemini CLI (requires 'gemini' command in PATH)")
	fmt.Println("  copilotcli - Copilot CLI direct execution (runs 'copilot' command)")
	fmt.Println("  copilotsdk - Copilot SDK programmatic access (uses SDK for session management)")
	fmt.Println("  claudecode - Claude Code CLI (requires 'claude' command in PATH)")
	fmt.Println("\nAuto-detection:")
	fmt.Println("  If provider is not specified, it will be auto-detected based on available tools/credentials")
	fmt.Println("\nExamples:")
	fmt.Println("  # Generate commit message with auto-detected provider")
	fmt.Println("  generate-auto-commit-message")
	fmt.Println()
	fmt.Println("  # Generate with specific provider")
	fmt.Println("  generate-auto-commit-message --provider=claude")
	fmt.Println()
	fmt.Println("  # Initialize config file")
	fmt.Println("  generate-auto-commit-message init")
	fmt.Println()
	fmt.Println("  # Initialize with custom path")
	fmt.Println("  generate-auto-commit-message init -f ./my-prompt.yaml")
	fmt.Println()
	fmt.Println("  # Overwrite existing config")
	fmt.Println("  generate-auto-commit-message init --force")
}

func runInit(args []string) {
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	fileShort := initFlags.String("f", "./prompt.yaml", "Output file path")
	fileLong := initFlags.String("file", "", "Output file path")
	force := initFlags.Bool("force", false, "Overwrite existing file")
	initFlags.Parse(args)

	// Use --file if specified, otherwise use -f
	outputPath := *fileShort
	if *fileLong != "" {
		outputPath = *fileLong
	}

	// Write default config
	if err := config.WriteDefaultConfig(outputPath, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Config file created: %s\n", outputPath)
	fmt.Println()
	fmt.Println("You can now:")
	fmt.Println("  1. Edit the config file to customize prompts")
	fmt.Println("  2. Use it with: generate-auto-commit-message --config=" + outputPath)
}

func runGenerate(args []string) {
	// Parse command line flags
	generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	modelID := generateFlags.String("model", "", "Model ID (default depends on provider)")
	region := generateFlags.String("region", "us-east-1", "AWS region (for bedrock provider)")
	provider := generateFlags.String("provider", "", "AI provider: 'bedrock', 'claude', 'geminicli', 'copilotcli', 'copilotsdk', or 'claudecode' (auto-detected if not specified)")
	configPath := generateFlags.String("config", "", "Path to config file (uses embedded default if not specified)")
	verbose := generateFlags.Bool("verbose", false, "Enable verbose output")
	help := generateFlags.Bool("help", false, "Show help")
	generateFlags.Parse(args)

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Auto-detect provider if not specified
	if *provider == "" {
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			*provider = "claude"
		} else if _, err := exec.LookPath("claude"); err == nil {
			*provider = "claudecode"
		} else if _, err := exec.LookPath("copilot"); err == nil {
			*provider = "copilotcli"
		} else if _, err := exec.LookPath("gemini"); err == nil {
			*provider = "geminicli"
		} else {
			*provider = "bedrock"
		}
	}

	// Validate provider
	*provider = strings.ToLower(*provider)
	if *provider != "bedrock" && *provider != "claude" && *provider != "geminicli" && *provider != "copilotcli" && *provider != "copilotsdk" && *provider != "claudecode" {
		fmt.Fprintf(os.Stderr, "Error: Invalid provider '%s'. Must be 'bedrock', 'claude', 'geminicli', 'copilotcli', 'copilotsdk', or 'claudecode'\n", *provider)
		os.Exit(1)
	}

	// Set default model ID based on provider
	if *modelID == "" {
		if *provider == "bedrock" {
			*modelID = "anthropic.claude-sonnet-4-5-20250929-v1:0"
		} else if *provider == "claude" {
			*modelID = "claude-sonnet-4-6"
		} else if *provider == "geminicli" {
			*modelID = "gemini-2.5-pro"
		} else if *provider == "copilotcli" {
			*modelID = "claude-sonnet-4.5"
		} else if *provider == "copilotsdk" {
			*modelID = "gpt-4o"
		} else if *provider == "claudecode" {
			*modelID = "claude-sonnet-4.5"
		}
	}

	// Initialize config
	if err := config.InitGlobal(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
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
	} else if *provider == "copilotcli" {
		aiClient, err = copilotcli.NewClient(*modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Copilot CLI client: %v\n", err)
			os.Exit(1)
		}
	} else if *provider == "copilotsdk" {
		aiClient, err = copilotsdk.NewClient(*modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Copilot SDK client: %v\n", err)
			os.Exit(1)
		}
	} else if *provider == "claudecode" {
		aiClient, err = claudecode.NewClient(*modelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Claude Code client: %v\n", err)
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
