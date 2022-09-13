package image

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/api/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Unable to test Generate because bootstrap.AddStorageFiles
// returns error in unit test:
//   open data/agent/files: no such file or directory
// Unit test working directory is ./pkg/asset/agent/image
// While normal execution working directory is ./data
// func TestIgnition_Generate(t *testing.T) {}

func TestIgnition_getTemplateData(t *testing.T) {
	clusterImageSet := &hivev1.ClusterImageSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openshift-v4.10.0",
		},
		Spec: hivev1.ClusterImageSetSpec{
			ReleaseImage: "quay.io:443/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64",
		},
	}
	pullSecret := "pull-secret"
	nodeZeroIP := "192.168.111.80"
	agentClusterInstall := &hiveext.AgentClusterInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-agent-cluster-install",
			Namespace: "cluster0",
		},
		Spec: hiveext.AgentClusterInstallSpec{
			APIVIP:       "192.168.111.2",
			SSHPublicKey: "ssh-rsa AAAAmyKey",
			ProvisionRequirements: hiveext.ProvisionRequirements{
				ControlPlaneAgents: 3,
				WorkerAgents:       5,
			},
		},
	}
	releaseImage := "quay.io:443/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64"
	releaseImageMirror := "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"
	infraEnvID := "random-infra-env-id"
	haveMirrorConfig := true

	releaseImageList, err := releaseImageList(clusterImageSet.Spec.ReleaseImage, "x86_64")
	assert.NoError(t, err)
	templateData := getTemplateData(pullSecret, nodeZeroIP, releaseImageList, releaseImage, releaseImageMirror, haveMirrorConfig, agentClusterInstall, infraEnvID)
	assert.Equal(t, "http", templateData.ServiceProtocol)
	assert.Equal(t, "http://"+nodeZeroIP+":8090/", templateData.ServiceBaseURL)
	assert.Equal(t, pullSecret, templateData.PullSecret)
	assert.Equal(t, "", templateData.PullSecretToken)
	assert.Equal(t, nodeZeroIP, templateData.NodeZeroIP)
	assert.Equal(t, nodeZeroIP+":8090", templateData.AssistedServiceHost)
	assert.Equal(t, agentClusterInstall.Spec.APIVIP, templateData.APIVIP)
	assert.Equal(t, agentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents, templateData.ControlPlaneAgents)
	assert.Equal(t, agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents, templateData.WorkerAgents)
	assert.Equal(t, releaseImageList, templateData.ReleaseImages)
	assert.Equal(t, releaseImage, templateData.ReleaseImage)
	assert.Equal(t, releaseImageMirror, templateData.ReleaseImageMirror)
	assert.Equal(t, haveMirrorConfig, templateData.HaveMirrorConfig)
	assert.Equal(t, infraEnvID, templateData.InfraEnvID)
}

func TestIgnition_addStaticNetworkConfig(t *testing.T) {
	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}

	cases := []struct {
		Name                string
		staticNetworkConfig []*models.HostStaticNetworkConfig
		expectedError       string
		expectedFileList    []string
	}{
		{
			Name: "default",
			staticNetworkConfig: []*models.HostStaticNetworkConfig{
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "52:54:01:aa:aa:a1"},
					},
					NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 52:54:01:aa:aa:a1\n  name: eth0\n  state: up\n  type: ethernet\n",
				},
			},
			expectedError: "",
			expectedFileList: []string{
				"/etc/assisted/network/host0/eth0.nmconnection",
				"/etc/assisted/network/host0/mac_interface.ini",
				"/usr/local/bin/pre-network-manager-config.sh",
			},
		},
		{
			Name:                "no-static-network-configs",
			staticNetworkConfig: []*models.HostStaticNetworkConfig{},
			expectedError:       "",
			expectedFileList:    nil,
		},
		{
			Name: "error-processing-config",
			staticNetworkConfig: []*models.HostStaticNetworkConfig{
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "52:54:01:aa:aa:a1"},
					},
					NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: bad-ip\n      prefix-length: 24\n    enabled: true\n  mac-address: 52:54:01:aa:aa:a1\n  name: eth0\n  state: up\n  type: ethernet\n",
				},
			},
			expectedError:    "'bad-ip' does not appear to be an IPv4 or IPv6 address",
			expectedFileList: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			config := igntypes.Config{}
			err := addStaticNetworkConfig(&config, tc.staticNetworkConfig)

			if tc.expectedError != "" {
				assert.Regexp(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			var fileList []string
			for _, file := range config.Storage.Files {
				fileList = append(fileList, file.Node.Path)
			}
			assert.Equal(t, tc.expectedFileList, fileList)
		})
	}
}

