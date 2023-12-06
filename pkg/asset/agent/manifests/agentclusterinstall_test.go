package manifests

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	configv1 "github.com/openshift/api/config/v1"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
	externaltype "github.com/openshift/installer/pkg/types/external"
)

func TestAgentClusterInstall_Generate(t *testing.T) {

	installConfigWithoutNetworkType := getValidOptionalInstallConfig()
	installConfigWithoutNetworkType.Config.NetworkType = ""

	installConfigWithFIPS := getValidOptionalInstallConfig()
	installConfigWithFIPS.Config.FIPS = true

	goodACI := getGoodACI()
	goodFIPSACI := getGoodACI()
	goodFIPSACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"fips":true}`,
	})

	installConfigWithProxy := getValidOptionalInstallConfig()
	installConfigWithProxy.Config.Proxy = (*types.Proxy)(getProxy(getProxyValidOptionalInstallConfig()))

	goodProxyACI := getGoodACI()
	goodProxyACI.Spec.Proxy = (*hiveext.Proxy)(getProxy(getProxyValidOptionalInstallConfig()))

	goodACIDualStackVIPs := getGoodACIDualStack()
	goodACIDualStackVIPs.Spec.APIVIPs = []string{"192.168.122.10", "2001:db8:1111:2222:ffff:ffff:ffff:cafe"}
	goodACIDualStackVIPs.Spec.IngressVIPs = []string{"192.168.122.11", "2001:db8:1111:2222:ffff:ffff:ffff:dead"}

	installConfigWithCapabilities := getValidOptionalInstallConfig()
	installConfigWithCapabilities.Config.Capabilities = &types.Capabilities{
		BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
		AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{
			configv1.ClusterVersionCapabilityMarketplace,
		},
	}

	goodCapabilitiesACI := getGoodACI()
	goodCapabilitiesACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"capabilities":{"baselineCapabilitySet":"None","additionalEnabledCapabilities":["marketplace"]}}`,
	})

	installConfigWithNetworkOverride := getValidOptionalInstallConfig()
	installConfigWithNetworkOverride.Config.Networking.NetworkType = "CustomNetworkType"

	goodNetworkOverrideACI := getGoodACI()
	goodNetworkOverrideACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"networking":{"networkType":"CustomNetworkType","machineNetwork":[{"cidr":"10.10.11.0/24"}],"clusterNetwork":[{"cidr":"192.168.111.0/24","hostPrefix":23}],"serviceNetwork":["172.30.0.0/16"]}}`,
	})

	installConfigWithCPUPartitioning := getValidOptionalInstallConfig()
	installConfigWithCPUPartitioning.Config.CPUPartitioning = types.CPUPartitioningAllNodes

	goodCPUPartitioningACI := getGoodACI()
	goodCPUPartitioningACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"cpuPartitioningMode":"AllNodes"}`,
	})

	installConfigWExternalPlatform := getValidOptionalInstallConfig()
	installConfigWExternalPlatform.Config.Platform = types.Platform{
		External: &externaltype.Platform{
			PlatformName:           "external",
			CloudControllerManager: "",
		},
	}

	goodExternalPlatformACI := getGoodACI()
	goodExternalPlatformACI.Spec.APIVIPs = nil
	goodExternalPlatformACI.Spec.IngressVIPs = nil
	val := true
	goodExternalPlatformACI.Spec.Networking.UserManagedNetworking = &val
	goodExternalPlatformACI.Spec.PlatformType = hiveext.ExternalPlatformType
	goodExternalPlatformACI.Spec.ExternalPlatformSpec = &hiveext.ExternalPlatformSpec{
		PlatformName: "external",
	}
	goodExternalPlatformACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"platform":{"external":{"platformName":"external"}}}`,
	})

	installConfigWExternalOCIPlatform := getValidOptionalInstallConfig()
	installConfigWExternalOCIPlatform.Config.Platform = types.Platform{
		External: &externaltype.Platform{
			PlatformName:           string(models.PlatformTypeOci),
			CloudControllerManager: externaltype.CloudControllerManagerTypeExternal,
		},
	}

	goodExternalOCIPlatformACI := getGoodACI()
	val = true
	goodExternalOCIPlatformACI.Spec.APIVIPs = nil
	goodExternalOCIPlatformACI.Spec.IngressVIPs = nil
	goodExternalOCIPlatformACI.Spec.Networking.UserManagedNetworking = &val
	goodExternalOCIPlatformACI.Spec.PlatformType = hiveext.ExternalPlatformType
	goodExternalOCIPlatformACI.Spec.ExternalPlatformSpec = &hiveext.ExternalPlatformSpec{
		PlatformName: string(models.PlatformTypeOci),
	}
	goodExternalOCIPlatformACI.SetAnnotations(map[string]string{
		installConfigOverrides: `{"platform":{"external":{"platformName":"oci","cloudControllerManager":"External"}}}`,
	})

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *hiveext.AgentClusterInstall
	}{
		{
			name: "missing install config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(),
			},
			expectedConfig: goodACI,
		},
		{
			name: "valid configuration with unspecified network type should result with ACI having default network type",
			dependencies: []asset.Asset{
				installConfigWithoutNetworkType,
			},
			expectedConfig: goodACI,
		},
		{
			name: "valid configuration with FIPS annotation",
			dependencies: []asset.Asset{
				installConfigWithFIPS,
			},
			expectedConfig: goodFIPSACI,
		},
		{
			name: "valid configuration with proxy",
			dependencies: []asset.Asset{
				installConfigWithProxy,
			},
			expectedConfig: goodProxyACI,
		},
		{
			name: "valid configuration dual stack",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigDualStack(),
			},
			expectedConfig: getGoodACIDualStack(),
		},
		{
			name: "valid configuration dual stack dual vips",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigDualStackDualVIPs(),
			},
			expectedConfig: goodACIDualStackVIPs,
		},
		{
			name: "valid configuration with capabilities",
			dependencies: []asset.Asset{
				installConfigWithCapabilities,
			},
			expectedConfig: goodCapabilitiesACI,
		},
		{
			name: "valid configuration with custom network type",
			dependencies: []asset.Asset{
				installConfigWithNetworkOverride,
			},
			expectedConfig: goodNetworkOverrideACI,
		},
		{
			name: "valid configuration with CPU Partitioning",
			dependencies: []asset.Asset{
				installConfigWithCPUPartitioning,
			},
			expectedConfig: goodCPUPartitioningACI,
		},
		{
			name: "valid configuration external generic platform",
			dependencies: []asset.Asset{
				installConfigWExternalPlatform,
			},
			expectedConfig: goodExternalPlatformACI,
		},
		{
			name: "valid configuration external OCI platform",
			dependencies: []asset.Asset{
				installConfigWExternalOCIPlatform,
			},
			expectedConfig: goodExternalOCIPlatformACI,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &AgentClusterInstall{}
			err := asset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/agent-cluster-install.yaml", configFile.Filename)

				var actualConfig hiveext.AgentClusterInstall
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestAgentClusterInstall_LoadedFromDisk(t *testing.T) {

	emptyACI := &hiveext.AgentClusterInstall{}
	emptyACI.Spec.Networking.NetworkType = "OVNKubernetes"

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *hiveext.AgentClusterInstall
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  platformType: BareMetal
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
    networkType: OVNKubernetes
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:       "192.168.111.5",
					IngressVIP:   "192.168.111.4",
					PlatformType: hiveext.BareMetalPlatformType,
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.10.0",
					},
					Networking: hiveext.Networking{
						MachineNetwork: []hiveext.MachineNetworkEntry{
							{
								CIDR: "10.10.11.0/24",
							},
						},
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
						},
						NetworkType: "OVNKubernetes",
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
			expectedError: "",
		},
		{
			name: "valid-config-file-external-oci-platform",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  platformType: External
  external:
    platformName: oci
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.14.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
    networkType: OVNKubernetes
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:       "192.168.111.5",
					IngressVIP:   "192.168.111.4",
					PlatformType: hiveext.ExternalPlatformType,
					ExternalPlatformSpec: &hiveext.ExternalPlatformSpec{
						PlatformName: string(models.PlatformTypeOci),
					},
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.14.0",
					},
					Networking: hiveext.Networking{
						MachineNetwork: []hiveext.MachineNetworkEntry{
							{
								CIDR: "10.10.11.0/24",
							},
						},
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
						},
						NetworkType:           "OVNKubernetes",
						UserManagedNetworking: func(b bool) *bool { return &b }(true),
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
			expectedError: "",
		},
		{
			name: "lowercase-platform-type-backwards-compat",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  platformType: baremetal
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
    networkType: OVNKubernetes
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:       "192.168.111.5",
					IngressVIP:   "192.168.111.4",
					PlatformType: hiveext.BareMetalPlatformType,
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.10.0",
					},
					Networking: hiveext.Networking{
						MachineNetwork: []hiveext.MachineNetworkEntry{
							{
								CIDR: "10.10.11.0/24",
							},
						},
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
						},
						NetworkType: "OVNKubernetes",
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
			expectedError: "",
		},
		{
			name: "valid-config-file-network-type-openshiftsdn",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
    networkType: OpenShiftSDN
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:     "192.168.111.5",
					IngressVIP: "192.168.111.4",
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.10.0",
					},
					Networking: hiveext.Networking{
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
						},
						NetworkType: "OpenShiftSDN",
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
			expectedError: "",
		},
		{
			name: "valid-config-file-no-network-type-specified-and-defaults-to-OVNKubernetes",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:     "192.168.111.5",
					IngressVIP: "192.168.111.4",
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.10.0",
					},
					Networking: hiveext.Networking{
						MachineNetwork: []hiveext.MachineNetworkEntry{
							{
								CIDR: "10.10.11.0/24",
							},
						},
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
						},
						NetworkType: "OVNKubernetes",
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
		},
		{
			name: "valid-config-file-dual-stack",
			data: `
