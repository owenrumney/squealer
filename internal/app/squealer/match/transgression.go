package match

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"hash/fnv"
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

func newTransgression(lineNo int, lineContent string, commit *object.Commit, filename, match string) Transgression {
	content := strings.TrimSpace(lineContent)

	return Transgression{
		lineNo:      lineNo,
		lineContent: content,
		filename:    filename,
		hash:        hashSimple(match),
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

func hashSimple(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum64())
}
