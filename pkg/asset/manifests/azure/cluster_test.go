package azure

import (
	"context"
	"testing"

	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

const (
	testNetworkResourceGroup = "test-network-rg"
	testVirtualNetwork       = "test-vnet"
)

// TestGetNextAvailableIPForLoadBalancer covers
// getNextAvailableIPForLoadBalancer. Add new cases to the table below.
func TestGetNextAvailableIPForLoadBalancer(t *testing.T) {
	controlPlaneSubnet := azuretypes.SubnetSpec{Name: "control-plane-subnet", Role: capz.SubnetControlPlane}

	platform := &azuretypes.Platform{
		CloudName:                azuretypes.PublicCloud,
		NetworkResourceGroupName: testNetworkResourceGroup,
		VirtualNetwork:           testVirtualNetwork,
		Subnets:                  []azuretypes.SubnetSpec{controlPlaneSubnet},
	}
	validInstallConfig := installconfig.MakeAsset(&types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
			},
		},
		Platform: types.Platform{
			Azure: platform,
		},
	})
	validInstallConfig.Azure = icazure.NewMetadata(platform, nil, nil)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// Control plane subnet without any address prefixes.
	azureClient.EXPECT().
		GetControlPlaneSubnet(gomock.Any(), testNetworkResourceGroup, testVirtualNetwork, controlPlaneSubnet.Name).
		Return(&aznetwork.Subnet{SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{}}, nil).
		AnyTimes()

	tests := []struct {
		name string
		// installConfig is the install config passed to the function.
		installConfig *installconfig.InstallConfig
		// lbip is the default load balancer IP.
		lbip string
		// want is the IP the function is expected to return.
		want string
		// wantErr is the exact error the function is expected to return.
		wantErr string
	}{
		{
			name:          "errors without wrapping a nil error when the control plane subnet has no address prefixes",
			installConfig: validInstallConfig,
			lbip:          "10.0.0.100",
			wantErr:       "failed to get control plane subnet addresses",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.installConfig.Azure.UseMockClient(azureClient)

			cidrs := capiutils.CIDRsFromInstallConfig(tc.installConfig)
			got, err := getNextAvailableIPForLoadBalancer(context.Background(), tc.installConfig, cidrs, tc.lbip)

			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
