package vsphere

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validIPIInstallConfig(dcName string, fName string) *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Cluster:          fmt.Sprintf("%s/%s_C0", fName, dcName),
				Datacenter:       fmt.Sprintf("%s/%s", fName, dcName),
				DefaultDatastore: "LocalDS_0",
				Network:          fmt.Sprintf("%s_DVPG0", dcName),
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
	dcName := "DC0"
	fName := "F0"
	dcName1 := "DC1"
	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*vim25.Client, Finder, *types.InstallConfig) error
		expectErr        string
	}{{
		name:             "valid IPI install config",
		installConfig:    validIPIInstallConfig(dcName, ""),
		validationMethod: validateProvisioning,
	}, {
		name:             "valid IPI install config - DC in folder",
		installConfig:    validIPIInstallConfig(dcName1, fName),
		validationMethod: validateProvisioning,
	}, {
		name: "invalid IPI - no network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Network = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.network: Required value: must specify the network$`,
	}, {
		name: "invalid IPI - invalid datacenter",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Datacenter = "invalid_dc"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_dc": datacenter './invalid_dc' not found`,
	}, {
		name: "invalid IPI - invalid network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
	}, {
		name: "invalid IPI - invalid network - DC in folder",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName1, fName)
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
	}, {
		name: "invalid IPI - no cluster",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Cluster = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.cluster: Required value: must specify the cluster$`,
	}}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	finder, err := mock.GetFinder(server)
	if err != nil {
		t.Error(err)
		return
	}
	client, _, err := mock.GetClient(server)
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(client, finder, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
