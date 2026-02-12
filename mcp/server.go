package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-ai-tool/bedrock"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/claude"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/claudecode"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/client"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/config"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/copilotcli"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/geminicli"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/git"
	"github.com/UNILORN/generative-commit-message-for-ai-tool/message"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP server for commit message generation
type Server struct {
	mcpServer *server.MCPServer
	provider  string
	modelID   string
	region    string
}

// NewServer creates a new MCP server instance
func NewServer(provider, modelID, region string) (*Server, error) {
	// Initialize config
	if err := config.InitGlobal(""); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	s := &Server{
		provider: provider,
		modelID:  modelID,
		region:   region,
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"gcm",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register tools
	s.registerTools(mcpServer)

	s.mcpServer = mcpServer
	return s, nil
}

// registerTools registers all available tools
func (s *Server) registerTools(mcpServer *server.MCPServer) {
	// Tool: get_staged_diff
	mcpServer.AddTool(
		mcp.NewTool("get_staged_diff",
			mcp.WithDescription("Get the diff of all staged changes in the current git repository"),
		),
		s.handleGetStagedDiff,
	)

	// Tool: get_staged_files
	mcpServer.AddTool(
		mcp.NewTool("get_staged_files",
			mcp.WithDescription("Get the list of staged files with their status (Added/Modified/Deleted)"),
		),
		s.handleGetStagedFiles,
	)

	// Tool: generate_commit_message
	mcpServer.AddTool(
		mcp.NewTool("generate_commit_message",
			mcp.WithDescription("Generate a commit message for staged changes using AI"),
			mcp.WithString("provider",
				mcp.Description("AI provider to use (bedrock, claude, geminicli, copilotcli, claudecode). If not specified, auto-detected."),
			),
			mcp.WithString("model",
				mcp.Description("Model ID to use. If not specified, uses default for the provider."),
			),
		),
		s.handleGenerateCommitMessage,
	)

	// Tool: commit
	mcpServer.AddTool(
		mcp.NewTool("commit",
			mcp.WithDescription("Create a git commit with the specified message"),
			mcp.WithString("message",
				mcp.Required(),
				mcp.Description("The commit message to use"),
			),
		),
		s.handleCommit,
	)

	// Tool: generate_and_commit
	mcpServer.AddTool(
		mcp.NewTool("generate_and_commit",
			mcp.WithDescription("Generate a commit message using AI and create a commit with it"),
			mcp.WithString("provider",
				mcp.Description("AI provider to use (bedrock, claude, geminicli, copilotcli, claudecode). If not specified, auto-detected."),
			),
			mcp.WithString("model",
				mcp.Description("Model ID to use. If not specified, uses default for the provider."),
			),
		),
		s.handleGenerateAndCommit,
	)
}

// handleGetStagedDiff handles the get_staged_diff tool call
func (s *Server) handleGetStagedDiff(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	diff, err := git.GetStagedDiff()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged diff: %v", err)), nil
	}

	if diff == "" {
		return mcp.NewToolResultText("No staged changes found. Please stage your changes with 'git add' first."), nil
	}

	return mcp.NewToolResultText(diff), nil
}

// handleGetStagedFiles handles the get_staged_files tool call
func (s *Server) handleGetStagedFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filesWithStatus, err := git.GetStagedFilesWithStatus()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged files: %v", err)), nil
	}

	if filesWithStatus == "" {
		return mcp.NewToolResultText("No staged files found."), nil
	}

	return mcp.NewToolResultText(filesWithStatus), nil
}

// handleGenerateCommitMessage handles the generate_commit_message tool call
func (s *Server) handleGenerateCommitMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get parameters
	provider := s.getStringParam(request, "provider", s.provider)
	modelID := s.getStringParam(request, "model", s.modelID)

	// Get staged diff
	diff, err := git.GetStagedDiff()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged diff: %v", err)), nil
	}

	if diff == "" {
		return mcp.NewToolResultError("No staged changes found. Please stage your changes with 'git add' first."), nil
	}

	// Get current branch
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current branch: %v", err)), nil
	}

	// Create AI client
	aiClient, err := s.createAIClient(provider, modelID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create AI client: %v", err)), nil
	}

	// Generate commit message
	commitMsg, err := message.Generate(aiClient, diff, branch)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to generate commit message: %v", err)), nil
	}

	return mcp.NewToolResultText(commitMsg), nil
}

