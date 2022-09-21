package manifests

import (
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"

	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestAgentClusterInstall_Generate(t *testing.T) {

	goodACI := &hiveext.AgentClusterInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getAgentClusterInstallName(getValidOptionalInstallConfig()),
			Namespace: getObjectMetaNamespace(getValidOptionalInstallConfig()),
		},
		Spec: hiveext.AgentClusterInstallSpec{
			ImageSetRef: &hivev1.ClusterImageSetReference{
				Name: getClusterImageSetReferenceName(),
			},
			ClusterDeploymentRef: corev1.LocalObjectReference{
				Name: getClusterDeploymentName(getValidOptionalInstallConfig()),
			},
			Networking: hiveext.Networking{
				ClusterNetwork: []hiveext.ClusterNetworkEntry{
					{
						CIDR:       "192.168.111.0/24",
						HostPrefix: 23,
					},
				},
				ServiceNetwork: []string{"172.30.0.0/16"},
				NetworkType:    "OVNKubernetes",
			},
			SSHPublicKey: strings.Trim(TestSSHKey, "|\n\t"),
			ProvisionRequirements: hiveext.ProvisionRequirements{
				ControlPlaneAgents: 3,
				WorkerAgents:       5,
			},
			APIVIP:     "192.168.122.10",
			IngressVIP: "192.168.122.11",
		},
	}

	installConfigWithoutNetworkType := getValidOptionalInstallConfig()
	installConfigWithoutNetworkType.Config.NetworkType = ""

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

// func TestAgentClusterInstall_Generate(t *testing.T) {

// 	installConfig := &agent.OptionalInstallConfig{
// 		Config: &types.InstallConfig{
// 			ObjectMeta: v1.ObjectMeta{
// 				Name:      "cluster0-name",
// 				Namespace: "cluster0-namespace",
// 			},
// 			SSHKey: "ssh-key",
// 			ControlPlane: &types.MachinePool{
// 				Name:     "master",
// 				Replicas: pointer.Int64Ptr(3),
// 				Platform: types.MachinePoolPlatform{},
// 			},
// 			Compute: []types.MachinePool{
// 				{
// 					Name:     "worker-machine-pool-1",
// 					Replicas: pointer.Int64Ptr(2),
// 				},
// 				{
// 					Name:     "worker-machine-pool-2",
// 					Replicas: pointer.Int64Ptr(3),
// 				},
// 			},
// 		},
// 	}

// 	parents := asset.Parents{}
// 	parents.Add(installConfig)

// 	asset := &AgentClusterInstall{}
// 	err := asset.Generate(parents)
// 	assert.NoError(t, err)

// 	assert.NotEmpty(t, asset.Files())
// 	aciFile := asset.Files()[0]
// 	assert.Equal(t, "cluster-manifests/agent-cluster-install.yaml", aciFile.Filename)

// 	aci := &hiveext.AgentClusterInstall{}
// 	err = yaml.Unmarshal(aciFile.Data, &aci)
// 	assert.NoError(t, err)

// 	assert.Equal(t, "agent-cluster-install", aci.Name)
// 	assert.Equal(t, "cluster0-namespace", aci.Namespace)
// 	assert.Equal(t, "cluster0-name", aci.Spec.ClusterDeploymentRef.Name)
// 	assert.Equal(t, 3, aci.Spec.ProvisionRequirements.ControlPlaneAgents)

// 	assert.Equal(t, 5, aci.Spec.ProvisionRequirements.WorkerAgents)
// 	assert.Equal(t, "ssh-key", aci.Spec.SSHPublicKey)
// }

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
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in AgentClusterInstall")
			}
		})
	}

}
