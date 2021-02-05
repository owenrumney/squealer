package match

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"
)

type Matcher struct {
	test *regexp.Regexp
}

type Matchers []*Matcher

type MatcherController struct {
	matchers       Matchers
	transgressions *Transgressions
}

func NewMatcherController() *MatcherController {
	return &MatcherController{
		matchers:       []*Matcher{},
		transgressions: newTransgressions(),
	}
}

func (mc *MatcherController) Add(regex string) error {
	compile, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("failed to compile the regex. %v", err.Error())
	}
	mc.matchers = append(mc.matchers, &Matcher{
		test: compile,
	})
	return nil
}

func (mc *MatcherController) Evaluate(file *object.File, commit *object.Commit) error {
	content, err := file.Contents()
	if err != nil {
		return err
	}
	for _, matcher := range mc.matchers {
		if matcher.test.MatchString(content) {
			mc.addTransgression(&content, file.Name, commit, matcher)
		}
	}
	return nil
}

func (mc *MatcherController) addTransgression(content *string, name string, commit *object.Commit, matcher *Matcher) {
	lines := strings.Split(*content, "\n")

	m := matcher.test.FindString(*content)
	if len(m) > 0 {
		lineNo, lineContent := lineInFile(m, lines)
		key := fmt.Sprintf("%s:%s", name, hashSimple(m))
		if !mc.transgressions.Exists(key) {
			mc.transgressions.Add(key, newTransgression(lineNo, lineContent, commit, name, m))
		}
	}
}

func lineInFile(m string, lines []string) (int, string) {
	for i, line := range lines {
		if strings.Contains(line, m) {
			return i + 1, line
		}
	}
	return -1, ""
}
