package validation

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

var (
	validCloud           = "valid-cloud"
	validExternalNetwork = "valid-external-network"
	validFIP1            = "128.35.27.8"
	validFIP2            = "128.35.27.13"
	validSubnetID        = "031a5b9d-5a89-4465-8d54-3517ec2bad48"
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

func validControlPlanePort() *openstack.PortTarget {
	fixedIP := openstack.FixedIP{
		Subnet: openstack.SubnetFilter{ID: validSubnetID, Name: "valid-subnet"},
	}
	controlPlanePort := &openstack.PortTarget{
		FixedIPs: []openstack.FixedIP{fixedIP},
	}
	return controlPlanePort
}

func validNetworking() *types.Networking {
	return &types.Networking{}
}

func withControlPlanePortSubnets(subnetCIDR, allocationPoolStart, allocationPoolEnd string) func(*CloudInfo) {
	return func(ci *CloudInfo) {
		subnet := subnets.Subnet{
			CIDR: subnetCIDR,
			AllocationPools: []subnets.AllocationPool{
				{Start: allocationPoolStart, End: allocationPoolEnd},
			},
		}
		allsubnets := []*subnets.Subnet{&subnet}
		ci.ControlPlanePortSubnets = allsubnets
	}
}
func validPlatformCloudInfo(options ...func(*CloudInfo)) *CloudInfo {
	ci := CloudInfo{
		ExternalNetwork: &Network{
			networks.Network{
				ID:           "71b97520-69af-4c35-8153-cdf827z96e60",
				Name:         validExternalNetwork,
				AdminStateUp: true,
				Status:       "ACTIVE",
			},
			mtu.NetworkMTUExt{},
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

	for _, apply := range options {
		apply(&ci)
	}

	return &ci
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
			name: "APIVIP inside subnet allocation pool",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.APIVIPs = []string{"10.0.128.10"}
				return p
			}(),
			cloudInfo: validPlatformCloudInfo(withControlPlanePortSubnets(
				"10.0.128.0/24",
				"10.0.128.8",
				"10.0.128.255",
			)),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: "platform.openstack.apiVIPs: Invalid value: \"10.0.128.10\": apiVIP can not fall in a MachineNetwork allocation pool",
		},
		{
			name: "ingressVIP inside subnet allocation pool",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.IngressVIPs = []string{"10.0.128.42"}
				return p
			}(),
			cloudInfo: validPlatformCloudInfo(withControlPlanePortSubnets(
				"10.0.128.0/24",
				"10.0.128.8",
				"10.0.128.255",
			)),
			networking:     validNetworking(),
			expectedError:  true,
			expectedErrMsg: "platform.openstack.ingressVIPs: Invalid value: \"10.0.128.42\": ingressVIP can not fall in a MachineNetwork allocation pool",
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
				p.ExternalDNS = append(p.ExternalDNS, "1.2.3.4")
				p.ControlPlanePort = validControlPlanePort()
				return p
			}(),
			cloudInfo:      validPlatformCloudInfo(),
			networking:     validNetworking(),
			expectedErrMsg: `platform.openstack.externalDNS: Invalid value: \[\]string{"1.2.3.4"}: externalDNS is set, externalDNS is not supported when ControlPlanePort is set`,
		},
		{
			name: "control plane port subnet not found",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ControlPlanePort = validControlPlanePort()
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID: "00000000-5a89-4465-8d54-3517ec2bad48",
				}
				allsubnets := []*subnets.Subnet{&subnet}
				ci.ControlPlanePortSubnets = allsubnets
				return ci
			}(),
			networking:     validNetworking(),
			expectedErrMsg: `platform.openstack.controlPlanePort.fixedIPs: Not found: "031a5b9d-5a89-4465-8d54-3517ec2bad48"`,
		},
		{
			name: "network does not contain subnets",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ControlPlanePort = validControlPlanePort()
				p.ControlPlanePort.Network.ID = "00000000-2a22-4465-8d54-3517ec2bad48"
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID:        "031a5b9d-5a89-4465-8d54-3517ec2bad48",
					NetworkID: "00000000-1a11-4465-8d54-3517ec2bad48",
					CIDR:      "172.0.0.1/24",
				}
				allSubnets := []*subnets.Subnet{&subnet}
				ci.ControlPlanePortSubnets = allSubnets
				network := Network{}
				network.ID = "00000000-2a22-4465-8d54-3517ec2bad48"
				ci.ControlPlanePortNetwork = &network
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
			expectedErrMsg: `platform.openstack.controlPlanePort.network: Invalid value: "00000000-2a22-4465-8d54-3517ec2bad48": network must contain subnets`,
		},
		{
			name: "doesn't match the CIDR",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ControlPlanePort = validControlPlanePort()
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID:   validSubnetID,
					CIDR: "172.0.0.1/16",
				}
				allsubnets := []*subnets.Subnet{&subnet}
				ci.ControlPlanePortSubnets = allsubnets
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
			expectedErrMsg: `platform.openstack.controlPlanePort.fixedIPs: Invalid value: "172.0.0.1/16": controlPlanePort CIDR does not match machineNetwork`,
		},
		{
			name: "control plane port subnets on different network",
			platform: func() *openstack.Platform {
				p := validPlatform()
				fixedIP := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-5a89-4465-8d54-3517ec2bad48"},
				}
				fixedIPv6 := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-1111-4465-8d54-3517ec2bad48"},
				}
				p.ControlPlanePort = &openstack.PortTarget{
					FixedIPs: []openstack.FixedIP{fixedIP, fixedIPv6},
				}
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID:        "00000000-5a89-4465-8d54-3517ec2bad48",
					NetworkID: "00000000-2222-4465-8d54-3517ec2bad48",
					CIDR:      "172.0.0.1/16",
					IPVersion: 4,
				}
				subnetv6 := subnets.Subnet{
					ID:        "00000000-1111-4465-8d54-3517ec2bad48",
					NetworkID: "00000000-3333-4465-8d54-3517ec2bad48",
					CIDR:      "2001:db8::/64",
					IPVersion: 6,
				}
				allsubnets := []*subnets.Subnet{&subnet, &subnetv6}
				ci.ControlPlanePortSubnets = allsubnets
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("172.0.0.1/16"),
				}
				machineNetworkEntryv6 := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("2001:db8::/64"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry, *machineNetworkEntryv6}
				return n
			}(),
			expectedErrMsg: `platform.openstack.controlPlanePort.fixedIPs: Invalid value: "00000000-3333-4465-8d54-3517ec2bad48": fixedIPs subnets must be on the same Network`,
		},
		{
			name: "valid control plane port",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ControlPlanePort = validControlPlanePort()
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID:   validSubnetID,
					CIDR: "172.0.0.1/16",
				}
				allsubnets := []*subnets.Subnet{&subnet}
				ci.ControlPlanePortSubnets = allsubnets
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("172.0.0.1/16"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: "",
		},
		{
			name: "control plane port multiple ipv4 subnets",
			platform: func() *openstack.Platform {
				p := validPlatform()
				fixedIP := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-5a89-4465-8d54-3517ec2bad48"},
				}
				fixedIPv6 := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-1111-4465-8d54-3517ec2bad48"},
				}
				p.ControlPlanePort = &openstack.PortTarget{
					FixedIPs: []openstack.FixedIP{fixedIP, fixedIPv6},
				}
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnet := subnets.Subnet{
					ID:        "00000000-5a89-4465-8d54-3517ec2bad48",
					CIDR:      "172.0.0.1/16",
					IPVersion: 4,
				}
				subnetv6 := subnets.Subnet{
					ID:        "00000000-1111-4465-8d54-3517ec2bad48",
					CIDR:      "10.0.0.0/16",
					IPVersion: 4,
				}
				allsubnets := []*subnets.Subnet{&subnet, &subnetv6}
				ci.ControlPlanePortSubnets = allsubnets
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("172.0.0.1/16"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: `[platform.openstack.controlPlanePort.fixedIPs: Internal error: controlPlanePort CIDRs does not match machineNetwork, platform.openstack.controlPlanePort.fixedIPs: Internal error: multiple IPv4 subnets is not supported]`,
		},
		{
			name: "control plane port no ipv4 subnets",
			platform: func() *openstack.Platform {
				p := validPlatform()
				fixedIPv6 := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-1111-4465-8d54-3517ec2bad48"},
				}
				p.ControlPlanePort = &openstack.PortTarget{
					FixedIPs: []openstack.FixedIP{fixedIPv6},
				}
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnetv6 := subnets.Subnet{
					ID:        "00000000-1111-4465-8d54-3517ec2bad48",
					CIDR:      "2001:db8::/64",
					IPVersion: 6,
				}
				allsubnets := []*subnets.Subnet{&subnetv6}
				ci.ControlPlanePortSubnets = allsubnets
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("2001:db8::/64"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: `platform.openstack.controlPlanePort.fixedIPs: Internal error: one IPv4 subnet must be specified`,
		},
		{
			name: "MTU too low",
			platform: func() *openstack.Platform {
				p := validPlatform()
				fixedIPv6 := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "00000000-1111-4465-8d54-3517ec2bad48"},
				}
				p.ControlPlanePort = &openstack.PortTarget{
					FixedIPs: []openstack.FixedIP{fixedIPv6},
					Network:  openstack.NetworkFilter{ID: "71b97520-69af-4c35-8153-cdf827z96e60"},
				}
				return p
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validPlatformCloudInfo()
				subnetv6 := subnets.Subnet{
					ID:        "00000000-1111-4465-8d54-3517ec2bad48",
					CIDR:      "2001:db8::/64",
					IPVersion: 6,
					NetworkID: "71b97520-69af-4c35-8153-cdf827z96e60",
				}
				allSubnets := []*subnets.Subnet{&subnetv6}
				ci.ControlPlanePortSubnets = allSubnets

				network := Network{
					networks.Network{
						ID:           "71b97520-69af-4c35-8153-cdf827z96e60",
						Name:         "too-low",
						AdminStateUp: true,
						Status:       "ACTIVE",
						Subnets:      []string{"00000000-1111-4465-8d54-3517ec2bad48"},
					},
					mtu.NetworkMTUExt{MTU: 1200},
				}
				ci.ControlPlanePortNetwork = &network
				return ci
			}(),
			networking: func() *types.Networking {
				n := validNetworking()
				machineNetworkEntry := &types.MachineNetworkEntry{
					CIDR: *ipnet.MustParseCIDR("2001:db8::/64"),
				}
				n.MachineNetwork = []types.MachineNetworkEntry{*machineNetworkEntry}
				return n
			}(),
			expectedErrMsg: "platform.openstack.controlPlanePort.network: Internal error: network should have an MTU of at least 1380",
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
