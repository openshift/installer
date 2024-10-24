package image

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/gencrypto"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types/agent"
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

	releaseImageList, err := releaseImageList(clusterImageSet.Spec.ReleaseImage, "x86_64", []string{"86_64"})
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

	publicKey := "-----BEGIN EC PUBLIC KEY-----\nMHcCAQEEIOSCfDNmx0qe6dncV4tg==\n-----END EC PUBLIC KEY-----\n"
	token := "someToken"
	rendezvousIP := "192.168.111.80"
	templateData := getTemplateData(clusterName, pullSecret, releaseImageList, releaseImage, releaseImageMirror, publicContainerRegistries, "minimal-iso", infraEnvID, publicKey, gencrypto.AuthType, token, "", "", haveMirrorConfig, agentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents, agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents, osImage, proxy, rendezvousIP)
	assert.Equal(t, clusterName, templateData.ClusterName)
	assert.Equal(t, "http", templateData.ServiceProtocol)
	assert.Equal(t, pullSecret, templateData.PullSecret)
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
	assert.Equal(t, publicKey, templateData.PublicKeyPEM)
	assert.Equal(t, gencrypto.AuthType, templateData.AuthType)
	assert.Equal(t, token, templateData.Token)
	assert.Equal(t, rendezvousIP, templateData.RendezvousIP)

}

func TestIgnition_getRendezvousHostEnv(t *testing.T) {
	nodeZeroIP := "2001:db8::dead:beef"
	token := "someToken"
	rendezvousHostEnv := getRendezvousHostEnv("http", nodeZeroIP, token, workflow.AgentWorkflowTypeInstall)
	assert.Equal(t,
		"NODE_ZERO_IP="+nodeZeroIP+"\nSERVICE_BASE_URL=http://["+nodeZeroIP+"]:8090/\nIMAGE_SERVICE_BASE_URL=http://["+nodeZeroIP+"]:8888/\nAGENT_AUTH_TOKEN="+token+"\nPULL_SECRET_TOKEN="+token+"\nWORKFLOW_TYPE=install\n",
		rendezvousHostEnv)
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
			expectedError:    ".*invalid IP address syntax",
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
			expectedError: "missing rendezvousIP in agent-config, at least one host networkConfig, or at least one NMStateConfig manifest",
		},
		{
			Name: "non-canonical-ipv6-address",
			agentConfig: &agent.Config{
				RendezvousIP: "fd2e:6f44:5dd8:c956:0000:0000:0000:0050",
				Hosts: []agent.Host{
					{
						Hostname: "control-0.example.org",
						Role:     "master",
					},
				},
			},
			expectedRendezvousIP: "fd2e:6f44:5dd8:c956::50",
		},
		{
			Name: "invalid-ipv6-address",
			agentConfig: &agent.Config{
				RendezvousIP: "fd2e:6f44:5dd8:c956::0000::0050",
				Hosts: []agent.Host{
					{
						Hostname: "control-0.example.org",
						Role:     "master",
					},
				},
			},
			expectedError: "invalid rendezvous IP: fd2e:6f44:5dd8:c956::0000::0050",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var hosts []agent.Host
			if tc.agentConfig != nil {
				hosts = tc.agentConfig.Hosts
			}
			rendezvousIP, err := RetrieveRendezvousIP(tc.agentConfig, hosts, tc.nmStateConfigs)
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
		agentHosts                      *agentconfig.AgentHosts
		expectedNumberOfHostConfigFiles int
	}{
		{
			Name: "one-host-role-defined",
			agentHosts: &agentconfig.AgentHosts{
				Hosts: []agent.Host{
					{
						Role: "master",
					},
				},
			},
			expectedNumberOfHostConfigFiles: 1,
		},
		{
			Name: "multiple-host-roles-defined",
			agentHosts: &agentconfig.AgentHosts{
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
			err := addHostConfig(config, tc.agentHosts)
			assert.NoError(t, err)
			assert.Equal(t, len(config.Storage.Files), tc.expectedNumberOfHostConfigFiles)
			for _, file := range config.Storage.Files {
				assert.Equal(t, true, strings.HasPrefix(file.Path, "/etc/assisted/hostconfig"))
				assert.Equal(t, true, strings.HasSuffix(file.Path, "role"))
			}
		})
	}
}

