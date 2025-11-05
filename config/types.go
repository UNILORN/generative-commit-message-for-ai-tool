package config

// Config represents the entire configuration
type Config struct {
	PromptTemplates        map[string]PromptTemplate `yaml:"prompt_templates"`
	SemanticReleasePrefixes []SemanticReleasePrefix  `yaml:"semantic_release_prefixes"`
}

// PromptTemplate represents a template for generating commit messages
type PromptTemplate struct {
	SystemInstruction string   `yaml:"system_instruction"`
	Guidelines        []string `yaml:"guidelines"`
	OutputFormat      string   `yaml:"output_format"`
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
	descriptions := make([]string, len(c.SemanticReleasePrefixes))
	for i, p := range c.SemanticReleasePrefixes {
		var desc string
		if lang == "ja" {
			desc = p.DescriptionJA
		} else {
			desc = p.DescriptionEN
		}
		descriptions[i] = "\t- \"" + p.Type + ": " + p.Emoji + "\" : " + desc
	}
	return descriptions
}