func TestRetrieveRendezvousIP(t *testing.T) {
	rawConfig := `interfaces:
  - ipv4:
      address:
        - ip: "192.168.122.21"`
	cases := []struct {
		Name                 string
		agentConfig          *agent.Config
		nmStateConfigs       []*v1beta1.NMStateConfig
		expectedRendezvousIP string
		expectedError        string
	}{
		{
			Name: "valid-agent-config-provided-with-RendezvousIP",
			agentConfig: &agent.Config{
				RendezvousIP: "192.168.122.21",
				Hosts: []agent.Host{
					{
						Hostname: "control-0.example.org",
						Role:     "master",
					},
				},
			},
			expectedRendezvousIP: "192.168.122.21",
		},
		{
			Name: "no-agent-config-provided-so-read-from-nmstateconfig",
			nmStateConfigs: []*v1beta1.NMStateConfig{
				{
					Spec: v1beta1.NMStateConfigSpec{
						NetConfig: v1beta1.NetConfig{
							Raw: []byte(rawConfig),
						},
					},
				},
			},
			expectedRendezvousIP: "192.168.122.21",
		},
		{
			Name: "neither-agent-config-was-provided-with-RendezvousIP-nor-nmstateconfig-manifest",
			agentConfig: &agent.Config{
				Hosts: []agent.Host{
					{
						Hostname: "control-0.example.org",
						Role:     "master",
					},
				},
			},
			expectedError: "missing rendezvousIP in agent-config or at least one NMStateConfig manifest",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			rendezvousIP, err := RetrieveRendezvousIP(tc.agentConfig, tc.nmStateConfigs)
			if tc.expectedError != "" {
				assert.Regexp(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRendezvousIP, rendezvousIP)
			}
		})
	}

}

func TestAddHostConfig_Roles(t *testing.T) {
	cases := []struct {
		Name                            string
		agentConfig                     *agentconfig.AgentConfig
		expectedNumberOfHostConfigFiles int
	}{
		{
			Name: "one-host-role-defined",
			agentConfig: &agentconfig.AgentConfig{
				Config: &agent.Config{
					Hosts: []agent.Host{
						{
							Role: "master",
						},
					},
				},
			},
			expectedNumberOfHostConfigFiles: 1,
		},
		{
			Name: "multiple-host-roles-defined",
			agentConfig: &agentconfig.AgentConfig{
				Config: &agent.Config{
					Hosts: []agent.Host{
						{
							Role: "master",
						},
						{
							Role: "master",
						},
						{
							Role: "master",
						},
						{
							Role: "worker",
						},
						{
							Role: "worker",
						},
					},
				},
			},
			expectedNumberOfHostConfigFiles: 5,
		},
		{
			Name:                            "zero-host-roles-defined",
			expectedNumberOfHostConfigFiles: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			config := &igntypes.Config{}
			err := addHostConfig(config, tc.agentConfig)
			assert.NoError(t, err)
			assert.Equal(t, len(config.Storage.Files), tc.expectedNumberOfHostConfigFiles)
			for _, file := range config.Storage.Files {
				assert.Equal(t, true, strings.HasPrefix(file.Path, "/etc/assisted/hostconfig"))
				assert.Equal(t, true, strings.HasSuffix(file.Path, "role"))
			}
		})
	}

}