func generatedFiles(otherFiles ...string) []string {
	files := []string{
		"/etc/assisted/rendezvous-host.env",
		"/etc/assisted/manifests/agent-config.yaml",
		// TODO: ZTP manifest files should also be present. Bug?
		// "/etc/assisted/manifests/cluster-deployment.yaml",
		// "/etc/assisted/manifests/agent-cluster-install.yaml",
		// "/etc/assisted/manifests/pull-secret.yaml",
		// "/etc/assisted/manifests/cluster-image-set.yaml",
		// "/etc/assisted/manifests/infraenv.yaml",
		"/etc/assisted/network/host0/eth0.nmconnection",
		"/etc/assisted/network/host0/mac_interface.ini",
		"/usr/local/bin/oci-eval-user-data.sh",
		"/usr/local/bin/pre-network-manager-config.sh",
		"/opt/agent/tls/kubeadmin-password.hash"}
	files = append(files, otherFiles...)
	return append(files, commonFiles()...)
}

func commonFiles() []string {
	return []string{
		"/etc/issue",
		"/etc/multipath.conf",
		"/etc/containers/containers.conf",
		"/etc/NetworkManager/conf.d/clientid.conf",
		"/root/.docker/config.json",
		"/root/assisted.te",
		"/usr/local/bin/agent-config-image-wait.sh",
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
		"/usr/local/share/start-cluster/start-cluster.env",
		"/usr/local/share/assisted-service/images.env",
		"/usr/local/bin/bootstrap-service-record.sh",
		"/usr/local/bin/release-image.sh",
		"/usr/local/bin/release-image-download.sh",
		"/etc/assisted/agent-installer.env",
		"/etc/motd.d/10-agent-installer",
		"/etc/systemd/system.conf.d/10-default-env.conf",
		"/usr/local/bin/install-status.sh",
		"/usr/local/bin/issue_status.sh",
		"/usr/local/bin/load-config-iso.sh",
		"/etc/udev/rules.d/80-agent-config-image.rules",
		"/usr/local/bin/add-node.sh",
		"/usr/local/bin/agent-auth-token-status.sh",
		"/usr/local/bin/common.sh",
	}
}

func TestIgnition_Generate(t *testing.T) {
	skipTestIfnmstatectlIsMissing(t)

	// This patch currently allows testing the Ignition asset using the embedded resources.
	// TODO: Replace it by mocking the filesystem in bootstrap.AddStorageFiles()
	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(path.Join(workingDirectory, "../../../../data"))
	assert.NoError(t, err)

	cases := []struct {
		name                string
		overrideAssets      map[reflect.Type]func(t *testing.T, dep asset.Asset) asset.Asset
		expectedError       string
		expectedFiles       []string
		expectedFileContent map[string]string
		serviceEnabledMap   map[string]bool
	}{
		{
			name: "default",
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": true,
				"agent-check-config-image.service":   false},
			expectedFiles: generatedFiles(),
		},
		{
			name: "with extra manifests",
			overrideAssets: map[reflect.Type]func(t *testing.T, dep asset.Asset) asset.Asset{
				reflect.TypeOf(manifests.ExtraManifests{}): func(t *testing.T, dep asset.Asset) asset.Asset {
					t.Helper()
					return &manifests.ExtraManifests{
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
					}
				},
			},
			expectedFiles: generatedFiles("/etc/assisted/extra-manifests/test-configmap-0.yaml", "/etc/assisted/extra-manifests/test-configmap-1.yaml"),
			expectedFileContent: map[string]string{
				"/etc/assisted/extra-manifests/test-configmap-0.yaml": `apiVersion: v1
kind: ConfigMap
metadata:
  name: agent-test-1`,
				"/etc/assisted/extra-manifests/test-configmap-1.yaml": `apiVersion: v1
kind: ConfigMap
metadata:
  name: agent-test-2
`,
			},
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": true,
				"agent-check-config-image.service":   false},
		},
		{
			name: "no nmstateconfigs defined, pre-network-manager-config.service should not be enabled",
			overrideAssets: map[reflect.Type]func(t *testing.T, dep asset.Asset) asset.Asset{
				reflect.TypeOf(manifests.AgentManifests{}): func(t *testing.T, dep asset.Asset) asset.Asset {
					t.Helper()
					am, ok := dep.(*manifests.AgentManifests)
					assert.True(t, ok)
					// remove nmstate configuration
					am.NMStateConfigs = []*aiv1beta1.NMStateConfig{}
					am.StaticNetworkConfigs = []*models.HostStaticNetworkConfig{}
					return am
				},
			},
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": false},
		},
		{
			name: "with additional ntp sources",
			overrideAssets: map[reflect.Type]func(t *testing.T, dep asset.Asset) asset.Asset{
				reflect.TypeOf(manifests.AgentManifests{}): func(t *testing.T, dep asset.Asset) asset.Asset {
					t.Helper()
					am, ok := dep.(*manifests.AgentManifests)
					assert.True(t, ok)
					am.InfraEnv.Spec.AdditionalNTPSources = []string{"0.clock.ntp.org", "1.clock.ntp.org"}
					return am
				},
			},
			expectedFiles: generatedFiles("/etc/chrony.conf"),
			expectedFileContent: map[string]string{
				"/etc/chrony.conf": `server 0.clock.ntp.org iburst
server 1.clock.ntp.org iburst
makestep 1.0 3
rtcsync
logdir /var/log/chrony`,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deps := buildIgnitionAssetDefaultDependencies(t)

			if tc.overrideAssets != nil {
				for i, dep := range deps {
					overrideFunc, found := tc.overrideAssets[reflect.TypeOf(dep).Elem()]
					if !found {
						continue
					}
					deps[i] = overrideFunc(t, dep)
				}
			}

			parents := asset.Parents{}
			parents.Add(deps...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assert.Len(t, ignitionAsset.Config.Storage.Directories, 1)
				assert.Equal(t, "/etc/assisted/extra-manifests", ignitionAsset.Config.Storage.Directories[0].Node.Path)

				assertExpectedFiles(t, ignitionAsset.Config, tc.expectedFiles, tc.expectedFileContent)

				assertServiceEnabled(t, ignitionAsset.Config, tc.serviceEnabledMap)
			}
		})
	}
}

