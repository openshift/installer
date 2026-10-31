package powervs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types/powervs"
)

func TestNetworkSelector(t *testing.T) {
	const clusterID = "test-abcde"

	cases := []struct {
		name          string
		networkType   powervs.NetworkType
		expectedType  machinev1.PowerVSResourceType
		expectedValue string
	}{
		{
			name:          "default network type uses DHCPSubnet network name",
			networkType:   "",
			expectedType:  machinev1.PowerVSResourceTypeName,
			expectedValue: "test-abcde-network",
		},
		{
			name:          "DHCPSubnet uses network name",
			networkType:   powervs.NetworkTypeDHCPSubnet,
			expectedType:  machinev1.PowerVSResourceTypeName,
			expectedValue: "test-abcde-network",
		},
		{
			name:          "DHCPServer uses legacy regex",
			networkType:   powervs.NetworkTypeDHCPServer,
			expectedType:  machinev1.PowerVSResourceTypeRegEx,
			expectedValue: "^DHCPSERVER.*test-abcde.*_Private$",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &powervs.Platform{NetworkType: tc.networkType}
			config, err := provider(clusterID, platform, &powervs.MachinePool{}, "worker-user-data", "rhcos-"+clusterID, networkName(clusterID, platform))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, config.Network.Type)
			if tc.expectedType == machinev1.PowerVSResourceTypeName {
				assert.Equal(t, tc.expectedValue, *config.Network.Name)
			} else {
				assert.Equal(t, tc.expectedValue, *config.Network.RegEx)
			}
		})
	}
}
