package image

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/imagebased/configimage"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/imagebased"
	"github.com/openshift/installer/pkg/types/none"
)

var (
	testCert = `-----BEGIN CERTIFICATE-----
MIICYTCCAcqgAwIBAgIJAI2kA+uXAbhOMA0GCSqGSIb3DQEBCwUAMEgxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJDQTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzEUMBIG
A1UECgwLUmVkIEhhdCBJbmMwHhcNMTkwMjEyMTkzMjUzWhcNMTkwMjEzMTkzMjUz
WjBIMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFu
Y2lzY28xFDASBgNVBAoMC1JlZCBIYXQgSW5jMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQC+HOC0mKig/oINAKPo88LqxDJ4l7lozdLtp5oGeqWrLUXSfkvXAkQY
2QYdvPAjpRfH7Ii7G0Asx+HTKdvula7B5fXDjc6NYKuEpTJZRV1ugntI97bozF/E
C2BBmxxEnJN3+Xe8RYXMjz5Q4aqPw9vZhlWN+0hrREl1Ea/zHuWFIQIDAQABo1Mw
UTAdBgNVHQ4EFgQUvTS1XjlvOdsufSyWxukyQu3LriEwHwYDVR0jBBgwFoAUvTS1
XjlvOdsufSyWxukyQu3LriEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsF
AAOBgQB9gFcOXnzJrM65QqxeCB9Z5l5JMjp45UFC9Bj2cgwDHP80Zvi4omlaacC6
aavmnLd67zm9PbYDWRaOIWAMeB916Iwaw/v6I0jwhAk/VxX5Fl6cGlZu9jZ3zbFE
2sDqkwzIuSjCG2A23s6d4M1S3IXCCydoCSLMu+WhLkbboK6jEg==
-----END CERTIFICATE-----
`
)

func TestIgnition_Generate(t *testing.T) {
	defer setupEmbeddedResources(t)()

	registriesConf := &RegistriesConf{
		Config: &sysregistriesv2.V2RegistriesConf{
			Registries: []sysregistriesv2.Registry{
				{
					Endpoint:           sysregistriesv2.Endpoint{Location: "quay.io"},
					MirrorByDigestOnly: true,
					Mirrors:            []sysregistriesv2.Endpoint{{Location: "mirror-quay.io"}},
				},
			},
		},
	}
	registriesConfData, err := toml.Marshal(registriesConf.Config)
	assert.NoError(t, err)

	testNetworkConfig := "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 00:01:02:03:04:05\n  name: eth0\n  state: up\n  type: ethernet\n\n"

	ignitionConfigOverride := `{"ignition": {"version": "3.2.0"}, "storage": {"files": [{"path": "/tmp/example", "contents": {"source": "data:text/plain;charset=utf-8;base64,aGVscGltdHJhcHBlZGluYXN3YWdnZXJzcGVj"}}]}}`

	postDeploymentScript := "#!/bin/bash\necho \"Post image-based installation script\""

	cases := []struct {
		name         string
		dependencies []asset.Asset

		expectedError        string
		expectedFiles        map[string]string
		expectedSystemdUnits map[string]bool
	}{
		{
			name: "comprehensive configuration",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						imageDigestSources([]types.ImageDigestSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						additionalTrustBundle(testCert).
						networkConfig(&aiv1beta1.NetConfig{Raw: []byte(testNetworkConfig)}).
						ignitionConfigOverride(ignitionConfigOverride).
						build(),
				},
				&RegistriesConf{Data: registriesConfData},
				&PostDeployment{File: &asset.File{
					Filename: "post.sh",
					Data:     []byte(postDeploymentScript),
				}},
				&configimage.InstallConfig{
					AssetBase: installconfig.AssetBase{
						Config: &types.InstallConfig{
							Platform: types.Platform{None: &none.Platform{}},
						},
					},
				},
			},

			expectedFiles: map[string]string{
				"/var/tmp/pull-secret.json": "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}\n",

				"/usr/local/bin/install-rhcos-and-restore-seed.sh": `#!/bin/bash

set -e # Halt on error

seed_image=${1:-$SEED_IMAGE}
authfile=${PULL_SECRET_FILE:-"/var/tmp/pull-secret.json"}
ibi_config=${IBI_CONFIGURATION_FILE:-"/var/tmp/ibi-configuration.json"}

# Copy the lca-cli binary to the host, pulling the seed image can sometimes fail
until podman create --authfile "${authfile}" --name lca-cli "${seed_image}" lca-cli ; do
    sleep 10
done
podman cp lca-cli:lca-cli /usr/local/bin/lca-cli
podman rm lca-cli

/usr/local/bin/lca-cli ibi -f "${ibi_config}"
`,

				"/var/tmp/ibi-configuration.json": "{\"extraPartitionLabel\":\"var-lib-containers\",\"extraPartitionNumber\":5,\"extraPartitionStart\":\"-40G\",\"installationDisk\":\"/dev/vda\",\"releaseRegistry\":\"mirror.quay.io\",\"seedImage\":\"quay.io/openshift-kni/seed-image:4.16.0\",\"seedVersion\":\"4.16.0\"}\n",

				"/etc/containers/registries.conf": "credential-helpers = []\nshort-name-mode = \"\"\nunqualified-search-registries = []\n\n[[registry]]\n  location = \"quay.io\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"mirror-quay.io\"\n",

				"/etc/pki/ca-trust/source/anchors/additional-trust-bundle.pem": testCert,

				"/var/tmp/network-config.yaml": testNetworkConfig + "\n",

				"/tmp/example": "helpimtrappedinaswaggerspec",

				"/var/tmp/post.sh": "#!/bin/bash\necho \"Post image-based installation script\"",
			},
			expectedSystemdUnits: map[string]bool{
				"install-rhcos-and-restore-seed.service": true,

				"network-config.service": true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assertFiles(t, ignitionAsset.Config, tc.expectedFiles)

				assertSystemdUnits(t, ignitionAsset.Config, tc.expectedSystemdUnits)
			}
		})
	}
}

