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

func TestGetVPCName(t *testing.T) {
	platform := Platform{}

	testCases := []struct {
		name           string
		vpcName        string
		expectedResult string
	}{
		{"no vpc name", "", ""},
		{"provided vpc name", "my-vpc", "my-vpc"},
	}

	for _, tc := range testCases {
		platform.VPCName = tc.vpcName
		assert.Equal(t, tc.expectedResult, platform.GetVPCName())
	}
}
