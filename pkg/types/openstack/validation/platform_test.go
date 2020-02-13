package validation

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

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

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name             string
		platform         *openstack.Platform
		noClouds         bool
		noNetworks       bool
		noFlavors        bool
		noNetExts        bool
		noServiceCatalog bool
		valid            bool
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "missing cloud",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.Cloud = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "missing external network",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "valid default machine pool",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &openstack.MachinePool{}
				return p
			}(),
			valid: true,
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
			valid: false,
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
			valid: true,
		},
		{
			name:     "clouds fetch failure",
			platform: validPlatform(),
			noClouds: true,
			valid:    false,
		},
		{
			name:       "networks fetch failure",
			platform:   validPlatform(),
			noNetworks: true,
			valid:      false,
		},
		{
			name:      "flavors fetch failure",
			platform:  validPlatform(),
			noFlavors: true,
			valid:     false,
		},
		{
			name:      "network extensions fetch failure",
			platform:  validPlatform(),
			noNetExts: true,
			valid:     true,
		},
		{
			name:             "service catalog fetch failure",
			platform:         validPlatform(),
			noServiceCatalog: true,
			valid:            true,
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

			testConfig := types.InstallConfig{}
			testConfig.ObjectMeta.Name = "test"

			err := ValidatePlatform(tc.platform, nil, field.NewPath("test-path"), fetcher, &testConfig).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
