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
	"github.com/vincent-petithory/dataurl"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/password"
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
	publicContainerRegistries := "quay.io,registry.ci.openshift.org"

	releaseImageList, err := releaseImageList(clusterImageSet.Spec.ReleaseImage, "x86_64")
	assert.NoError(t, err)

	arch := "x86_64"
	ov := "4.12"
	isoURL := "https://rhcos.mirror.openshift.com/art/storage/releases/rhcos-4.12/412.86.202208101039-0/x86_64/rhcos-412.86.202208101039-0-live.x86_64.iso"
	ver := "412.86.202208101039-0"
	osImage := &models.OsImage{
		CPUArchitecture:  &arch,
		OpenshiftVersion: &ov,
		URL:              &isoURL,
		Version:          &ver,
	}

	proxy := &aiv1beta1.Proxy{
		HTTPProxy:  "http://1.1.1.1:80",
		HTTPSProxy: "https://1.1.1.1:443",
		NoProxy:    "valid-proxy.com,172.30.0.0/16",
	}
	clusterName := "test-agent-cluster-install.test"

	templateData := getTemplateData(clusterName, pullSecret, nodeZeroIP, releaseImageList, releaseImage, releaseImageMirror, haveMirrorConfig, publicContainerRegistries, agentClusterInstall, infraEnvID, osImage, proxy)
	assert.Equal(t, clusterName, templateData.ClusterName)
	assert.Equal(t, "http", templateData.ServiceProtocol)
	assert.Equal(t, "http://"+nodeZeroIP+":8090/", templateData.ServiceBaseURL)
	assert.Equal(t, pullSecret, templateData.PullSecret)
	assert.Equal(t, nodeZeroIP, templateData.NodeZeroIP)
	assert.Equal(t, nodeZeroIP+":8090", templateData.AssistedServiceHost)
	assert.Equal(t, agentClusterInstall.Spec.APIVIP, templateData.APIVIP)
	assert.Equal(t, agentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents, templateData.ControlPlaneAgents)
	assert.Equal(t, agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents, templateData.WorkerAgents)
	assert.Equal(t, releaseImageList, templateData.ReleaseImages)
	assert.Equal(t, releaseImage, templateData.ReleaseImage)
	assert.Equal(t, releaseImageMirror, templateData.ReleaseImageMirror)
	assert.Equal(t, haveMirrorConfig, templateData.HaveMirrorConfig)
	assert.Equal(t, publicContainerRegistries, templateData.PublicContainerRegistries)
	assert.Equal(t, infraEnvID, templateData.InfraEnvID)
	assert.Equal(t, osImage, templateData.OSImage)
	assert.Equal(t, proxy, templateData.Proxy)
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
		nmStateConfigs       []*aiv1beta1.NMStateConfig
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
			nmStateConfigs: []*aiv1beta1.NMStateConfig{
				{
					Spec: aiv1beta1.NMStateConfigSpec{
						NetConfig: aiv1beta1.NetConfig{
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

func defaultGeneratedFiles() []string {
	return []string{
		"/etc/issue",
		"/etc/multipath.conf",
		"/etc/containers/containers.conf",
		"/etc/motd",
		"/root/.docker/config.json",
		"/root/assisted.te",
		"/usr/local/bin/common.sh",
		"/usr/local/bin/agent-gather",
		"/usr/local/bin/extract-agent.sh",
		"/usr/local/bin/get-container-images.sh",
		"/usr/local/bin/set-hostname.sh",
		"/usr/local/bin/start-agent.sh",
		"/usr/local/bin/start-cluster-installation.sh",
		"/usr/local/bin/wait-for-assisted-service.sh",
		"/usr/local/bin/set-node-zero.sh",
		"/usr/local/share/assisted-service/assisted-db.env",
		"/usr/local/share/assisted-service/assisted-service.env",
		"/usr/local/share/assisted-service/images.env",
		"/usr/local/bin/bootstrap-service-record.sh",
		"/usr/local/bin/release-image.sh",
		"/usr/local/bin/release-image-download.sh",
		"/etc/assisted/manifests/agent-config.yaml",
		"/etc/assisted/network/host0/eth0.nmconnection",
		"/etc/assisted/network/host0/mac_interface.ini",
		"/usr/local/bin/pre-network-manager-config.sh",
		"/opt/agent/tls/kubeadmin-password.hash",
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
		expectedFileContent                   map[string]string
		preNetworkManagerConfigServiceEnabled bool
	}{
		{
			name:                                  "no-extra-manifests",
			preNetworkManagerConfigServiceEnabled: true,

			expectedFiles: defaultGeneratedFiles(),
		},
		{
			name: "default",
			overrideDeps: []asset.Asset{
				&manifests.ExtraManifests{
					FileList: []*asset.File{
						{
							Filename: "openshift/test-configmap.yaml",
							Data: []byte(`
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: agent-test-1

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: agent-test-2`),
						},
					},
				},
			},
			expectedFiles: []string{
				"/etc/issue",
				"/etc/multipath.conf",
				"/etc/containers/containers.conf",
				"/etc/motd",
				"/root/.docker/config.json",
				"/root/assisted.te",
				"/usr/local/bin/common.sh",
				"/usr/local/bin/agent-gather",
				"/usr/local/bin/extract-agent.sh",
				"/usr/local/bin/get-container-images.sh",
				"/usr/local/bin/install-status.sh",
				"/usr/local/bin/issue_status.sh",
				"/usr/local/bin/set-hostname.sh",
				"/usr/local/bin/start-agent.sh",
				"/usr/local/bin/start-cluster-installation.sh",
				"/usr/local/bin/wait-for-assisted-service.sh",
				"/usr/local/bin/set-node-zero.sh",
				"/usr/local/share/assisted-service/assisted-db.env",
				"/usr/local/share/assisted-service/assisted-service.env",
				"/usr/local/share/assisted-service/images.env",
				"/usr/local/bin/bootstrap-service-record.sh",
				"/usr/local/bin/release-image.sh",
				"/usr/local/bin/release-image-download.sh",
				"/etc/assisted/manifests/agent-config.yaml",
				"/etc/assisted/network/host0/eth0.nmconnection",
				"/etc/assisted/network/host0/mac_interface.ini",
				"/usr/local/bin/pre-network-manager-config.sh",
				"/opt/agent/tls/kubeadmin-password.hash",
				"/etc/assisted/extra-manifests/test-configmap-0.yaml",
				"/etc/assisted/extra-manifests/test-configmap-1.yaml",
			},
			expectedFileContent: map[string]string{
				"/usr/local/share/assisted-service/images.env": `ASSISTED_SERVICE_HOST=192.168.111.80:8090
ASSISTED_SERVICE_SCHEME=http
OS_IMAGES=\[\{"openshift_version":"was not built correctly","cpu_architecture":"x86_64","url":"https://rhcos.mirror.openshift.com/art/storage/releases/rhcos-.*.x86_64.iso","version":".*"\}\]
`,
			},
			preNetworkManagerConfigServiceEnabled: true,
		},
		{
			name: "no nmstateconfigs defined, pre-network-manager-config.service should not be enabled",
			overrideDeps: []asset.Asset{
				&manifests.AgentManifests{
					InfraEnv: &aiv1beta1.InfraEnv{
						Spec: aiv1beta1.InfraEnvSpec{
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

				if len(tc.expectedFiles) > 0 {
					assert.Equal(t, len(tc.expectedFiles), len(ignitionAsset.Config.Storage.Files))

					for _, f := range tc.expectedFiles {
						found := false
						for _, i := range ignitionAsset.Config.Storage.Files {
							if i.Node.Path == f {
								if expectedData, ok := tc.expectedFileContent[i.Node.Path]; ok {
									actualData, err := dataurl.DecodeString(*i.FileEmbedded1.Contents.Source)
									assert.NoError(t, err)
									assert.Regexp(t, expectedData, string(actualData.Data))
								}

								found = true
								break
							}
						}
						assert.True(t, found, fmt.Sprintf("Expected file %s not found", f))
					}
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
			InfraEnv: &aiv1beta1.InfraEnv{
				Spec: aiv1beta1.InfraEnvSpec{
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
		&password.KubeadminPassword{},
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

			mirror := mirror.GetMirrorFromRelease(tc.release, &tc.registriesConf)

			assert.Equal(t, tc.expectedMirror, mirror)

		})
	}
}

func TestIgnition_getPublicContainerRegistries(t *testing.T) {

	cases := []struct {
		name               string
		registriesConf     mirror.RegistriesConf
		expectedRegistries string
	}{
		{
			name:               "no-mirror",
			registriesConf:     mirror.RegistriesConf{},
			expectedRegistries: "quay.io",
		},
		{
			name: "mirror-one-entry",
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
			expectedRegistries: "some.registry.org",
		},
		{
			name: "mirror-multiple-entries",
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
						Location: "registry.ci.openshift.org/ocp/release",
						Mirror:   "localhost:5000/openshift-release-dev/ocp-release",
					},
				},
			},
			expectedRegistries: "quay.io,registry.ci.openshift.org",
		},
		{
			name: "duplicate-entries",
			registriesConf: mirror.RegistriesConf{
				File: &asset.File{
					Filename: "registries.conf",
					Data:     []byte(""),
				},
				MirrorConfig: []mirror.RegistriesConfig{
					{
						Location: "registry.ci.openshift.org/ocp-v4.0-art-dev",
						Mirror:   "localhost:5000/openshift4/openshift/release",
					},
					{
						Location: "registry.ci.openshift.org/ocp/release",
						Mirror:   "localhost:5000/openshift-release-dev/ocp-release",
					},
				},
			},
			expectedRegistries: "registry.ci.openshift.org",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			publicContainerRegistries := getPublicContainerRegistries(&tc.registriesConf)

			assert.Equal(t, tc.expectedRegistries, publicContainerRegistries)

		})
	}
}
