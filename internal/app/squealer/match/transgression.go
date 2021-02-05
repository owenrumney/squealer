package match

import (
	"fmt"
	"strings"
)

type Transgression struct {
	lineNo         int
	lineContent    string
	commit         string
	committerEmail string
	committer      string
	filename       string
	hash           string
	match          string
	redacted       string
}

func newTransgression(lineContent, filename, match, hash string) Transgression {
	content := strings.TrimSpace(lineContent)

	return Transgression{
		lineContent: content,
		filename:    filename,
		hash:        hash,
		match:       match,
		redacted:    strings.ReplaceAll(content, match, "REDACTED"),
	}
}

func (t Transgression) String() string {
	return fmt.Sprintf(`content:  | %s
filename: | %s
hash:     | %s
rule:     | %s:%s
	`, t.lineContent, t.filename, t.hash, t.filename, t.hash)
}

func (t Transgression) Redacted() string {
	return fmt.Sprintf(`content:  | %s
filename: | %s
hash:     | %s
rule:     | %s:%s
	`, t.redacted, t.filename, t.hash, t.filename, t.hash)
}
