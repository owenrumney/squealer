package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransgressionMap(t *testing.T) {
	tm := newTransgressions()

	assert.NotNil(t, tm)
	assert.Equal(t, 0, tm.count())
}

func TestAddItemToTransgressionMap(t *testing.T) {
	tm := newTransgressions()

	assert.NotNil(t, tm)
	assert.Equal(t, 0, tm.count())

	tm.add("test1", Transgression{
		LineContent:     "testing",
		Filename:        "",
		Hash:            "",
		Match:           "",
		RedactedContent: "",
	})

	tm.add("test1", Transgression{
		LineContent:     "testing2",
		Filename:        "",
		Hash:            "",
		Match:           "",
		RedactedContent: "",
	})
}
