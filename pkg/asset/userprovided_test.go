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
		expectedOutput string
	}{
		{
			name:           "single word",
			input:          "test",
			expectedOutput: "test",
		},
		{
			name: "with newline",
			input: `test
`,
			expectedOutput: "test",
		},
		{
			name:           "multiple words",
			input:          "this is a test",
			expectedOutput: "this is a test",
		},
		{
			name: "multiple lines",
			input: `first
second`,
			expectedOutput: "first",
		},
		{
			name:           "empty",
			input:          "",
			expectedOutput: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := bufio.NewReader(strings.NewReader(tc.input))
			actual := QueryUser(inputReader, "prompt")
			assert.Equal(t, tc.expectedOutput, actual, "unexpected output")
		})
	}
}
