# コミットメッセージ自動生成

GitリポジトリのコミットメッセージをAI（AWS Bedrock、Claude API、またはローカルGemini CLI）を使用して自動生成するGoツールです。

## 特徴

- Gitのステージング差分を読み取り
- 複数のAIプロバイダーに対応（AWS Bedrock、Claude API、Gemini CLI）
- 簡潔で有益なコミットメッセージを生成
- コミット粒度を評価
- クロスプラットフォーム対応（Goで構築）
- 自動プロバイダー検出機能
- Anthropic Claudeモデル対応


## インストール

- [go install](#go-install)
- [binary install（mac arm64）](#バイナリをダウンロードして利用する)

### go install

go がインストールされている環境では、go install コマンドでインストールできます。
`$GOPATH/bin`配下にバイナリが配置されます。

```sh
go install github.com/UNILORN/generative-commit-message-for-bedrock.git
```

### バイナリをダウンロードして利用する

※Mac arm64のみ対応

#### 1. バイナリをダウンロード

[generate-auto-commit-message](/uploads/a3724435d66999c7c98250feca8af38b/generate-auto-commit-message)

#### 2. 配置,権限付与

```sh
sudo mv ~/Downloads/generate-auto-commit-message /usr/local/bin/generate-auto-commit-message
chmod +x /usr/local/bin/generate-auto-commit-message
```

#### 3. 一度利用出来る状態にし、Macのセキュリティ許可を実行する

- お好きなGitRepositoryへ移動

```sh
generate-auto-commit-message --model "us.anthropic.claude-3-7-sonnet-20250219-v1:0"
```

- 警告がでるので「完了」を押下
- システム設定 -> プライバシーとセキュリティ -> 最下部で「このまま許可」を押下
- 再度コマンド実行
- このまま許可

```sh
generate-auto-commit-message --model "us.anthropic.claude-3-7-sonnet-20250219-v1:0"

No staged changes found. Please stage your changes with 'git add' first.
```

このような出力になっていれば完了

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

### 選択肢2: Claude API を使用する（推奨）

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

### 選択肢3: AWS Bedrock を使用する

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