func skipTestIfnmstatectlIsMissing(t *testing.T) {
	t.Helper()
	// Generate calls addStaticNetworkConfig which calls nmstatectl
	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}
}

func overrideDeps(deps []asset.Asset, overrides []asset.Asset) {
	for _, od := range overrides {
		for i, d := range deps {
			if d.Name() == od.Name() {
				deps[i] = od
				break
			}
		}
	}
}

func assertServiceEnabled(t *testing.T, config *igntypes.Config, serviceEnabledMap map[string]bool) {
	t.Helper()
	for serviceName, enabled := range serviceEnabledMap {
		for _, unit := range config.Systemd.Units {
			if unit.Name == serviceName {
				if unit.Enabled == nil {
					assert.Equal(t, enabled, false)
				} else {
					assert.Equal(t, enabled, *unit.Enabled)
				}
			}
		}
	}
}

func assertExpectedFiles(t *testing.T, config *igntypes.Config, expectedFiles []string, expectedFileContent map[string]string) {
	t.Helper()
	if len(expectedFiles) > 0 {
		assert.Equal(t, len(expectedFiles), len(config.Storage.Files))

		for _, f := range expectedFiles {
			found := false
			for _, i := range config.Storage.Files {
				if i.Node.Path == f {
					if expectedData, ok := expectedFileContent[i.Node.Path]; ok {
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
}

// This test util create the minimum valid set of dependencies for the
// Ignition asset
func buildIgnitionAssetDefaultDependencies(t *testing.T) []asset.Asset {
	t.Helper()
	secretDataBytes, err := base64.StdEncoding.DecodeString("c3VwZXItc2VjcmV0Cg==")
	assert.NoError(t, err)

	return []asset.Asset{
		&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
		&joiner.ClusterInfo{},
		&joiner.AddNodesConfig{},
		&manifests.AgentManifests{
			InfraEnv: &aiv1beta1.InfraEnv{
				Spec: aiv1beta1.InfraEnvSpec{
					SSHAuthorizedKey: "my-ssh-key",
				},
			},
			ClusterDeployment: &hivev1.ClusterDeployment{
				Spec: hivev1.ClusterDeploymentSpec{
					ClusterName: "ostest",
					BaseDomain:  "ostest",
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
		&agentconfig.AgentHosts{},
		&manifests.ExtraManifests{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
		&password.KubeadminPassword{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
		&tls.AdminKubeConfigClientCertKey{},
		&gencrypto.AuthConfig{},
		&common.InfraEnvID{},
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
