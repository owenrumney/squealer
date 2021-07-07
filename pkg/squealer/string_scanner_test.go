package squealer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquealerStrings(t *testing.T) {

	var tests = []struct {
		name                      string
		content                   string
		expectedResultDescription string
		expectedResultStatus      bool
	}{
		{
			name:                      "github token is found to have an issue",
			content:                   "ghp_dsflkj234ASF34wdfkjbslf1234dsfsdfSDFDDSF",
			expectedResultDescription: "Check for new Github Token",
			expectedResultStatus:      true,
		},
		{
			name:                      "string with password has an issue",
			content:                   `password="Password1234"`,
			expectedResultDescription: "Password literal text",
			expectedResultStatus:      true,
		},
		{
			name: "string with password has an issue",
			content: `
			# user data example
			DB_USERNAME="database_user"
			DB_PASSWORD="Password1234"
			`,
			expectedResultDescription: "Password literal text",
			expectedResultStatus:      true,
		},
		{
			name:                      "normal string has no issue",
			content:                   "The quick brown fox jumps over the lazy dog",
			expectedResultDescription: "",
			expectedResultStatus:      false,
		},
	}

	stringSquealer := NewStringScanner()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := stringSquealer.Scan(test.content)
			assert.Equal(t, test.expectedResultDescription, result.Description)
			assert.Equal(t, test.expectedResultStatus, result.TransgressionFound)
		})
	}
}
