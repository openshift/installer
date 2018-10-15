package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	config := &InstallConfig{
		ClusterID: "123",
		Platform: Platform{
			AWS: &AWSPlatform{
				UserTags: map[string]string{
					"abc": "def",
				},
			},
		},
	}

	expected := map[string]string{
		"abc":               "def",
		"tectonicClusterID": "123",
	}

	assert.Equal(t, expected, config.Tags())
}
