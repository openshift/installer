package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterResourceGroupName(t *testing.T) {
	infraID := "infra-id"
	platform := Platform{}
	platform.ResourceGroupName = ""
	assert.Equal(t, infraID, platform.ClusterResourceGroupName(infraID))
	platform.ResourceGroupName = "test-cluster"
	assert.Equal(t, "test-cluster", platform.ClusterResourceGroupName(infraID))
}