func TestIgnition_Generate(t *testing.T) {
	// Generate calls addStaticNetworkConfig which calls nmstatectl
	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}

	// This patch currently allows testing the Ignition asset using the embedded resources.
	// TODO: Replace it by mocking the filesystem in bootstrap.AddStorageFiles()
	workingDirectory, _ := os.Getwd()
	os.Chdir(path.Join(workingDirectory, "../../../../data"))
	secretDataBytes, _ := base64.StdEncoding.DecodeString("super-secret")

	cases := []struct {
		name                                  string
		overrideDeps                          []asset.Asset
		expectedError                         string
		expectedFiles                         []string
		preNetworkManagerConfigServiceEnabled bool
	}{
		{
			name:                                  "no-extra-manifests",
			expectedFiles:                         []string{},
			preNetworkManagerConfigServiceEnabled: true,
		},
		{
			name: "default",
			overrideDeps: []asset.Asset{
				&manifests.ExtraManifests{
					FileList: []*asset.File{
						{
							Filename: "openshift/test-configmap.yaml",
						},
					},
				},
			},
			expectedFiles: []string{
				"/etc/assisted/extra-manifests/test-configmap.yaml",
			},
			preNetworkManagerConfigServiceEnabled: true,
		},
		{
			name: "no nmstateconfigs defined, pre-network-manager-config.service should not be enabled",
			overrideDeps: []asset.Asset{
				&manifests.AgentManifests{
					InfraEnv: &v1beta1.InfraEnv{
						Spec: v1beta1.InfraEnvSpec{
							SSHAuthorizedKey: "my-ssh-key",
						},
					},
					ClusterImageSet: &hivev1.ClusterImageSet{
						Spec: hivev1.ClusterImageSetSpec{
							ReleaseImage: "registry.ci.openshift.org/origin/release:4.11",
						},
					},
					PullSecret: &v1.Secret{
						Data: map[string][]byte{
							".dockerconfigjson": secretDataBytes,
						},
					},
					AgentClusterInstall: &hiveext.AgentClusterInstall{
						Spec: hiveext.AgentClusterInstallSpec{
							APIVIP: "192.168.111.5",
							ProvisionRequirements: hiveext.ProvisionRequirements{
								ControlPlaneAgents: 3,
								WorkerAgents:       5,
							},
						},
					},
				},
			},
			expectedFiles:                         []string{},
			preNetworkManagerConfigServiceEnabled: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			deps := buildIgnitionAssetDefaultDependencies()

			for _, od := range tc.overrideDeps {
				for i, d := range deps {
					if d.Name() == od.Name() {
						deps[i] = od
						break
					}
				}
			}

			parents := asset.Parents{}
			parents.Add(deps...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assert.Len(t, ignitionAsset.Config.Storage.Directories, 1)
				assert.Equal(t, "/etc/assisted/extra-manifests", ignitionAsset.Config.Storage.Directories[0].Node.Path)

				for _, f := range tc.expectedFiles {
					found := false
					for _, i := range ignitionAsset.Config.Storage.Files {
						if i.Node.Path == f {
							found = true
							break
						}
					}
					assert.True(t, found, fmt.Sprintf("Expected file %s not found", f))
				}

				for _, unit := range ignitionAsset.Config.Systemd.Units {
					if unit.Name == "pre-network-manager-config.service" {
						if unit.Enabled == nil {
							assert.Equal(t, tc.preNetworkManagerConfigServiceEnabled, false)
						} else {
							assert.Equal(t, tc.preNetworkManagerConfigServiceEnabled, *unit.Enabled)
						}
					}
				}
			}
		})
	}
}

