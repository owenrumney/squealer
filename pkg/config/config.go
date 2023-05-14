package config

type Config struct {
	IncludeDefault   bool            `yaml:"include_default" json:"include_default"`
	Rules            []MatchRule     `yaml:"rules" json:"rules"`
	IgnorePaths      []string        `yaml:"ignore_paths" json:"ignore_paths"`
	IgnoreExtensions []string        `yaml:"ignore_extensions" json:"ignore_extensions"`
	Exceptions       []RuleException `yaml:"exceptions" json:"exceptions"`
}

type MatchRule struct {
	Rule        string `yaml:"rule" json:"rule"`
	Description string `yaml:"description" json:"description"`
	FileFilter  string `yaml:"file" json:"file"`
	Entropy     string
}

type RuleException struct {
	ExceptionString string `yaml:"exception" json:"exception"`
	Reason          string `yaml:"reason" json:"reason"`
}
