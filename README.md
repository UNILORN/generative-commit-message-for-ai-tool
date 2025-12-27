# generative-commit-message

[English](README.en.md) | 日本語

Gitのステージング済み変更を分析し、AIを使用して意味のあるコミットメッセージを生成するツールです。

## 特徴

- 🤖 複数のAIプロバイダーに対応（AWS Bedrock、Claude API、Gemini CLI、Copilot CLI、Claude Code）
- 🔍 環境に応じた自動プロバイダー検出
- 📝 簡潔で意味のあるコミットメッセージを生成
- ⚡ クロスプラットフォーム対応（Linux、macOS、Windows）
- 🎯 コミット粒度の評価機能

## インストール

### go install を使用（推奨）

```sh
# 最新版をインストール
go install github.com/UNILORN/generative-commit-message-for-ai-tool@latest

# 特定のバージョンをインストール（例: v1.0.0）
go install github.com/UNILORN/generative-commit-message-for-ai-tool@v1.0.0
```

バイナリは `$GOPATH/bin` にインストールされます。このディレクトリが `PATH` に含まれていることを確認してください。

### ビルド済みバイナリをダウンロード

各プラットフォーム（Linux、macOS、Windows）向けのビルド済みバイナリは [GitHub Releases](https://github.com/UNILORN/generative-commit-message-for-ai-tool/releases) からダウンロードできます。

### バージョン確認

```sh
generate-auto-commit-message version
# または
generate-auto-commit-message --version
# または
generate-auto-commit-message -v
```

## クイックスタート

このツールは利用可能な最適なAIプロバイダーを自動検出します。変更をステージして実行するだけです：

```sh
git add .
generate-auto-commit-message
```

## 使用方法

### 自動プロバイダー検出

ツールは以下の優先順位でAIプロバイダーを自動選択します：

1. **Claude API** - `ANTHROPIC_API_KEY` が設定されている場合
2. **Claude Code** - `claude` コマンドが利用可能な場合
3. **Gemini CLI** - `gemini` コマンドが利用可能な場合
4. **Copilot CLI** - `copilot` コマンドが利用可能な場合
5. **AWS Bedrock** - AWS認証情報が設定されている場合

### 手動でプロバイダーを指定

#### Gemini CLI（最も簡単）

```sh
# PATH に 'gemini' コマンドが必要
git add .
generate-auto-commit-message --provider geminicli --model "gemini-2.5-pro"
```

#### Claude Code

```sh
# PATH に 'claude' コマンドが必要
git add .
generate-auto-commit-message --provider claudecode --model "claude-sonnet-4.5"
```

#### Copilot CLI

```sh
# PATH に 'copilot' コマンドが必要
git add .
generate-auto-commit-message --provider copilotcli --model "gpt-5"
```

#### Claude API

```sh
# APIキーを設定
export ANTHROPIC_API_KEY="your-api-key"

git add .
generate-auto-commit-message --provider claude --model "claude-3-5-sonnet-20241022"
```

#### AWS Bedrock

```sh
# AWS認証情報を設定
aws sso login --profile="bedrock"
export AWS_PROFILE="bedrock"

git add .
generate-auto-commit-message --provider bedrock --model "us.anthropic.claude-3-5-sonnet-20241022-v2:0"
```

### 実行例

```sh
$ git add .
$ generate-auto-commit-message
feat: :sparkles: Gemini CLIプロバイダー対応を追加

ローカルのgeminiコマンドを統合したマルチプロバイダーアーキテクチャを実装し、自動検出機能を強化

---
コミット粒度は適切です。Gemini CLIプロバイダー機能の追加は関連性が高く、1つのコミットにまとめることが妥当です。
```

## 設定

### 環境変数

- `ANTHROPIC_API_KEY` - Claude API の直接アクセス用APIキー
- `AWS_PROFILE` - Bedrock アクセス用のAWSプロファイル
- `AWS_REGION` - Bedrock用のAWSリージョン（デフォルト: us-east-1）

### コマンドラインオプション

```sh
generate-auto-commit-message [options]

Options:
  --provider string    AIプロバイダー (bedrock, claude, geminicli, copilotcli, claudecode)
  --model string       使用するモデルID
  --region string      AWSリージョン（Bedrock用）
  --verbose            詳細な出力を有効化
  -v, --version        バージョンを表示
  version              バージョンを表示
```

## 必要要件

ステージング済みの変更があるGitリポジトリ内で実行する必要があります。

## コントリビューション

コントリビューションを歓迎します！開発のセットアップとガイドラインについては [CONTRIBUTING.md](CONTRIBUTING.md) をご覧ください。

## ライセンス

詳細は [LICENSE](LICENSE) を参照してください。
