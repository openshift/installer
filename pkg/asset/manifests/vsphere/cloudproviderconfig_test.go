package vsphere

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	vsphere "github.com/openshift/installer/pkg/types/vsphere"
)

var (
	expectedIniConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[VirtualCenter "test-vcenter"]
port = "443"

datacenters = "test-datacenter,test-datacenter2"

[Workspace]
server = "test-vcenter"
datacenter = "test-datacenter"
default-datastore = "test-datastore"
folder = "/test-datacenter/vm/test-folder"
resourcepool-path = "/test-datacenter/host/cluster/Resources/test-resourcepool"

`

	expectIniLabelsSection = `[Labels]
region = "openshift-region"
zone = "openshift-zone"
`

	// expectedYamlConfig is used by the basic yaml path test (feature gate
	// enabled via inDefault, machine network 10.0.0.0/24 → nodes section
	// populated with that CIDR).
	expectedYamlConfig = `global:
  insecureFlag: true
  secretName: vsphere-creds
  secretNamespace: kube-system
labels:
  region: openshift-region
  zone: openshift-zone
nodes:
  externalNetworkSubnetCidr: 10.0.0.0/24
  internalNetworkSubnetCidr: 10.0.0.0/24
vcenter:
  test-vcenter:
    datacenters:
    - test-datacenter
    - test-datacenter2
    port: 443
    server: test-vcenter
`

	// expectedYamlConfigWithNodes is used by tests that exercise the yaml path
	// with the VSphereMultiNetworks feature gate enabled and a machine network
	// CIDR of 10.0.0.0/24 (the fallback path).
	expectedYamlConfigWithNodes = `global:
  insecureFlag: true
  secretName: vsphere-creds
  secretNamespace: kube-system
labels:
  region: openshift-region
  zone: openshift-zone
nodes:
  externalNetworkSubnetCidr: 10.0.0.0/24
  internalNetworkSubnetCidr: 10.0.0.0/24
vcenter:
  test-vcenter:
    datacenters:
    - test-datacenter
    - test-datacenter2
    port: 443
    server: test-vcenter
`

	// expectedYamlConfigWithExplicitNodes is used by the test that sets
	// NodeNetworking explicitly in the install-config.
	expectedYamlConfigWithExplicitNodes = `global:
  insecureFlag: true
  secretName: vsphere-creds
  secretNamespace: kube-system
labels:
  region: openshift-region
  zone: openshift-zone
nodes:
  externalNetworkSubnetCidr: 192.168.10.0/24
  internalNetworkSubnetCidr: 192.168.20.0/24
vcenter:
  test-vcenter:
    datacenters:
    - test-datacenter
    - test-datacenter2
    port: 443
    server: test-vcenter
`
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"test-datacenter",
					"test-datacenter2",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-dz-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore",
					ResourcePool: "/test-datacenter/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-dz-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter2",
					ComputeCluster: "/test-datacenter2/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore2",
					ResourcePool: "/test-datacenter2/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter2/vm/test-folder",
				},
			},
			{
				Name:   "test-dz-east-3a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter2",
					ComputeCluster: "/test-datacenter2/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore2",
					ResourcePool: "/test-datacenter2/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter2/vm/test-folder",
				},
			},
		},
	}
}

// makeInstallConfig builds a minimal installconfig.InstallConfig for use in
// tests. featureSet controls which feature gates are active.
func makeInstallConfig(platform *vsphere.Platform, featureSet configv1.FeatureSet, machineNetworkCIDR string) *installconfig.InstallConfig {
	_, cidr, err := net.ParseCIDR(machineNetworkCIDR)
	if err != nil {
		panic(err)
	}
	cfg := &types.InstallConfig{
		FeatureSet: featureSet,
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: ipnet.IPNet{IPNet: *cidr}},
			},
		},
		Platform: types.Platform{
			VSphere: platform,
		},
	}
	return installconfig.MakeAsset(cfg)
}

func TestCloudProviderConfig(t *testing.T) {
	cases := []struct {
		name                string
		platform            *vsphere.Platform
		cloudProviderFunc   func(string, *vsphere.Platform) (string, error)
		installConfig       *installconfig.InstallConfig
		useYaml             bool
		expectedCloudConfig string
	}{
		{
			name:                "valid intree cloud provider config",
			platform:            validPlatform(),
			cloudProviderFunc:   CloudProviderConfigIni,
			expectedCloudConfig: expectedIniConfig + expectIniLabelsSection,
		},
		{
			name: "valid single failure domain intree cloud provider config",
			platform: func() *vsphere.Platform {
				p := validPlatform()

				p.FailureDomains = p.FailureDomains[0:1]
				p.VCenters[0].Datacenters = p.VCenters[0].Datacenters[0:1]

				return p
			}(),
			cloudProviderFunc: CloudProviderConfigIni,
			expectedCloudConfig: func() string {
				// only a single datacenter would be provided to the datacenters
				ini := strings.ReplaceAll(expectedIniConfig, ",test-datacenter2", "")
				return ini
			}(),
		},
		{
			// VSphereMultiNetworks is enabled in inDefault() so the nodes
			// section is always populated from the machine network CIDR when
			// nodeNetworking is not explicitly set.
			name:                "valid out of tree yaml cloud provider config",
			platform:            validPlatform(),
			installConfig:       makeInstallConfig(validPlatform(), "", "10.0.0.0/24"),
			useYaml:             true,
			expectedCloudConfig: expectedYamlConfig,
		},
		{
			// TechPreviewNoUpgrade also enables VSphereMultiNetworks – same
			// result as the default case above.
			name:                "valid out of tree yaml cloud provider config with node networking from machine network",
			platform:            validPlatform(),
			installConfig:       makeInstallConfig(validPlatform(), configv1.TechPreviewNoUpgrade, "10.0.0.0/24"),
			useYaml:             true,
			expectedCloudConfig: expectedYamlConfigWithNodes,
		},
		{
			// With VSphereMultiNetworks enabled and explicit nodeNetworking set
			// in the install-config, the explicit values take precedence.
			name: "valid out of tree yaml cloud provider config with explicit node networking",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.NodeNetworking = &configv1.VSpherePlatformNodeNetworking{
					External: configv1.VSpherePlatformNodeNetworkingSpec{
						NetworkSubnetCIDR: []string{"192.168.10.0/24"},
					},
					Internal: configv1.VSpherePlatformNodeNetworkingSpec{
						NetworkSubnetCIDR: []string{"192.168.20.0/24"},
					},
				}
				return p
			}(),
			installConfig: func() *installconfig.InstallConfig {
				p := validPlatform()
				p.NodeNetworking = &configv1.VSpherePlatformNodeNetworking{
					External: configv1.VSpherePlatformNodeNetworkingSpec{
						NetworkSubnetCIDR: []string{"192.168.10.0/24"},
					},
					Internal: configv1.VSpherePlatformNodeNetworkingSpec{
						NetworkSubnetCIDR: []string{"192.168.20.0/24"},
					},
				}
				return makeInstallConfig(p, configv1.TechPreviewNoUpgrade, "10.0.0.0/24")
			}(),
			useYaml:             true,
			expectedCloudConfig: expectedYamlConfigWithExplicitNodes,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cloudConfig string
			var err error
			if tc.useYaml {
				cloudConfig, err = CloudProviderConfigYaml("infraID", tc.installConfig)
			} else {
				cloudConfig, err = tc.cloudProviderFunc("infraID", tc.platform)
			}
			assert.NoError(t, err, "failed to create cloud provider config")
			assert.Equal(t, tc.expectedCloudConfig, cloudConfig, "unexpected cloud provider config")
		})
	}
}
