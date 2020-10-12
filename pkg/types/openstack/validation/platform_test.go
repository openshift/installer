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
		FlavorName:      "test-flavor",
	}
}

func validNetworking() *types.Networking {
	return &types.Networking{
		NetworkType: "OpenShiftSDN",
		MachineNetwork: []types.MachineNetworkEntry{{
			CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
		}},
	}
}

func TestValidateClusterName(t *testing.T) {
	testConfig := types.InstallConfig{}

	// valid platform
	testConfig.ObjectMeta.Name = "test"
	errs := ValidatePlatform(validPlatform(), validNetworking(), field.NewPath("test-path"), &testConfig)
	assert.NoError(t, errs.ToAggregate())

	// too long cluster name (more than 14 chars)
	testConfig.ObjectMeta.Name = "0123456789abcde"
	errs = ValidatePlatform(validPlatform(), validNetworking(), field.NewPath("test-path"), &testConfig)
	assert.True(t, len(errs) == 1)
	assert.Equal(t, "cluster name is too long, please restrict it to 14 characters", errs[0].Detail)

	// . in the name
	testConfig.ObjectMeta.Name = "test.cluster"
	errs = ValidatePlatform(validPlatform(), validNetworking(), field.NewPath("test-path"), &testConfig)
	assert.True(t, len(errs) == 1)
	assert.Equal(t, "cluster name can't contain \".\" character", errs[0].Detail)
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
		{
			name: "valid custom API vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.APIVIP = "10.0.0.9"
				return p
			}(),
			networking: validNetworking(),
			valid:      true,
		},
		{
			name: "incorrect network custom API vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.APIVIP = "11.1.0.5"
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "valid custom ingress vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.IngressVIP = "10.0.0.9"
				return p
			}(),
			networking: validNetworking(),
			valid:      true,
		},
		{
			name: "incorrect network custom ingress vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.IngressVIP = "11.1.0.5"
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "invalid network custom ingress vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.IngressVIP = "banana"
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "invalid network custom API vip",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.APIVIP = "banana"
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name:     "cluster network overlaps machine network",
			platform: validPlatform(),
			networking: func() *types.Networking {
				n := validNetworking()
				n.ClusterNetwork = []types.ClusterNetworkEntry{{
					CIDR: *ipnet.MustParseCIDR("10.0.0.0/6"),
				}}
				return n
			}(),
			valid: true,
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
