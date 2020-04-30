package vsphere

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
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
				Network:          "valid_network",
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid_vcenter",
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
				Datacenter:       "valid_dc",
				DefaultDatastore: "valid_ds",
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid_vcenter",
			},
		},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*types.InstallConfig) error
		expectErr        string
	}{{
		name:             "valid UPI install config",
		installConfig:    validUPIInstallConfig(),
		validationMethod: Validate,
	}, {
		name:             "valid IPI install config",
		installConfig:    validIPIInstallConfig(),
		validationMethod: ValidateForProvisioning,
	}, {
		name: "invalid IPI - no network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Network = ""
			return c
		}(),
		validationMethod: ValidateForProvisioning,
		expectErr:        `^platform\.vsphere\.network: Required value: must specify the network$`,
	}, {
		name: "invalid IPI - no cluster",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig()
			c.Platform.VSphere.Cluster = ""
			return c
		}(),
		validationMethod: ValidateForProvisioning,
		expectErr:        `^platform\.vsphere\.cluster: Required value: must specify the cluster$`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err.Error())
			}
		})
	}
}
