package version //nolint:revive

import (
	"testing"
)

func Test_removeGoVersionPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "go module version v1.5.0",
			input:    "v1.5.0",
			expected: "5.0",
		},
		{
			name:     "go module version v1.4.22",
			input:    "v1.4.22",
			expected: "4.22",
		},
		{
			name:     "go module version with prerelease v1.5.0-alpha.0",
			input:    "v1.5.0-alpha.0",
			expected: "5.0-alpha.0",
		},
		{
			name:     "go module version from git describe",
			input:    "v1.5.0-alpha.0-106-g01a3a762a23abc",
			expected: "5.0-alpha.0-106-g01a3a762a23abc",
		},
		{
			name:     "ART build version v5.0.0",
			input:    "v5.0.0",
			expected: "5.0.0",
		},
		{
			name:     "ART build version v4.18.0",
			input:    "v4.18.0",
			expected: "4.18.0",
		},
		{
			name:     "no prefix",
			input:    "5.0.0",
			expected: "5.0.0",
		},
		{
			name:     "bare version without v",
			input:    "was not built correctly",
			expected: "was not built correctly",
		},
		{
			name:     "v only",
			input:    "v",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeGoVersionPrefix(tt.input)
			if got != tt.expected {
				t.Errorf("removeGoVersionPrefix(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
