package match

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"

	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
)

type Matcher struct {
	test        *regexp.Regexp
	description string
}

type Matchers []*Matcher

type MatcherController struct {
	matchers       Matchers
	exclusions     []config.RuleException
	transgressions *transgressionMap
	metrics        *mertics.Metrics
	redacted       bool
}

func NewMatcherController(cfg *config.Config, metrics *mertics.Metrics, redacted bool) *MatcherController {
	mc := &MatcherController{
		matchers:       []*Matcher{},
		transgressions: newTransgressions(),
		exclusions:     cfg.Exceptions,
		metrics:        metrics,
		redacted:       redacted,
	}

	for _, rule := range cfg.Rules {
		err := mc.add(rule)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return mc
}

func (mc *MatcherController) add(rule config.MatchRule) error {
	compile, err := regexp.Compile(rule.Rule)
	if err != nil {
		return fmt.Errorf("failed to compile the regex. %v", err.Error())
	}
	mc.matchers = append(mc.matchers, &Matcher{
		test:        compile,
		description: rule.Description,
	})
	return nil
}

func (mc *MatcherController) Evaluate(file *object.File) error {
	content, err := file.Contents()
	if err != nil {
		return err
	}
	for _, matcher := range mc.matchers {
		if matcher.test.MatchString(content) {
			mc.addTransgression(&content, file.Name, matcher)
		}
	}
	return nil
}

func (mc *MatcherController) addTransgression(content *string, name string, matcher *Matcher) {
	lines := strings.Split(*content, "\n")

	m := matcher.test.FindString(*content)
	if len(m) > 0 {
		lineContent := lineInFile(m, lines)
		secretHash := mc.newHash(m)
		key := fmt.Sprintf("%s:%s", name, secretHash)
		mc.metrics.IncrementTransgressionsFound()
		for _, exclusion := range mc.exclusions {
			if exclusion.ExceptionString == key {
				mc.metrics.IncrementTransgressionsIgnored()
				return
			}
		}

		if !mc.transgressions.exists(key) {
			mc.metrics.IncrementTransgressionsReported()
			transgression := newTransgression(lineContent, name, m, secretHash)
			mc.transgressions.add(key, transgression)
			if mc.redacted {
				fmt.Printf(transgression.Redacted())
			} else {
				fmt.Printf(transgression.String())
			}
		}
	}
}

func (mc *MatcherController) newHash(secret string) string {
	hasher := sha1.New()
	hasher.Write([]byte(secret))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

func lineInFile(m string, lines []string) string {
	for _, line := range lines {
		if strings.Contains(line, m) {
			return line
		}
	}
	return ""
}
