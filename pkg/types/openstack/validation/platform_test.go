package validation

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/validation/mock"
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

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name                  string
		platform              *openstack.Platform
		networking            *types.Networking
		noClouds              bool
		noNetworks            bool
		noFlavors             bool
		noNetExts             bool
		noServiceCatalog      bool
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
			name: "missing cloud",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.Cloud = ""
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "missing external network",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				return p
			}(),
			networking: validNetworking(),
			valid:      false,
		},
		{
			name: "valid default machine pool",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &openstack.MachinePool{}
				return p
			}(),
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
			name:       "clouds fetch failure",
			platform:   validPlatform(),
			networking: validNetworking(),
			noClouds:   true,
			valid:      false,
		},
		{
			name:       "networks fetch failure",
			platform:   validPlatform(),
			networking: validNetworking(),
			noNetworks: true,
			valid:      false,
		},
		{
			name:       "flavors fetch failure",
			platform:   validPlatform(),
			networking: validNetworking(),
			noFlavors:  true,
			valid:      false,
		},
		{
			name:       "network extensions fetch failure",
			platform:   validPlatform(),
			networking: validNetworking(),
			noNetExts:  true,
			valid:      true,
		},
		{
			name:             "service catalog fetch failure",
			platform:         validPlatform(),
			networking:       validNetworking(),
			noServiceCatalog: true,
			valid:            true,
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
			name: "valid MachinesSubnet",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "c664df47-4f7e-4852-819e-e66f9882b7b3"
				return p
			}(),
			networking:          validNetworking(),
			validMachinesSubnet: true,
			valid:               true,
		},
		{
			name: "valid MachinesSubnet invalid machineNetwork",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "c664df47-4f7e-4852-819e-e66f9882b7b3"
				return p
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				n.MachineNetwork[0].CIDR = *ipnet.MustParseCIDR("105.90.0.0/16")
				return n
			}(),
			validMachinesSubnet: true,
			valid:               false,
		},
		{
			name: "invalid MachinesSubnet",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "subnet-c17b"
				return p
			}(),
			networking:            validNetworking(),
			invalidMachinesSubnet: true,
			valid:                 false,
		},
		{
			name: "valid MachinesSubnet externalDNS set",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "c664df47-4f7e-4852-819e-e66f9882b7b3"
				p.ExternalDNS = []string{
					"192.168.1.12",
					"10.0.5.16",
				}
				return p
			}(),
			networking:          validNetworking(),
			validMachinesSubnet: true,
			valid:               false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fetcher := mock.NewMockValidValuesFetcher(mockCtrl)
			if tc.noClouds {
				fetcher.EXPECT().GetCloudNames().
					Return(nil, errors.New("no clouds"))
			} else {
				fetcher.EXPECT().GetCloudNames().
					Return([]string{"test-cloud"}, nil)
			}
			if tc.noNetworks {
				fetcher.EXPECT().GetNetworkNames(tc.platform.Cloud).
					Return(nil, errors.New("no networks")).
					MaxTimes(1)
			} else {
				fetcher.EXPECT().GetNetworkNames(tc.platform.Cloud).
					Return([]string{"test-network"}, nil).
					MaxTimes(1)
			}
			if tc.noFlavors {
				fetcher.EXPECT().GetFlavorNames(tc.platform.Cloud).
					Return(nil, errors.New("no flavors")).
					MaxTimes(1)
			} else {
				fetcher.EXPECT().GetFlavorNames(tc.platform.Cloud).
					Return([]string{"test-flavor"}, nil).
					MaxTimes(1)
			}
			if tc.noNetExts {
				fetcher.EXPECT().GetNetworkExtensionsAliases(tc.platform.Cloud).
					Return(nil, errors.New("no network extensions")).
					MaxTimes(1)
			} else {
				fetcher.EXPECT().GetNetworkExtensionsAliases(tc.platform.Cloud).
					Return([]string{"trunk"}, nil).
					MaxTimes(1)
			}
			if tc.noServiceCatalog {
				fetcher.EXPECT().GetServiceCatalog(tc.platform.Cloud).
					Return(nil, errors.New("no service catalog")).
					MaxTimes(1)
			} else {
				fetcher.EXPECT().GetServiceCatalog(tc.platform.Cloud).
					Return([]string{"octavia"}, nil).
					MaxTimes(1)
			}
			if tc.validMachinesSubnet {
				fetcher.EXPECT().GetSubnetCIDR(tc.platform.Cloud, tc.platform.MachinesSubnet).
					Return("10.0.0.0/16", nil).
					MaxTimes(1)
			}
			if tc.invalidMachinesSubnet {
				fetcher.EXPECT().GetSubnetCIDR(tc.platform.Cloud, tc.platform.MachinesSubnet).
					Return("", errors.New("invalid machinesSubnet")).
					MaxTimes(1)
			}

			testConfig := types.InstallConfig{}
			testConfig.ObjectMeta.Name = "test"

			err := ValidatePlatform(tc.platform, tc.networking, field.NewPath("test-path"), fetcher, &testConfig).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
