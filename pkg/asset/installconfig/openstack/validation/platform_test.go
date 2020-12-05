package validation

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
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
)

// Returns a default install
func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:             validCloud,
		ExternalNetwork:   validExternalNetwork,
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
		APIFIP: &floatingips.FloatingIP{
			ID:     validFIP1,
			Status: "DOWN",
		},
		IngressFIP: &floatingips.FloatingIP{
			ID:     validFIP2,
			Status: "DOWN",
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
		{
			name:     "not found api FIP",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.APIFIP = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `platform.openstack.lbFloatingIP: Not found: "128.35.27.8"`,
		},
		{
			name:     "not found ingress FIP",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.IngressFIP = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `platform.openstack.ingressFloatingIP: Not found: "128.35.27.13"`,
		},
		{
			name:     "not found both FIPs",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.IngressFIP = nil
				ci.APIFIP = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `[platform.openstack.lbFloatingIP: Not found: "128.35.27.8", platform.openstack.ingressFloatingIP: Not found: "128.35.27.13"]`,
		},
		{
			name:     "in use ingress FIP",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.IngressFIP.Status = "ACTIVE"
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `platform.openstack.ingressFloatingIP: Invalid value: "128.35.27.13": Floating IP already in use`,
		},
		{
			name:     "in use api FIP",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.APIFIP.Status = "ACTIVE"
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `platform.openstack.lbFloatingIP: Invalid value: "128.35.27.8": Floating IP already in use`,
		},
		{
			name: "invalid usage both FIPs",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `[platform.openstack.ingressFloatingIP: Invalid value: "128.35.27.13": Cannot set floating ips when external network not specified, platform.openstack.lbFloatingIP: Invalid value: "128.35.27.8": Cannot set floating ips when external network not specified]`,
		},
		{
			name: "no external network provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				p.LbFloatingIP = ""
				p.IngressFloatingIP = ""
				p.APIVIP = ""
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.ExternalNetwork = nil
				ci.IngressFIP = nil
				ci.APIFIP = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "valid external network",
			platform:       validPlatform(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:     "external network not found",
			platform: validPlatform(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.ExternalNetwork = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: "platform.openstack.externalNetwork: Not found: \"valid-external-network\"",
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