// This test util create the minimum valid set of dependencies for the
// Ignition asset
func buildIgnitionAssetDefaultDependencies() []asset.Asset {
	secretDataBytes, _ := base64.StdEncoding.DecodeString("super-secret")

	return []asset.Asset{
		&manifests.AgentManifests{
			InfraEnv: &v1beta1.InfraEnv{
				Spec: v1beta1.InfraEnvSpec{
					SSHAuthorizedKey: "my-ssh-key",
				},
			},
			ClusterImageSet: &hivev1.ClusterImageSet{
				Spec: hivev1.ClusterImageSetSpec{
					ReleaseImage: "registry.ci.openshift.org/origin/release:4.11",
				},
			},
			PullSecret: &v1.Secret{
				Data: map[string][]byte{
					".dockerconfigjson": secretDataBytes,
				},
			},
			AgentClusterInstall: &hiveext.AgentClusterInstall{
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP: "192.168.111.5",
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       5,
					},
				},
			},
			NMStateConfigs: []*aiv1beta1.NMStateConfig{
				{
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "eth0",
								MacAddress: "00:01:02:03:04:05",
							},
						},
					},
				},
			},
			StaticNetworkConfigs: []*models.HostStaticNetworkConfig{
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "00:01:02:03:04:05"},
					},
					NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 00:01:02:03:04:05\n  name: eth0\n  state: up\n  type: ethernet\n",
				},
			},
		},
		&agentconfig.AgentConfig{
			Config: &agent.Config{
				RendezvousIP: "192.168.111.80",
			},
			File: &asset.File{
				Filename: "/cluster-manifests/agent-config.yaml",
			},
		},
		&manifests.ExtraManifests{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
		&tls.AdminKubeConfigClientCertKey{},
	}
}

func TestIgnition_getMirrorFromRelease(t *testing.T) {

	cases := []struct {
		name           string
		release        string
		registriesConf mirror.RegistriesConf
		expectedMirror string
	}{
		{
			name:           "no-mirror",
			release:        "registry.ci.openshift.org/ocp/release:latest",
			registriesConf: mirror.RegistriesConf{},
			expectedMirror: "",
		},
		{
			name:    "mirror-no-match",
			release: "registry.ci.openshift.org/ocp/release:4.11.0-0.nightly-foo",
			registriesConf: mirror.RegistriesConf{
				File: &asset.File{
					Filename: "registries.conf",
					Data:     []byte(""),
				},
				MirrorConfig: []mirror.RegistriesConfig{
					{
						Location: "some.registry.org/release",
						Mirror:   "some.mirror.org",
					},
				},
			},
			expectedMirror: "",
		},
		{
			name:    "mirror-match",
			release: "registry.ci.openshift.org/ocp/release:4.11.0-0.nightly-foo",
			registriesConf: mirror.RegistriesConf{
				File: &asset.File{
					Filename: "registries.conf",
					Data:     []byte(""),
				},
				MirrorConfig: []mirror.RegistriesConfig{
					{
						Location: "registry.ci.openshift.org/ocp/release",
						Mirror:   "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image",
					},
				},
			},
			expectedMirror: "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image:4.11.0-0.nightly-foo",
		},
		{
			name:    "mirror-match-with-checksum",
			release: "quay.io/openshift-release-dev/ocp-release@sha256:300bce8246cf880e792e106607925de0a404484637627edf5f517375517d54a4",
			registriesConf: mirror.RegistriesConf{
				File: &asset.File{
					Filename: "registries.conf",
					Data:     []byte(""),
				},
				MirrorConfig: []mirror.RegistriesConfig{
					{
						Location: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
						Mirror:   "localhost:5000/openshift4/openshift/release",
					},
					{
						Location: "quay.io/openshift-release-dev/ocp-release",
						Mirror:   "localhost:5000/openshift-release-dev/ocp-release",
					},
				},
			},
			expectedMirror: "localhost:5000/openshift-release-dev/ocp-release@sha256:300bce8246cf880e792e106607925de0a404484637627edf5f517375517d54a4",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mirror := getMirrorFromRelease(tc.release, &tc.registriesConf)

			assert.Equal(t, tc.expectedMirror, mirror)

		})
	}
}
