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

	commitHash := commit.Hash.String()
	commitAuthor := commit.Author.Name
	commitEmail := commit.Author.Email

	return Transgression{
		lineNo:         lineNo,
		lineContent:    lineContent,
		commit:         commitHash,
		committer:      commitAuthor,
		committerEmail: commitEmail,
		filename:       filename,
		hash:           hashSimple(match),
		match:          match,
		redacted:       strings.ReplaceAll(lineContent, match, "REDACTED"),
	}
}

func (t Transgression) String() string {
		return fmt.Sprintf(`line:      |%d
	content:   |%s
	commit:    |%s
	committer: |%s (%s)
	filename:  |%s
	hash:      |%s
	`, t.lineNo, t.redacted, t.commit, t.committer, t.committerEmail, t.filename, t.hash)
}

func hashSimple(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum64())
}
