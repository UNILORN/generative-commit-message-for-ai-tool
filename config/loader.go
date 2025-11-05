package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed prompt.yaml
var defaultConfigData []byte

var globalConfig *Config

// Load loads the configuration from a file or uses the default embedded config
func Load(configPath string) (*Config, error) {
	var data []byte
	var err error

	if configPath != "" {
		// Load from specified file
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Use embedded default config
		data = defaultConfigData
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// LoadDefault loads the default embedded configuration
func LoadDefault() (*Config, error) {
	return Load("")
}

// InitGlobal initializes the global configuration
func InitGlobal(configPath string) error {
	config, err := Load(configPath)
	if err != nil {
		return err
	}
	globalConfig = config
	return nil
}

// Get returns the global configuration instance
func Get() *Config {
	if globalConfig == nil {
		// Fallback to default if not initialized
		config, err := LoadDefault()
		if err != nil {
			panic(fmt.Sprintf("failed to load default config: %v", err))
		}
		globalConfig = config
	}
	return globalConfig
}

// BuildPrompt builds a prompt for commit message generation
func (c *Config) BuildPrompt(lang string, branch string, diff string) string {
	template, ok := c.PromptTemplates[lang]
	if !ok {
		// Fallback to Japanese if language not found
		template = c.PromptTemplates["japanese"]
	}

	var sb strings.Builder

	// System instruction
	sb.WriteString(template.SystemInstruction)
	sb.WriteString("\n")

	// Guidelines
	sb.WriteString("コミットメッセージは以下のガイドラインに従ってください：\n")
	for _, guideline := range template.Guidelines {
		sb.WriteString("- ")
		sb.WriteString(guideline)
		sb.WriteString("\n")
	}

	// Current branch
	sb.WriteString(fmt.Sprintf("- 現在のブランチ名は '%s' です\n", branch))
	sb.WriteString("\n")

	// Semantic Release prefixes
	sb.WriteString("- Semantic Release の記法では以下のルールに従ってください\n")
	sb.WriteString("\t- 以下は 「\"Prefixのテキスト\": 解説」の形で表記しています\n")
	prefixDescriptions := c.GetPrefixDescription(lang)
	for _, desc := range prefixDescriptions {
		sb.WriteString(desc)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Git diff
	sb.WriteString("以下が git diff です：\n")
	sb.WriteString("\n")
	sb.WriteString(diff)
	sb.WriteString("\n\n")

	// Output format
	sb.WriteString(template.OutputFormat)

	return sb.String()
}

// BuildPromptEnglish builds an English prompt for commit message generation
func (c *Config) BuildPromptEnglish(branch string, diff string) string {
	template, ok := c.PromptTemplates["english"]
	if !ok {
		// Fallback to Japanese if English not found
		return c.BuildPrompt("japanese", branch, diff)
	}

	var sb strings.Builder

	// System instruction
	sb.WriteString(template.SystemInstruction)
	sb.WriteString("\n\n")

	// Guidelines (already formatted)
	for _, guideline := range template.Guidelines {
		sb.WriteString(guideline)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Current branch
	sb.WriteString(fmt.Sprintf("Current branch: %s\n\n", branch))

	// Semantic Release prefixes
	sb.WriteString("Semantic Release Prefixes:\n")
	for _, p := range c.SemanticReleasePrefixes {
		sb.WriteString(fmt.Sprintf("- \"%s: %s\" : %s\n", p.Type, p.Emoji, p.DescriptionEN))
	}
	sb.WriteString("\n")

	// Git diff
	sb.WriteString("Git Diff:\n")
	sb.WriteString(diff)
	sb.WriteString("\n\n")

	// Output format
	sb.WriteString(template.OutputFormat)

	return sb.String()
}
