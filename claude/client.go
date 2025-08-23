package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
)

// Client represents a Claude API client
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Claude API client
func NewClient(model string) (*Client, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.anthropic.com",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// ClaudeMessage represents a message in the Claude API format
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeRequest represents a request to the Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeResponseContent represents content in the Claude API response
type ClaudeResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeUsage represents usage information in the Claude API response
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ClaudeResponse represents a response from the Claude API
type ClaudeResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Role       string                  `json:"role"`
	Model      string                  `json:"model"`
	Content    []ClaudeResponseContent `json:"content"`
	StopReason string                  `json:"stop_reason"`
	Usage      ClaudeUsage             `json:"usage"`
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Create the prompt (same as bedrock for consistency)
	prompt := fmt.Sprintf(`あなたは提供された diff に基づいて、簡潔で有益な git コミットメッセージを生成する役立つアシスタントです。
コミットメッセージは以下のガイドラインに従ってください：
- 短い要約行で始める（50〜72文字）
- 命令形を使用する（例：「機能を追加する」であって「機能を追加した」ではない）
- 必要に応じて空白行の後に詳細な説明を含める
- 「どのように」よりも「なぜ」と「何を」に焦点を当てる
- Semantic Release の記法でPrefixをつける
- ブランチ名に数字が含まれていた場合、'feat: 本文 #1234' のように記入してください。
- 現在のブランチ名は '%s' です

- Semantic Release の記法では以下のルールに従ってください
	- 以下は 「"Prefixのテキスト": 解説」の形で表記しています
	- "feat: :sparkles:" : 新機能追加
	- "fix: :bug:" : バグ修正
	- "refactor: :hammer:" : レビューや仕様変更によるコード修正
	- "test: :white_check_mark:" : テストの追加や既存テストの修正
	- "docs: :memo:" : ドキュメントのみの変更
	- "config: :wrench:" : 設定ファイルの追加・更新
	- "lint: :rotating_light:" : リンターの警告を修正
	- "ci: :construction_worker:" : CIの追加・修正
	- "remove: :wastebasket:" : 削除
	- "improve: :zap:" : パフォーマンス改善のためのコード修正
	- "try: :bulb:" : 検証や試行錯誤のコード修正
	- "wip: :construction:" : WIP
	- "update: :up:" : ライブラリのアップデート
	- "release: :rocket:" : リリース
	- "merge: :twisted_rightwards_arrows:" : マージ・ブランチ統合

以下が git diff です：

%s

日本語でコミットメッセージを生成してください。
その際、コードブロック文字は不要です。コミットメッセージのみを出力してください。`, branch, diff)

	// Create the request
	request := ClaudeRequest{
		Model:     c.model,
		MaxTokens: 10000,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal the request to JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/v1/messages", bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers according to Claude API documentation
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the response
	var response ClaudeResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the commit message
	if len(response.Content) > 0 && len(response.Content[0].Text) > 0 {
		return response.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}