// handleCommit handles the commit tool call
func (s *Server) handleCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get commit message from parameters
	msg := s.getStringParam(request, "message", "")
	if msg == "" {
		return mcp.NewToolResultError("Commit message is required"), nil
	}

	// Verify staged changes exist before committing (prevent race condition)
	diff, err := git.GetStagedDiff()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to verify staged changes: %v", err)), nil
	}
	if diff == "" {
		return mcp.NewToolResultError("No staged changes found. Please stage your changes with 'git add' first."), nil
	}

	// Execute git commit
	cmd := exec.Command("git", "commit", "-m", msg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to commit: %v\n%s", err, string(output))), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Commit created successfully:\n%s", string(output))), nil
}

// handleGenerateAndCommit handles the generate_and_commit tool call
func (s *Server) handleGenerateAndCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get parameters
	provider := s.getStringParam(request, "provider", s.provider)
	modelID := s.getStringParam(request, "model", s.modelID)

	// Get staged diff
	diff, err := git.GetStagedDiff()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged diff: %v", err)), nil
	}

	if diff == "" {
		return mcp.NewToolResultError("No staged changes found. Please stage your changes with 'git add' first."), nil
	}

	// Get current branch
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current branch: %v", err)), nil
	}

	// Create AI client
	aiClient, err := s.createAIClient(provider, modelID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create AI client: %v", err)), nil
	}

	// Generate commit message
	commitMsg, err := message.Generate(aiClient, diff, branch)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to generate commit message: %v", err)), nil
	}

	// Execute git commit
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to commit: %v\n%s", err, string(output))), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Commit created with message:\n%s\n\nOutput:\n%s", commitMsg, string(output))), nil
}

// getStringParam gets a string parameter from the request
func (s *Server) getStringParam(request mcp.CallToolRequest, name, defaultValue string) string {
	if request.Params.Arguments == nil {
		return defaultValue
	}

	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return defaultValue
	}

	if val, ok := args[name]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return strVal
		}
	}
	return defaultValue
}

// createAIClient creates an AI client based on the provider
func (s *Server) createAIClient(provider, modelID string) (client.AIClient, error) {
	// Auto-detect provider if not specified
	if provider == "" {
		provider = s.autoDetectProvider()
	}

	// Set default model ID if not specified
	if modelID == "" {
		modelID = s.getDefaultModelID(provider)
	}

	provider = strings.ToLower(provider)

	switch provider {
	case "bedrock":
		region := s.region
		if region == "" {
			region = "us-east-1"
		}
		return bedrock.NewClient(region, modelID)
	case "claude":
		return claude.NewClient(modelID)
	case "geminicli":
		return geminicli.NewClient(modelID)
	case "copilotcli":
		return copilotcli.NewClient(modelID)
	case "claudecode":
		return claudecode.NewClient(modelID)
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

// autoDetectProvider detects the best available provider
func (s *Server) autoDetectProvider() string {
	if s.provider != "" {
		return s.provider
	}

	// Check environment variables first (same logic as main.go)
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return "claude"
	}

	// Check available CLI tools
	if _, err := exec.LookPath("claude"); err == nil {
		return "claudecode"
	}
	if _, err := exec.LookPath("copilot"); err == nil {
		return "copilotcli"
	}
	if _, err := exec.LookPath("gemini"); err == nil {
		return "geminicli"
	}
	return "bedrock"
}

// getDefaultModelID returns the default model ID for a provider
func (s *Server) getDefaultModelID(provider string) string {
	switch strings.ToLower(provider) {
	case "bedrock":
		return "anthropic.claude-3-sonnet-20240229-v1:0"
	case "claude":
		return "claude-3-5-sonnet-20241022"
	case "geminicli":
		return "gemini-2.5-pro"
	case "copilotcli", "claudecode":
		return "claude-sonnet-4.5"
	default:
		return ""
	}
}

// ServeStdio starts the MCP server using stdio transport
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.mcpServer)
}
