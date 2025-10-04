package copilotcli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
)

// Client represents a Copilot CLI client
type Client struct {
	model string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Copilot CLI client
func NewClient(model string) (*Client, error) {
	// Check if copilot command is available
	if _, err := exec.LookPath("copilot"); err != nil {
		return nil, fmt.Errorf("copilot command not found in PATH: %w", err)
	}

	// Set default model if not provided
	if model == "" {
		model = "gpt-5"
	}

	return &Client{
		model: model,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Create the prompt (same as other providers for consistency)
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

	// Execute copilot command with -p flag for prompt and --model for model specification
	cmd := exec.Command("copilot", "-p", prompt, "--model", c.model)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute copilot command: %w\nstderr: %s", err, stderr.String())
	}

	response := strings.TrimSpace(stdout.String())
	if response == "" {
		return "", fmt.Errorf("empty response from copilot command")
	}

	return response, nil
}
