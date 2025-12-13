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
	// Normalize language code
	normalizedLang := normalizeLangCode(lang)

	template, ok := c.PromptTemplates[normalizedLang]
	if !ok {
		// Fallback to Japanese if language not found
		template = c.PromptTemplates["japanese"]
	}

	var sb strings.Builder

	// System instruction
	sb.WriteString(template.SystemInstruction)
	sb.WriteString("\n")

	// Guidelines header
	if template.GuidelinesHeader != "" {
		sb.WriteString(template.GuidelinesHeader)
		sb.WriteString("\n")
	}
	for _, guideline := range template.Guidelines {
		sb.WriteString("- ")
		sb.WriteString(guideline)
		sb.WriteString("\n")
	}

	// Current branch
	if template.BranchFormat != "" {
		sb.WriteString(fmt.Sprintf(template.BranchFormat, branch))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Semantic Release prefixes
	if template.SemanticReleaseHeader != "" {
		sb.WriteString(template.SemanticReleaseHeader)
		sb.WriteString("\n")
	}
	if template.SemanticReleaseNote != "" {
		sb.WriteString(template.SemanticReleaseNote)
		sb.WriteString("\n")
	}
	prefixDescriptions := c.GetPrefixDescription(normalizedLang)
	for _, desc := range prefixDescriptions {
		sb.WriteString(desc)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Git diff
	if template.DiffHeader != "" {
		sb.WriteString(template.DiffHeader)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(diff)
	sb.WriteString("\n\n")

	// Output format
	sb.WriteString(template.OutputFormat)

	return sb.String()
}

// BuildPromptEnglish builds an English prompt for commit message generation
// This is a convenience method that calls BuildPrompt with "english" language
func (c *Config) BuildPromptEnglish(branch string, diff string) string {
	return c.BuildPrompt("english", branch, diff)
}
