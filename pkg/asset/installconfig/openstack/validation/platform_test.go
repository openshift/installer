package validation

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/openshift/installer/pkg/ipnet"
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
		APIFloatingIP:     validFIP1,
		Cloud:             validCloud,
		ExternalNetwork:   validExternalNetwork,
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
			expectedErrMsg: `platform.openstack.apiFloatingIP: Not found: "128.35.27.8"`,
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
			expectedErrMsg: `[platform.openstack.apiFloatingIP: Not found: "128.35.27.8", platform.openstack.ingressFloatingIP: Not found: "128.35.27.13"]`,
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
			expectedErrMsg: `platform.openstack.apiFloatingIP: Invalid value: "128.35.27.8": Floating IP already in use`,
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
			expectedErrMsg: `[platform.openstack.ingressFloatingIP: Invalid value: "128.35.27.13": Cannot set floating ips when external network not specified, platform.openstack.apiFloatingIP: Invalid value: "128.35.27.8": Cannot set floating ips when external network not specified]`,
		},
		{
			name: "ingress and API FIPs identical",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.IngressFloatingIP = p.APIFloatingIP
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: `platform.openstack.ingressFloatingIP: Invalid value: "128.35.27.8": ingressFloatingIP can not be the same as apiFloatingIP`,
		},
		{
			name: "no external network provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalNetwork = ""
				p.APIFloatingIP = ""
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
		{
			name: "APIVIP and IngressVIP are the same",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.APIVIP = "128.35.27.8"
				p.IngressVIP = "128.35.27.8"
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: "platform.openstack.ingressVIP: Invalid value: \"128.35.27.8\": ingressVIP can not be the same as apiVIP",
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

func TestClusterOSImage(t *testing.T) {
	cases := []struct {
		name           string
		platform       *openstack.Platform
		cloudInfo      *CloudInfo
		networking     *types.Networking
		expectedErrMsg string // NOTE: this is a REGEXP
	}{
		{
			name:           "no image provided",
			platform:       validPlatform(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: "",
		},
		{
			name: "HTTP address instead of the image name",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "http://example.com/myrhcos.iso"
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: "",
		},
		{
			name: "file location instead of the image name",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "file:///home/user/myrhcos.iso"
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: "",
		},
		{
			name: "valid image",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "my-rhcos"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.OSImage = &images.Image{
					Name:   "my-rhcos",
					Status: images.ImageStatusActive,
				}
				return ci
			}(),
			networking:     validNetworking(),
			expectedErrMsg: "",
		},
		{
			name: "image with invalid status",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "my-rhcos"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.OSImage = &images.Image{
					Name:   "my-rhcos",
					Status: images.ImageStatusSaving,
				}
				return ci
			}(),
			networking:     validNetworking(),
			expectedErrMsg: "platform.openstack.clusterOSImage: Invalid value: \"my-rhcos\": OS image must be active but its status is 'saving'",
		},
		{
			name: "image not found",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "my-rhcos"
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: "platform.openstack.clusterOSImage: Not found: \"my-rhcos\"",
		},
		{
			name: "Unsupported image URL scheme",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ClusterOSImage = "s3://mybucket/myrhcos.iso"
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: "platform.openstack.clusterOSImage: Invalid value: \"s3://mybucket/myrhcos.iso\": URL scheme should be either http\\(s\\) or file but it is 's3'",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aggregatedErrors := ValidatePlatform(tc.platform, tc.networking, tc.cloudInfo).ToAggregate()
			if tc.expectedErrMsg != "" {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestMachineSubnet(t *testing.T) {
	cases := []struct {
		name           string
		platform       *openstack.Platform
		cloudInfo      *CloudInfo
		networking     *types.Networking
		expectedErrMsg string // NOTE: this is a REGEXP
	}{
		{
			name: "external dns is not supported",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "031a5b9d-5a89-4465-8d54-3517ec2bad48"
				p.ExternalDNS = append(p.ExternalDNS, "1.2.3.4")
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: `platform.openstack.externalDNS: Invalid value: \[\]string{"1.2.3.4"}: externalDNS is set, externalDNS is not supported when machinesSubnet is set`,
		},
		{
			name: "machine subnet not found",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "031a5b9d-5a89-4465-8d54-3517ec2bad48"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.MachinesSubnet = nil
				return ci
			}(),
			networking:     validNetworking(),
			expectedErrMsg: `platform.openstack.machinesSubnet: Not found: "031a5b9d-5a89-4465-8d54-3517ec2bad48"`,
		},
		{
			name: "invalid subnet ID",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "fake"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.MachinesSubnet = &subnets.Subnet{
					ID: "031a5b9d-5a89-4465-8d54-3517ec2bad48",
				}
				return ci
			}(),
			networking:     validNetworking(),
			expectedErrMsg: `platform.openstack.machinesSubnet: Internal error: invalid subnet ID`,
		},
		{
			name: "doesn't match the CIDR",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "031a5b9d-5a89-4465-8d54-3517ec2bad48"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.MachinesSubnet = &subnets.Subnet{
					ID: "031a5b9d-5a89-4465-8d54-3517ec2bad48",
				}
				ci.MachinesSubnet.CIDR = "172.0.0.1/16"
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("172.0.0.1/24"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: `platform.openstack.machinesSubnet: Internal error: the first CIDR in machineNetwork, 172.0.0.1/24, doesn't match the CIDR of the machineSubnet, 172.0.0.1/16`,
		},
		{
			name: "valid machine subnet",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.MachinesSubnet = "031a5b9d-5a89-4465-8d54-3517ec2bad48"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				ci.MachinesSubnet = &subnets.Subnet{
					ID: "031a5b9d-5a89-4465-8d54-3517ec2bad48",
				}
				ci.MachinesSubnet.CIDR = "172.0.0.1/24"
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("172.0.0.1/24"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aggregatedErrors := ValidatePlatform(tc.platform, tc.networking, tc.cloudInfo).ToAggregate()
			if tc.expectedErrMsg != "" {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
