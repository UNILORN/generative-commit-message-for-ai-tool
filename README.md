# コミットメッセージ自動生成

GitリポジトリのコミットメッセージをAI（AWS Bedrock、Claude API、Gemini CLI、Copilot CLI、またはClaude Code）を使用して自動生成するGoツールです。

## 特徴

- Gitのステージング差分を読み取り
- 複数のAIプロバイダーに対応（AWS Bedrock、Claude API、Gemini CLI、Copilot CLI、Claude Code）
- 簡潔で有益なコミットメッセージを生成
- コミット粒度を評価
- クロスプラットフォーム対応（Goで構築）
- 自動プロバイダー検出機能
- Anthropic Claudeモデル対応


## インストール

### go install

go がインストールされている環境では、go install コマンドでインストールできます。
`$GOPATH/bin`配下にバイナリが配置されます。

```sh
go install github.com/UNILORN/generative-commit-message-for-bedrock.git
```

### バイナリをダウンロードして利用する

[GitHub Releases](https://github.com/UNILORN/generative-commit-message-for-bedrock/releases) から最新版をダウンロードしてください。

## 使用方法

### 選択肢1: Gemini CLI を使用する（最も簡単）

1. ローカルにgeminiコマンドがインストールされていることを確認
```sh
which gemini
# /opt/homebrew/bin/gemini などが表示される
```

2. git addして実行
```sh
$ git add .
$ generate-auto-commit-message
# 自動的にGemini CLIが選択されます（Claude APIキーがない場合）

# または明示的に指定
$ generate-auto-commit-message --provider geminicli --model "gemini-2.5-pro"
```

### 選択肢2: Claude Code を使用する

1. ローカルにclaudeコマンドがインストールされていることを確認
```sh
which claude
# /usr/local/bin/claude などが表示される
```

2. git addして実行
```sh
$ git add .
$ generate-auto-commit-message
# 自動的にClaude Codeが選択されます（Claude APIキーがなく、claudeコマンドが利用可能な場合）

# または明示的に指定
$ generate-auto-commit-message --provider claudecode --model "claude-sonnet-4.5"
```

### 選択肢3: Copilot CLI を使用する

1. ローカルにcopilotコマンドがインストールされていることを確認
```sh
which copilot
# /usr/local/bin/copilot などが表示される
```

2. git addして実行
```sh
$ git add .
$ generate-auto-commit-message
# 自動的にCopilot CLIが選択されます（Claude APIキーがなく、claudeコマンドがない場合）

# または明示的に指定
$ generate-auto-commit-message --provider copilotcli --model "gpt-5"
```

### 選択肢4: Claude API を使用する（推奨）

1. Claude API キーを環境変数に設定
```sh
export ANTHROPIC_API_KEY="your-api-key"
```

2. git addして実行
```sh
$ git add .
$ generate-auto-commit-message
# 自動的にClaude APIが選択されます

# または明示的に指定
$ generate-auto-commit-message --provider claude --model "claude-3-5-sonnet-20241022"
```

### 選択肢5: AWS Bedrock を使用する

1. AWS Bedrockを使える状態にする

Continueを利用する際に設定したBedrockのProfileを利用

```sh
aws sso login --profile="bedrock"
export AWS_PROFILE="bedrock"
```

2. 不要な環境変数はクリアして実行する

```sh
AWS_ACCESS_KEY_ID=""
AWS_SECRET_ACCESS_KEY=""
AWS_SESSION_TOKEN="" 
```

3. git addして実行

```sh
$ git add .
$ generate-auto-commit-message --provider bedrock --model "us.anthropic.claude-3-5-sonnet-20241022-v2:0"
```

### 実行例

```sh
$ git add .
$ generate-auto-commit-message
feat: :sparkles: Gemini CLIプロバイダー対応を追加する

ローカルのgeminiコマンドを利用したマルチプロバイダー構成を実装し、自動検出機能を強化

---
コミット粒度は適切です。Gemini CLIプロバイダー機能の追加は関連性が高く、1つのコミットにまとめることが妥当です。
```

### 4. 標準出力のコメントをよしなにする

おわり

## Develop

```
make help
```