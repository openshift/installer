package vsphere

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
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
				Cluster:          "DC0_C0",
				Datacenter:       "DC0",
				DefaultDatastore: "LocalDS_0",
				Network:          "DC0_DVPG0",
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid-vcenter",
				APIVIP:           "192.168.111.0",
				IngressVIP:       "192.168.111.1",
			},
		},
	}
}

func TestValidate(t *testing.T) {
	server := mock.StartSimulator()
	defer server.Close()

	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(Finder, *types.InstallConfig) error
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
		name: "invalid IPI - invalid datacenter",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Datacenter = "invalid_dc"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_dc": 404$`,
	}, {
		name: "invalid IPI - invalid network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
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
	finder, _ := mock.GetFinder(server)

	/*vsphereClient := mock.NewMockFinder(mockCtrl)
	vsphereClient.EXPECT().Datacenter(gomock.Any(), "./valid_dc").Return(&object.Datacenter{Common: object.Common{InventoryPath: "valid_dc"}}, nil).AnyTimes()
	vsphereClient.EXPECT().Datacenter(gomock.Any(), gomock.Not("./valid_dc")).Return(nil, fmt.Errorf("404")).AnyTimes()

	vsphereClient.EXPECT().Network(gomock.Any(), "valid_dc/network/valid_network").Return(nil, nil).AnyTimes()
	vsphereClient.EXPECT().Network(gomock.Any(), gomock.Not("valid_dc/network/valid_network")).Return(nil, fmt.Errorf("404")).AnyTimes()*/
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(finder, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