// TestGetDHCPKernelArgs tests the dual-stack kernel argument logic.
func TestGetDHCPKernelArgs(t *testing.T) {
	cases := []struct {
		name           string
		machineNetwork []types.MachineNetworkEntry
		networkConfig  string
		userArgs       []string
		expectedArgs   string
	}{
		{
			name:         "no machine networks",
			expectedArgs: "",
		},
		{
			name: "single IPv4 machine network, no networkConfig (DHCP)",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
			},
			expectedArgs: "", // Only add args for dual-stack
		},
		{
			name: "single IPv6 machine network, no networkConfig (DHCP)",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			expectedArgs: "", // Only add args for dual-stack
		},
		{
			name: "dual-stack machine networks, no networkConfig (both DHCP)",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			expectedArgs: "ip=dhcp,dhcp6",
		},
		{
			name: "dual-stack machine networks, both static",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  type: ethernet
  state: up
  ipv4:
    enabled: true
    address:
    - ip: 10.0.0.100
      prefix-length: 24
    dhcp: false
  ipv6:
    enabled: true
    address:
    - ip: 2001:db8::100
      prefix-length: 64
    dhcp: false`,
			expectedArgs: "",
		},
		{
			name: "dual-stack machine networks, IPv4 static, IPv6 DHCP",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  type: ethernet
  state: up
  ipv4:
    enabled: true
    address:
    - ip: 10.0.0.100
      prefix-length: 24
    dhcp: false
  ipv6:
    enabled: true
    dhcp: true`,
			expectedArgs: "ip=dhcp6",
		},
		{
			name: "dual-stack machine networks, IPv4 DHCP, IPv6 static",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  type: ethernet
  state: up
  ipv4:
    enabled: true
    dhcp: true
  ipv6:
    enabled: true
    address:
    - ip: 2001:db8::100
      prefix-length: 64
    dhcp: false`,
			expectedArgs: "ip=dhcp",
		},
		{
			name: "dual-stack machine networks, both DHCP explicitly configured",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  type: ethernet
  state: up
  ipv4:
    enabled: true
    dhcp: true
  ipv6:
    enabled: true
    dhcp: true`,
			expectedArgs: "ip=dhcp,dhcp6",
		},
		{
			name: "dual-stack, user already provided kernel args",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			userArgs:     []string{"--append-karg", "ip=dhcp"},
			expectedArgs: "",
		},
		{
			name: "dual-stack, user provided different kernel args",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			userArgs:     []string{"--append-karg", "console=ttyS0"},
			expectedArgs: "ip=dhcp,dhcp6",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create install config with machine networks using the builder pattern from existing tests
			installConfig := &configimage.InstallConfig{
				AssetBase: installconfig.AssetBase{
					Config: &types.InstallConfig{
						Platform: types.Platform{None: &none.Platform{}},
					},
				},
			}
			if len(tc.machineNetwork) > 0 {
				installConfig.Config.Networking = &types.Networking{
					MachineNetwork: tc.machineNetwork,
				}
			}

			// Create IBI config with network config and user args
			ibiConfig := &imagebased.InstallationConfig{
				CoreosInstallerArgs: tc.userArgs,
			}
			if tc.networkConfig != "" {
				ibiConfig.NetworkConfig = &aiv1beta1.NetConfig{
					Raw: []byte(tc.networkConfig),
				}
			}

			result := getDHCPKernelArgs(ibiConfig, installConfig)
			assert.Equal(t, tc.expectedArgs, result, "unexpected kernel args")
		})
	}
}

// setupEmbeddedResources changes to the data directory for testing embedded resources.
func setupEmbeddedResources(t *testing.T) func() {
	t.Helper()
	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(path.Join(workingDirectory, "../../../../data")))
	return func() {
		assert.NoError(t, os.Chdir(workingDirectory))
	}
}

