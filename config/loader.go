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

// BuildPrompt builds a prompt for commit message generation using template replacement
func (c *Config) BuildPrompt(lang string, branch string, diff string) string {
	// Normalize language code
	normalizedLang := normalizeLangCode(lang)

	promptTemplate, ok := c.PromptTemplates[normalizedLang]
	if !ok {
		// Fallback to Japanese if language not found
		promptTemplate = c.PromptTemplates["japanese"]
	}

	// Format guidelines
	guidelinesText := formatGuidelines(promptTemplate.Guidelines)

	// Format semantic release prefixes
	prefixesText := formatSemanticReleasePrefixes(c.SemanticReleasePrefixes, normalizedLang)

	// Replace template variables
	result := promptTemplate.Template
	result = strings.ReplaceAll(result, "{guidelines}", guidelinesText)
	result = strings.ReplaceAll(result, "{branch}", branch)
	result = strings.ReplaceAll(result, "{semantic_release_prefixes}", prefixesText)
	result = strings.ReplaceAll(result, "{diff}", diff)

	return result
}

// formatGuidelines formats guidelines as a bulleted list
func formatGuidelines(guidelines []string) string {
	var sb strings.Builder
	for _, guideline := range guidelines {
		sb.WriteString("- ")
		sb.WriteString(guideline)
		sb.WriteString("\n")
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

// formatSemanticReleasePrefixes formats semantic release prefixes for the prompt
func formatSemanticReleasePrefixes(prefixes []SemanticReleasePrefix, lang string) string {
	var sb strings.Builder
	for _, p := range prefixes {
		var desc string
		if lang == "ja" {
			desc = p.DescriptionJA
		} else {
			desc = p.DescriptionEN
		}
		sb.WriteString(fmt.Sprintf("\t- \"%s: %s\" : %s\n", p.Type, p.Emoji, desc))
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

// BuildPromptEnglish builds an English prompt for commit message generation
// This is a convenience method that calls BuildPrompt with "english" language
func (c *Config) BuildPromptEnglish(branch string, diff string) string {
	return c.BuildPrompt("english", branch, diff)
}

// normalizeLangCode normalizes language codes to a standard format
func normalizeLangCode(lang string) string {
	switch strings.ToLower(lang) {
	case "ja", "japanese", "jp", "jpn":
		return "ja"
	case "en", "english", "eng":
		return "en"
	default:
		// Explicit default to English
		return "en"
	}
}
