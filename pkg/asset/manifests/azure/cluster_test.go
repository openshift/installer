package azure

import (
	"context"
	"testing"

	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

func TestSetFailureDomains(t *testing.T) {
	tests := []struct {
		name         string
		hasDefault   bool
		defaultZones []string
		controlZones []string
		computeZones [][]string
		want         clusterv1.FailureDomains
	}{
		{
			name: "no zones",
		},
		{
			name:         "control plane zones only",
			controlZones: []string{"1", "2"},
		},
		{
			name:         "compute zones without default machine platform",
			computeZones: [][]string{{"3"}},
		},
		{
			name:       "default machine platform without zones",
			hasDefault: true,
		},
		{
			name:         "default machine platform without zones includes compute zones",
			hasDefault:   true,
			computeZones: [][]string{{"3"}},
			want:         clusterv1.FailureDomains{"3": {}},
		},
		{
			name:         "default machine platform zones",
			hasDefault:   true,
			defaultZones: []string{"1", "2"},
			want: clusterv1.FailureDomains{
				"1": {},
				"2": {},
			},
		},
		{
			name:         "default and control plane zones are combined",
			hasDefault:   true,
			defaultZones: []string{"1", "2"},
			controlZones: []string{"3"},
			want: clusterv1.FailureDomains{
				"1": {},
				"2": {},
				"3": {},
			},
		},
		{
			name:         "default, control plane and compute zones are combined",
			hasDefault:   true,
			defaultZones: []string{"1", "2"},
			controlZones: []string{"2", "3"},
			computeZones: [][]string{{"3", "4"}, {"5"}},
			want: clusterv1.FailureDomains{
				"1": {},
				"2": {},
				"3": {},
				"4": {},
				"5": {},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			platform := &azuretypes.Platform{}
			if tc.hasDefault {
				platform.DefaultMachinePlatform = &azuretypes.MachinePool{Zones: tc.defaultZones}
			}
			ic := &types.InstallConfig{Platform: types.Platform{Azure: platform}}
			if tc.controlZones != nil {
				ic.ControlPlane = &types.MachinePool{Platform: types.MachinePoolPlatform{
					Azure: &azuretypes.MachinePool{Zones: tc.controlZones},
				}}
			}
			for _, zones := range tc.computeZones {
				ic.Compute = append(ic.Compute, types.MachinePool{Platform: types.MachinePoolPlatform{
					Azure: &azuretypes.MachinePool{Zones: zones},
				}})
			}

			assert.Equal(t, tc.want, setFailureDomains(installconfig.MakeAsset(ic)))
		})
	}
}

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