// TestIgnition_Generate_DualStack tests ignition generation with dual-stack configurations.
func TestIgnition_Generate_DualStack(t *testing.T) {
	defer setupEmbeddedResources(t)()

	cases := []struct {
		name           string
		machineNetwork []types.MachineNetworkEntry
		networkConfig  string
		userArgs       []string
		expectedKArgs  []string // Expected kernel arguments in CoreOS installer args
	}{
		{
			name: "dual-stack with dynamic networking",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			expectedKArgs: []string{"--append-karg", "ip=dhcp,dhcp6"},
		},
		{
			name: "dual-stack with static networking",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  ipv4:
    enabled: true
    address:
    - ip: 10.0.0.100
      prefix-length: 24
    dhcp: false
  ipv6:
    enabled: true
    address:
    - ip: 2001:db8::100
      prefix-length: 64
    dhcp: false`,
			expectedKArgs: []string{}, // No kernel args for static
		},
		{
			name: "dual-stack with mixed static/dynamic networking",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			networkConfig: `interfaces:
- name: eth0
  ipv4:
    enabled: true
    dhcp: true
  ipv6:
    enabled: true
    address:
    - ip: 2001:db8::100
      prefix-length: 64
    dhcp: false`,
			expectedKArgs: []string{"--append-karg", "ip=dhcp"},
		},
		{
			name: "single-stack IPv4 with dynamic networking",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
			},
			expectedKArgs: []string{}, // No kernel args for single-stack
		},
		{
			name: "dual-stack with user-provided kernel args",
			machineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
				{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
			},
			userArgs:      []string{"--append-karg", "ip=static"},
			expectedKArgs: []string{"--append-karg", "ip=static"}, // User args preserved
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create install config
			installConfig := &configimage.InstallConfig{
				AssetBase: installconfig.AssetBase{
					Config: &types.InstallConfig{
						Platform: types.Platform{None: &none.Platform{}},
					},
				},
			}
			if len(tc.machineNetwork) > 0 {
				installConfig.Config.Networking = &types.Networking{
					MachineNetwork: tc.machineNetwork,
				}
			}

			// Create IBI config
			ibiConfig := &imagebased.InstallationConfig{
				SeedImage:           "quay.io/openshift-kni/seed-image:4.16.0",
				SeedVersion:         "4.16.0",
				InstallationDisk:    "/dev/vda",
				PullSecret:          `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
				CoreosInstallerArgs: tc.userArgs,
			}
			if tc.networkConfig != "" {
				ibiConfig.NetworkConfig = &aiv1beta1.NetConfig{
					Raw: []byte(tc.networkConfig),
				}
			}

			dependencies := []asset.Asset{
				&ImageBasedInstallationConfig{Config: ibiConfig},
				installConfig,
				&RegistriesConf{},
				&PostDeployment{},
			}

			parents := asset.Parents{}
			parents.Add(dependencies...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(context.Background(), parents)
			assert.NoError(t, err)

			// Check CoreOS installer args contain expected kernel arguments
			foundArgs := false
			for i := 0; i < len(ignitionAsset.Config.Storage.Files); i++ {
				file := ignitionAsset.Config.Storage.Files[i]
				if file.Node.Path == "/var/tmp/ibi-configuration.json" {
					actualData, err := dataurl.DecodeString(*file.FileEmbedded1.Contents.Source)
					assert.NoError(t, err)

					err = fmt.Errorf("found config file")
					if err != nil {
						// Parse and check the CoreOS installer args
						if len(tc.expectedKArgs) > 0 {
							// Should contain the expected kernel args
							configStr := string(actualData.Data)
							for _, expectedArg := range tc.expectedKArgs {
								assert.Contains(t, configStr, expectedArg, "missing expected kernel arg")
							}
						}
						foundArgs = true
						break
					}
				}
			}

			if len(tc.expectedKArgs) > 0 {
				assert.True(t, foundArgs, "expected to find configuration file")
			}
		})
	}
}

func assertFiles(t *testing.T, config *igntypes.Config, files map[string]string) {
	t.Helper()

	if len(files) > 0 {
		assert.Equal(t, len(files), len(config.Storage.Files))

		for name, content := range files {
			found := false
			for _, i := range config.Storage.Files {
				if i.Node.Path == name {
					actualData, err := dataurl.DecodeString(*i.FileEmbedded1.Contents.Source)
					assert.NoError(t, err)
					assert.Equal(t, content, string(actualData.Data))

					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("expected file %s not found", name))
		}
	}
}

func assertSystemdUnits(t *testing.T, config *igntypes.Config, units map[string]bool) {
	t.Helper()

	for name, enabled := range units {
		for _, unit := range config.Systemd.Units {
			if unit.Name != name {
				continue
			}

			if unit.Enabled == nil {
				assert.Equal(t, enabled, false)
			} else {
				assert.Equal(t, enabled, *unit.Enabled)
			}
		}
	}
}
