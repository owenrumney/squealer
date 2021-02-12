package match

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
)

type Transgression struct {
	lineContent    string
	filename       string
	hash           string
	match          string
	redacted       string
	committer      string
	committerEmail string
	commitHash     string
	excludeRule    string
	committed      string
}

func newTransgression(lineContent, filename, match, hash string, commit *object.Commit) Transgression {
	content := strings.TrimSpace(lineContent)

	commitHash := "-- not applicable --"
	committerName := "-- not applicable --"
	committerEmail := ""
	when := "-- not applicable --"
	if commit != nil {
		commitHash = commit.Hash.String()
		committerEmail = commit.Committer.Email
		committerName = commit.Committer.Name
		when = commit.Committer.When.String()
	}

	return Transgression{
		lineContent:    content,
		filename:       filename,
		hash:           hash,
		match:          match,
		redacted:       strings.ReplaceAll(content, match, "REDACTED"),
		committer:      committerName,
		committerEmail: committerEmail,
		committed:      when,
		commitHash:     commitHash,
		excludeRule:    fmt.Sprintf("%s:%s", filename, hash),
	}
}

func (t Transgression) String() string {
	return fmt.Sprintf(`
content:      | %s
filename:     | %s
secret hash:  | %s
commit:       | %s
committer:    | %s (%s)
committed:    | %s
exclude rule: | %s
	`, t.lineContent, t.filename, t.hash, t.commitHash, t.committer, t.committerEmail, t.committed, t.excludeRule)
}

func (t Transgression) Redacted() string {
	return fmt.Sprintf(`
content:      | %s
filename:     | %s
secret hash:  | %s
commit:       | %s
committer:    | %s (%s)
committed:    | %s
exclude rule: | %s
	`, t.redacted, t.filename, t.hash, t.commitHash, t.committer, t.committerEmail, t.committed, t.excludeRule)
}

func (t *Transgression) update(t2 Transgression) {
	t.committer = t2.committer
	t.committerEmail = t2.committerEmail
	t.commitHash = t2.commitHash
	t.committed = t2.committed
}
