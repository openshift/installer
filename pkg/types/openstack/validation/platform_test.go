package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:           "test-cloud",
		ExternalNetwork: "test-network",
		DefaultMachinePlatform: &openstack.MachinePool{
			FlavorName: "test-flavor",
		},
	}
}

func validNetworking() *types.Networking {
	return &types.Networking{
		NetworkType: "OVNKubernetes",
		MachineNetwork: []types.MachineNetworkEntry{{
			CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
		}},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name                  string
		platform              *openstack.Platform
		networking            *types.Networking
		noClouds              bool
		noNetworks            bool
		noFlavors             bool
		validMachinesSubnet   bool
		invalidMachinesSubnet bool
		valid                 bool
	}{
		{
			name:       "minimal",
			platform:   validPlatform(),
			networking: validNetworking(),
			valid:      true,
		},
		{
			name: "non IP external dns",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalDNS = []string{
					"invalid",
				}
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "valid external dns",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalDNS = []string{
					"192.168.1.1",
				}
				return p
			}(),
			networking: validNetworking(),
			valid:      true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testConfig := types.InstallConfig{}
			testConfig.ObjectMeta.Name = "test"

			err := ValidatePlatform(tc.platform, tc.networking, field.NewPath("test-path"), &testConfig).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