metadata:
  name: test-agent-cluster-install-dual-stack
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    - cidr: 2001:db8:5dd8:c956::/64
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    - cidr: 2001:db8:1111:2222::/64
      hostPrefix: 64
    serviceNetwork:
    - 172.30.0.0/16
    - fd02::/112
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-agent-cluster-install-dual-stack",
					Namespace: "cluster0",
				},
				Spec: hiveext.AgentClusterInstallSpec{
					APIVIP:     "192.168.111.5",
					IngressVIP: "192.168.111.4",
					ClusterDeploymentRef: corev1.LocalObjectReference{
						Name: "ostest",
					},
					ImageSetRef: &hivev1.ClusterImageSetReference{
						Name: "openshift-v4.10.0",
					},
					Networking: hiveext.Networking{
						MachineNetwork: []hiveext.MachineNetworkEntry{
							{
								CIDR: "10.10.11.0/24",
							},
							{
								CIDR: "2001:db8:5dd8:c956::/64",
							},
						},
						ClusterNetwork: []hiveext.ClusterNetworkEntry{
							{
								CIDR:       "10.128.0.0/14",
								HostPrefix: 23,
							},
							{
								CIDR:       "2001:db8:1111:2222::/64",
								HostPrefix: 64,
							},
						},
						ServiceNetwork: []string{
							"172.30.0.0/16",
							"fd02::/112",
						},
						NetworkType: "OVNKubernetes",
					},
					ProvisionRequirements: hiveext.ProvisionRequirements{
						ControlPlaneAgents: 3,
						WorkerAgents:       2,
					},
					SSHPublicKey: "ssh-rsa AAAAmyKey",
				},
			},
			expectedError: "",
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "failed to unmarshal cluster-manifests/agent-cluster-install.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.AgentClusterInstall",
		},
		{
			name:           "empty",
			data:           "",
			expectedFound:  true,
			expectedConfig: emptyACI,
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-manifests/agent-cluster-install.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  wrongField: wrongValue`,
			expectedError: "failed to unmarshal cluster-manifests/agent-cluster-install.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "network-ip-address-incompatible-with-network-type",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    clusterNetwork:
    - cidr: fd01::/48
      hostPrefix: 23
    serviceNetwork:
    - fd02::/112
    - 172.30.0.0/16
    networkType: "OpenShiftSDN"
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedError: "invalid NetworkType configured: [spec.networking.networkType: Required value: clusterNetwork CIDR is IPv6 and is not compatible with networkType OpenShiftSDN, spec.networking.networkType: Required value: serviceNetwork CIDR is IPv6 and is not compatible with networkType OpenShiftSDN]",
		},
		{
			name: "invalid-config-file",
			data: `
metadata:
  name: test-agent-cluster-install
  namespace: cluster0
spec:
  apiVIP: 192.168.111.5
  ingressVIP: 192.168.111.4
  platformType: aws
  clusterDeploymentRef:
    name: ostest
  imageSetRef:
    name: openshift-v4.10.0
  networking:
    machineNetwork:
    - cidr: 10.10.11.0/24
    clusterNetwork:
    - cidr: 10.128.0.0/14
      hostPrefix: 23
    serviceNetwork:
    - 172.30.0.0/16
    networkType: OVNKubernetes
  provisionRequirements:
    controlPlaneAgents: 3
    workerAgents: 2
  sshPublicKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: false,
			expectedError: "invalid PlatformType configured: spec.platformType: Unsupported value: \"aws\": supported values: \"BareMetal\", \"VSphere\", \"None\", \"External\"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(agentClusterInstallFilename).
				Return(
					&asset.File{
						Filename: agentClusterInstallFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &AgentClusterInstall{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.Equal(t, nil, err)
			}

			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in AgentClusterInstall")
			}
		})
	}

}
