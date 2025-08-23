package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/bedrock"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/git"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/message"
)

func main() {
	// Parse command line flags
	modelID := flag.String("model", "anthropic.claude-3-sonnet-20240229-v1:0", "AWS Bedrock model ID")
	region := flag.String("region", "us-east-1", "AWS region")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("generate-auto-commit-message - Generate commit messages using AWS Bedrock")
		fmt.Println("\nUsage:")
		fmt.Println("  generate-auto-commit-message [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(0)
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

	// Initialize Bedrock client
	client, err := bedrock.NewClient(*region, *modelID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AWS Bedrock client: %v\n", err)
		os.Exit(1)
	}

	// Generate commit message
	commitMsg, err := message.Generate(client, diff, branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// デバッグ情報の出力
	if *verbose {
		fmt.Println("=== Debug Information ===")
		fmt.Printf("Model ID: %s\n", *modelID)
		fmt.Printf("Region: %s\n", *region)
		fmt.Printf("Diff size: %d bytes\n", len(diff))
		fmt.Println("========================")
	}

	// Print the generated commit message
	fmt.Println(commitMsg)
}
