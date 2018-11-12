package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlatformNamesSorted(t *testing.T) {
	sorted := make([]string, len(PlatformNames))
	copy(sorted, PlatformNames)
	sort.Strings(sorted)
	assert.Equal(t, sorted, PlatformNames)
}
