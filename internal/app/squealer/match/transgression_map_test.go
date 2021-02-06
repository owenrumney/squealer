package match

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
		lineContent: "testing",
		filename:    "",
		hash:        "",
		match:       "",
		redacted:    "",
	})
}
