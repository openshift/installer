package version //nolint:revive

import (
	"testing"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name          string
		versionPadded string
		raw           string
		expected      string
	}{
		{
			name:          "injected release version takes precedence",
			versionPadded: "4.18.7\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			raw:           "v5.0.0",
			expected:      "4.18.7",
		},
		{
			name:          "ART build version",
			versionPadded: defaultVersionPadded,
			raw:           "v5.1.0",
			expected:      "5.1.0",
		},
		{
			name:          "git describe version",
			versionPadded: defaultVersionPadded,
			raw:           "v1.5.0-alpha.2-144-g0de7706b19ecffda804f93c08edbbe091abe71ed",
			expected:      "5.0-alpha.2-144-g0de7706b19ecffda804f93c08edbbe091abe71ed",
		},
		{
			name:          "invalid injected version falls back to build version",
			versionPadded: "0.0\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			raw:           "v5.1.0",
			expected:      "5.1.0",
		},
		{
			name:          "injected version with invalid minor falls back to build version",
			versionPadded: "5.invalid\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			raw:           "v5.1.0",
			expected:      "5.1.0",
		},
		{
			name:          "corrupt injected version falls back to build version",
			versionPadded: "5.1.0-without-null-terminator",
			raw:           "v5.1.0",
			expected:      "5.1.0",
		},
		{
			name:          "invalid build version uses default",
			versionPadded: defaultVersionPadded,
			raw:           "was not built correctly",
			expected:      fallbackVersion,
		},
		{
			name:          "zero build version uses default",
			versionPadded: defaultVersionPadded,
			raw:           "0.0",
			expected:      fallbackVersion,
		},
		{
			name:          "build version with invalid minor uses default",
			versionPadded: defaultVersionPadded,
			raw:           "5.invalid",
			expected:      fallbackVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if actual := resolveVersion(tt.versionPadded, tt.raw); actual != tt.expected {
				t.Errorf("resolveVersion() = %q, want %q", actual, tt.expected)
			}
		})
	}
}

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
