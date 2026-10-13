package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
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

func TestCheckServiceEndpointOverride(t *testing.T) {
	endpoints := []configv1.IBMCloudServiceEndpoint{
		{
			Name: configv1.IBMCloudServiceIAM,
			URL:  "e2e.unittest.local",
		},
		{
			Name: configv1.IBMCloudServiceCOS,
			URL:  "e2e.unittest.local",
		},
		{
			Name: configv1.IBMCloudServiceVPC,
			URL:  "e2e.unittest.local",
		},
	}

	testCases := []struct {
		name           string
		service        configv1.IBMCloudServiceName
		expectedResult string
	}{
		{
			name:           "service not found",
			service:        "new-service",
			expectedResult: "",
		},
		{
			name:           "service found",
			service:        configv1.IBMCloudServiceVPC,
			expectedResult: "e2e.unittest.local",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, CheckServiceEndpointOverride(tc.service, endpoints))
	}
}
