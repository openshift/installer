package validation

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/stretchr/testify/assert"
)

var (
	validCloud           = "valid-cloud"
	validExternalNetwork = "valid-external-network"
	validFIP1            = "128.35.27.8"
	validFIP2            = "128.35.27.13"
	validCtrlPlaneFlavor = "valid-control-plane-flavor"
)

// Returns a default install
func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:             validCloud,
		ExternalNetwork:   validExternalNetwork,
		FlavorName:        validCtrlPlaneFlavor,
		LbFloatingIP:      validFIP1,
		IngressFloatingIP: validFIP2,
	}
}

func validNetworking() *types.Networking {
	return &types.Networking{}
}

func validPlatformCloudInfo() *CloudInfo {
	return &CloudInfo{
		ExternalNetwork: &networks.Network{
			ID:           "71b97520-69af-4c35-8153-cdf827z96e60",
			Name:         validExternalNetwork,
			AdminStateUp: true,
			Status:       "ACTIVE",
		},
		PlatformFlavor: &flavors.Flavor{
			Name: validCtrlPlaneFlavor,
		},
	}
}

func TestOpenStackPlatformValidation(t *testing.T) {
	cases := []struct {
		name           string
		platform       *openstack.Platform
		cloudInfo      *CloudInfo
		networking     *types.Networking
		expectedError  bool
		expectedErrMsg string // NOTE: this is a REGEXP
	}{
		{
			name:           "valid platform",
			platform:       validPlatform(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedError:  false,
			expectedErrMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aggregatedErrors := ValidatePlatform(tc.platform, tc.networking, tc.cloudInfo).ToAggregate()
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
