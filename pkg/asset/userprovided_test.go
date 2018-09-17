package asset

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryUser(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		defaultValue   string
		expectedOutput string
	}{
		{
			name:           "single word",
			input:          "test",
			defaultValue:   "default",
			expectedOutput: "test",
		},
		{
			name: "with newline",
			input: `test
`,
			defaultValue:   "default",
			expectedOutput: "test",
		},
		{
			name:           "multiple words",
			input:          "this is a test",
			defaultValue:   "default",
			expectedOutput: "this is a test",
		},
		{
			name: "multiple lines",
			input: `first
second`,
			defaultValue:   "default",
			expectedOutput: "first",
		},
		{
			name:           "empty",
			input:          "",
			defaultValue:   "",
			expectedOutput: "",
		},
		{
			name:           "default value",
			input:          "",
			defaultValue:   "default",
			expectedOutput: "default",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := bufio.NewReader(strings.NewReader(tc.input))
			actual := QueryUser(inputReader, "prompt", tc.defaultValue)
			assert.Equal(t, tc.expectedOutput, actual, "unexpected output")
		})
	}
}
