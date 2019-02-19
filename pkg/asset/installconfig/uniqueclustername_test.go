package installconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateName(t *testing.T) {
	cases := []struct {
		name     string
		base     string
		expected string
	}{
		{
			name:     "empty",
			base:     "",
			expected: "^-[[:alnum:]]{5}$",
		},
		{
			name:     "short",
			base:     "test-name",
			expected: "^test-name-[[:alnum:]]{5}$",
		},
		{
			name:     "long",
			base:     "012345678901234567890123456789",
			expected: "^0123456789012345678901-[[:alnum:]]{5}$",
		},
	}
	for _, tc := range cases {
		actual := generateName(tc.base)
		assert.Regexp(t, tc.expected, actual)
		assert.True(t, len(actual) <= maxNameLen)
	}
}
