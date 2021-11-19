package vsphere

import (
	"fmt"
	"testing"

	"github.com/vmware/govmomi/object"
	vim25types "github.com/vmware/govmomi/vim25/types"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/mock"
	"github.com/stretchr/testify/assert"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validIPIInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Cluster:          "valid_cluster",
				Datacenter:       "valid_dc",
				DefaultDatastore: "valid_ds",
				Folder:           "valid_folder",
				Network:          "valid_network",
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid-vcenter",
				APIVIP:           "192.168.111.0",
				IngressVIP:       "192.168.111.1",
			},
		},
	}
}

func validUPIInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Cluster:          "valid_cluster",
				Datacenter:       "valid_dc",
				DefaultDatastore: "valid_ds",
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid-vcenter",
				Network:          "valid_network",
			},
		},
	}
}

func TestValidateProvision(t *testing.T) {
	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(vsphere.Finder, *types.InstallConfig) error
		expectErr        string
	}{{
		name:             "valid IPI install config",
		installConfig:    validIPIInstallConfig(),
		validationMethod: validateProvisioning,
	}, {
		name: "invalid IPI - no network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Network = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.network: Required value: must specify the network$`,
	}, {
		name: "invalid IPI - no cluster",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Cluster = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.cluster: Required value: must specify the cluster$`,
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	vsphereClient := mock.NewMockFinder(mockCtrl)
	vsphereClient.EXPECT().Folder(gomock.Any(), "valid_folder").AnyTimes().Return(nil, nil)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(vsphereClient, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}

func TestValidateResources(t *testing.T) {
	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(vsphere.Finder, vsphere.NetworkIdentifier, *types.InstallConfig) error
		expectErr        string
	}{{
		name:             "valid UPI install config",
		installConfig:    validUPIInstallConfig(),
		validationMethod: validateResources,
	}, {
		name: "invalid IPI - invalid datacenter",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Datacenter = "invalid_dc"
			return c
		}(),
		validationMethod: validateResources,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_dc": 404$`,
	}, {
		name: "invalid IPI - invalid network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateResources,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
	}}

	ccr := object.ClusterComputeResource{}
	ccr.InventoryPath = "valid_cluster"

	networks := []vim25types.ManagedObjectReference{
		{
			Value: "valid",
		},
		{
			Value: "other",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vsphereClient := mock.NewMockFinder(mockCtrl)
	networkIdentifier := mock.NewMockNetworkIdentifier(mockCtrl)

	vsphereClient.EXPECT().Datacenter(gomock.Any(), "./valid_dc").Return(&object.Datacenter{Common: object.Common{InventoryPath: "valid_dc"}}, nil).AnyTimes()
	vsphereClient.EXPECT().Datacenter(gomock.Any(), gomock.Not("./valid_dc")).Return(nil, fmt.Errorf("404")).AnyTimes()
	vsphereClient.EXPECT().ClusterComputeResource(gomock.Any(), "/valid_dc/host/valid_cluster").Return(&ccr, nil).AnyTimes()
	vsphereClient.EXPECT().ClusterComputeResource(gomock.Any(), gomock.Not("/valid_dc/host/valid_cluster")).Return(nil, fmt.Errorf("404")).AnyTimes()

	networkIdentifier.EXPECT().GetNetworks(&ccr).Return(networks, nil).AnyTimes()
	networkIdentifier.EXPECT().GetNetworkName(networks[0]).Return("valid_network", nil).AnyTimes()
	networkIdentifier.EXPECT().GetNetworkName(networks[1]).Return("other_network", nil).AnyTimes()
	networkIdentifier.EXPECT().GetNetworkName(gomock.Not(networks[0])).Return("", fmt.Errorf("404")).AnyTimes()
	networkIdentifier.EXPECT().GetNetworkName(gomock.Not(networks[1])).Return("", fmt.Errorf("404")).AnyTimes()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(vsphereClient, networkIdentifier, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
