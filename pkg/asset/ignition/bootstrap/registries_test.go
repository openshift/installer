package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func TestMergedMirrorSets(t *testing.T) {
	tests := []struct {
		name     string
		input    []types.ImageDigestSource
		expected []types.ImageDigestSource
	}{{
		input: []types.ImageDigestSource{{
			Source: "a",
		}, {
			Source: "b",
		}},
		expected: []types.ImageDigestSource{{
			Source: "a",
		}, {
			Source: "b",
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source: "a",
		}, {
			Source: "a",
		}},
		expected: []types.ImageDigestSource{{
			Source: "a",
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"mc", "mc", "md"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "mc", "md"},
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "b",
			Mirrors: []string{"mc", "md"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "b",
			Mirrors: []string{"mc", "md"},
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"ma", "md"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "md"},
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"md", "ma"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "md"},
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"md", "ma"},
		}, {
			Source:  "a",
			Mirrors: []string{"me", "mb"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "md", "me"},
		}},
	}, {
		input: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma"},
		}, {
			Source:  "b",
			Mirrors: []string{"md", "mc"},
		}, {
			Source:  "a",
			Mirrors: []string{"mb", "ma"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb"},
		}, {
			Source:  "b",
			Mirrors: []string{"md", "mc"},
		}},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, MergedMirrorSets(test.input))
		})
	}
}

func TestContentSourceToDigestMirror(t *testing.T) {
	tests := []struct {
		name     string
		input    []types.ImageContentSource
		expected []types.ImageDigestSource
	}{{
		input: []types.ImageContentSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"mc", "mc", "md"},
		}},
		expected: []types.ImageDigestSource{{
			Source:  "a",
			Mirrors: []string{"ma", "mb", "mb"},
		}, {
			Source:  "a",
			Mirrors: []string{"mc", "mc", "md"},
		}},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, ContentSourceToDigestMirror(test.input))
		})
	}
}
