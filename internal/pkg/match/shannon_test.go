package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Shannon_EntropyCheck(t *testing.T) {

	tests := []struct {
		testValue      string
		bounds         string
		expectedInside bool
		errorExpected  bool
	}{
		{
			testValue:      "1223334444",
			bounds:         "1,2",
			expectedInside: true,
		},
		{
			testValue:      "23sad12233as2342rwer34d323423444erer4",
			bounds:         "1,2",
			expectedInside: false,
		},
		{
			testValue:      "1223334444",
			bounds:         "4.3,7.2",
			expectedInside: false,
		},
		{
			testValue:     "1223334444",
			bounds:        "4.37.2",
			errorExpected: true,
		},
	}

	for _, tt := range tests {
		inBounds, err := entropyCheck(tt.testValue, tt.bounds)
		assert.Equal(t, tt.expectedInside, inBounds)
		if tt.errorExpected {
			assert.Error(t, err)
		}
	}

}
