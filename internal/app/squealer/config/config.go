package config

type Config struct {
	Rules            []MatchRule     `yaml:"rules" json:"rules"`
	IgnorePrefixes   []string        `yaml:"ignore_prefixes" json:"ignore_prefixes"`
	IgnoreExtensions []string        `yaml:"ignore_extensions" json:"ignore_extensions"`
	Exceptions       []RuleException `yaml:"exceptions" json:"exceptions"`
}

type MatchRule struct {
	Rule        string `yaml:"rule" json:"rule"`
	Description string `yaml:"description" json:"description"`
}

type RuleException struct {
	ExceptionString string `yaml:"exception" json:"exception"`
	Reason          string `yaml:"reason" json:"reason"`
}
