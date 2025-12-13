package config

import "strings"

// Config represents the entire configuration
type Config struct {
	PromptTemplates        map[string]PromptTemplate `yaml:"prompt_templates"`
	SemanticReleasePrefixes []SemanticReleasePrefix  `yaml:"semantic_release_prefixes"`
}

// PromptTemplate represents a template for generating commit messages
type PromptTemplate struct {
	SystemInstruction      string   `yaml:"system_instruction"`
	GuidelinesHeader       string   `yaml:"guidelines_header"`
	BranchFormat           string   `yaml:"branch_format"`
	SemanticReleaseHeader  string   `yaml:"semantic_release_header"`
	SemanticReleaseNote    string   `yaml:"semantic_release_note"`
	DiffHeader             string   `yaml:"diff_header"`
	Guidelines             []string `yaml:"guidelines"`
	OutputFormat           string   `yaml:"output_format"`
}

// SemanticReleasePrefix represents a semantic release prefix type
type SemanticReleasePrefix struct {
	Type          string `yaml:"type"`
	Emoji         string `yaml:"emoji"`
	DescriptionJA string `yaml:"description_ja"`
	DescriptionEN string `yaml:"description_en"`
}

// GetPrefixList returns a list of prefix strings (e.g., "feat:", "fix:")
func (c *Config) GetPrefixList() []string {
	prefixes := make([]string, len(c.SemanticReleasePrefixes))
	for i, p := range c.SemanticReleasePrefixes {
		prefixes[i] = p.Type + ":"
	}
	return prefixes
}

// GetPrefixDescription returns the description lines for a given language
func (c *Config) GetPrefixDescription(lang string) []string {
	normalizedLang := normalizeLangCode(lang)
	descriptions := make([]string, len(c.SemanticReleasePrefixes))
	for i, p := range c.SemanticReleasePrefixes {
		var desc string
		if normalizedLang == "ja" {
			desc = p.DescriptionJA
		} else {
			desc = p.DescriptionEN
		}
		descriptions[i] = "\t- \"" + p.Type + ": " + p.Emoji + "\" : " + desc
	}
	return descriptions
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
