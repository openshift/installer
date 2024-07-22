package joiner

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	fakeclientconfig "github.com/openshift/client-go/config/clientset/versioned/fake"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestClusterInfo_Generate(t *testing.T) {
	cases := []struct {
		name                string
		workflow            workflow.AgentWorkflowType
		nodesConfig         AddNodesConfig
		objects             []runtime.Object
		openshiftObjects    []runtime.Object
		expectedClusterInfo ClusterInfo
	}{
		{
			name:                "skip if not add-nodes workflow",
			workflow:            workflow.AgentWorkflowTypeInstall,
			expectedClusterInfo: ClusterInfo{},
		},
		{
			name:     "default",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			openshiftObjects: []runtime.Object{
				&configv1.ClusterVersion{
					ObjectMeta: v1.ObjectMeta{
						Name: "version",
					},
					Spec: configv1.ClusterVersionSpec{
						ClusterID: "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
					},
					Status: configv1.ClusterVersionStatus{
						History: []configv1.UpdateHistory{
							{
								Image:   "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
								Version: "4.15.0",
							},
						},
					},
				},
				&configv1.Proxy{
					ObjectMeta: v1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProxySpec{
						HTTPProxy:  "http://proxy",
						HTTPSProxy: "https://proxy",
						NoProxy:    "localhost",
					},
				},
			},
			objects: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "pull-secret",
						Namespace: "openshift-config",
					},
					Data: map[string][]byte{
						".dockerconfigjson": []byte("c3VwZXJzZWNyZXQK"),
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "user-ca-bundle",
						Namespace: "openshift-config",
					},
					Data: map[string]string{
						"ca-bundle.crt": "--- bundle ---",
					},
				},
				&corev1.Node{
					ObjectMeta: v1.ObjectMeta{
						Labels: map[string]string{
							"node-role.kubernetes.io/master": "",
						},
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							Architecture: "amd64",
						},
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "cluster-config-v1",
						Namespace: "kube-system",
					},
					Data: map[string]string{
						"install-config": makeInstallConfig(t),
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "coreos-bootimages",
						Namespace: "openshift-machine-config-operator",
					},
					Data: map[string]string{
						"stream": makeCoreOsBootImages(t, buildStreamData()),
					},
				},
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "worker-user-data-managed",
						Namespace: "openshift-machine-api",
					},
					Data: map[string][]byte{
						"userData": []byte(`{"ignition":{"config":{"merge":[{"source":"https://192.168.111.5:22623/config/worker","verification":{}}],
"replace":{"verification":{}}},"proxy":{},"security":{"tls":{"certificateAuthorities":[{"source":"data:text/plain;charset=utf-8;base64,LS0tL_FakeCertificate_LS0tCg==",
"verification":{}}]}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{},"storage":{},"systemd":{}}`),
					},
				},
			},
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
				ClusterName:  "ostest",
				ReleaseImage: "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
				Version:      "4.15.0",
				APIDNSName:   "api.ostest.test.metalkube.org",
				Namespace:    "cluster0",
				PullSecret:   "c3VwZXJzZWNyZXQK",
				UserCaBundle: "--- bundle ---",
				Architecture: "amd64",
				Proxy: &types.Proxy{
					HTTPProxy:  "http://proxy",
					HTTPSProxy: "https://proxy",
					NoProxy:    "localhost",
				},
				ImageDigestSources: []types.ImageDigestSource{
					{
						Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
						Mirrors: []string{
							"registry.example.com:5000/ocp4/openshift4",
						},
					},
				},
				PlatformType:    v1beta1.BareMetalPlatformType,
				SSHKey:          "my-ssh-key",
				OSImage:         buildStreamData(),
				OSImageLocation: "http://my-coreosimage-url/416.94.202402130130-0",
				IgnitionEndpointWorker: &models.IgnitionEndpoint{
					URL:           ptr.To("https://192.168.111.5:22623/config/worker"),
					CaCertificate: ptr.To("LS0tL_FakeCertificate_LS0tCg=="),
				},
			},
		},
		{
			name:     "architecture specified in nodesConfig as arm64 and target cluster is amd64",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			nodesConfig: AddNodesConfig{
				Config: Config{
					CPUArchitecture: "arm64",
				},
			},
			openshiftObjects: []runtime.Object{
				&configv1.ClusterVersion{
					ObjectMeta: v1.ObjectMeta{
						Name: "version",
					},
					Spec: configv1.ClusterVersionSpec{
						ClusterID: "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
					},
					Status: configv1.ClusterVersionStatus{
						History: []configv1.UpdateHistory{
							{
								Image:   "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
								Version: "4.15.0",
							},
						},
					},
				},
				&configv1.Proxy{
					ObjectMeta: v1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProxySpec{
						HTTPProxy:  "http://proxy",
						HTTPSProxy: "https://proxy",
						NoProxy:    "localhost",
					},
				},
			},
			objects: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "pull-secret",
						Namespace: "openshift-config",
					},
					Data: map[string][]byte{
						".dockerconfigjson": []byte("c3VwZXJzZWNyZXQK"),
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "user-ca-bundle",
						Namespace: "openshift-config",
					},
					Data: map[string]string{
						"ca-bundle.crt": "--- bundle ---",
					},
				},
				&corev1.Node{
					ObjectMeta: v1.ObjectMeta{
						Labels: map[string]string{
							"node-role.kubernetes.io/master": "",
						},
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							Architecture: "amd64",
						},
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "cluster-config-v1",
						Namespace: "kube-system",
					},
					Data: map[string]string{
						"install-config": makeInstallConfig(t),
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "coreos-bootimages",
						Namespace: "openshift-machine-config-operator",
					},
					Data: map[string]string{
						"stream": makeCoreOsBootImages(t, buildStreamData()),
					},
				},
			},
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
				ClusterName:  "ostest",
				ReleaseImage: "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
				Version:      "4.15.0",
				APIDNSName:   "api.ostest.test.metalkube.org",
				Namespace:    "cluster0",
				PullSecret:   "c3VwZXJzZWNyZXQK",
				UserCaBundle: "--- bundle ---",
				Architecture: "arm64",
				Proxy: &types.Proxy{
					HTTPProxy:  "http://proxy",
					HTTPSProxy: "https://proxy",
					NoProxy:    "localhost",
				},
				ImageDigestSources: []types.ImageDigestSource{
					{
						Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
						Mirrors: []string{
							"registry.example.com:5000/ocp4/openshift4",
						},
					},
				},
				PlatformType:    v1beta1.BareMetalPlatformType,
				SSHKey:          "my-ssh-key",
				OSImage:         buildStreamData(),
				OSImageLocation: "http://my-coreosimage-url/416.94.202402130130-1",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agentWorkflow := &workflow.AgentWorkflow{Workflow: tc.workflow}
			addNodesConfig := &tc.nodesConfig
			parents := asset.Parents{}
			parents.Add(agentWorkflow)
			parents.Add(addNodesConfig)

			fakeClient := fake.NewSimpleClientset(tc.objects...)
			fakeOCClient := fakeclientconfig.NewSimpleClientset(tc.openshiftObjects...)

			clusterInfo := &ClusterInfo{
				Client:          fakeClient,
				OpenshiftClient: fakeOCClient,
			}
			err := clusterInfo.Generate(context.Background(), parents)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedClusterInfo.ClusterID, clusterInfo.ClusterID)
			assert.Equal(t, tc.expectedClusterInfo.ClusterName, clusterInfo.ClusterName)
			assert.Equal(t, tc.expectedClusterInfo.Version, clusterInfo.Version)
			assert.Equal(t, tc.expectedClusterInfo.ReleaseImage, clusterInfo.ReleaseImage)
			assert.Equal(t, tc.expectedClusterInfo.APIDNSName, clusterInfo.APIDNSName)
			assert.Equal(t, tc.expectedClusterInfo.PullSecret, clusterInfo.PullSecret)
			assert.Equal(t, tc.expectedClusterInfo.Namespace, clusterInfo.Namespace)
			assert.Equal(t, tc.expectedClusterInfo.UserCaBundle, clusterInfo.UserCaBundle)
			assert.Equal(t, tc.expectedClusterInfo.Proxy, clusterInfo.Proxy)
			assert.Equal(t, tc.expectedClusterInfo.Architecture, clusterInfo.Architecture)
			assert.Equal(t, tc.expectedClusterInfo.ImageDigestSources, clusterInfo.ImageDigestSources)
			assert.Equal(t, tc.expectedClusterInfo.DeprecatedImageContentSources, clusterInfo.DeprecatedImageContentSources)
			assert.Equal(t, tc.expectedClusterInfo.PlatformType, clusterInfo.PlatformType)
			assert.Equal(t, tc.expectedClusterInfo.SSHKey, clusterInfo.SSHKey)
			assert.Equal(t, tc.expectedClusterInfo.OSImageLocation, clusterInfo.OSImageLocation)
			assert.Equal(t, tc.expectedClusterInfo.OSImage, clusterInfo.OSImage)
			assert.Equal(t, tc.expectedClusterInfo.IgnitionEndpointWorker, clusterInfo.IgnitionEndpointWorker)
		})
	}
}

func buildStreamData() *stream.Stream {
	return &stream.Stream{
		Architectures: map[string]stream.Arch{
			"x86_64": {
				Artifacts: map[string]stream.PlatformArtifacts{
					"metal": {
						Release: "416.94.202402130130-0",
						Formats: map[string]stream.ImageFormat{
							"iso": {
								Disk: &stream.Artifact{
									Location: "http://my-coreosimage-url/416.94.202402130130-0",
								},
							},
						},
					},
				},
			},
			"aarch64": {
				Artifacts: map[string]stream.PlatformArtifacts{
					"metal": {
						Release: "416.94.202402130130-0",
						Formats: map[string]stream.ImageFormat{
							"iso": {
								Disk: &stream.Artifact{
									Location: "http://my-coreosimage-url/416.94.202402130130-1",
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeCoreOsBootImages(t *testing.T, st *stream.Stream) string {
	t.Helper()
	data, err := json.Marshal(st)
	if err != nil {
		t.Error(err)
	}

	return string(data)
}

func makeInstallConfig(t *testing.T) string {
	t.Helper()
	ic := &types.InstallConfig{
		ObjectMeta: v1.ObjectMeta{
			Name: "ostest",
		},
		BaseDomain: "test.metalkube.org",
		ImageDigestSources: []types.ImageDigestSource{
			{
				Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
				Mirrors: []string{
					"registry.example.com:5000/ocp4/openshift4",
				},
			},
		},
		Platform: types.Platform{
			BareMetal: &baremetal.Platform{},
		},
		SSHKey: "my-ssh-key",
	}
	data, err := yaml.Marshal(ic)
	if err != nil {
		t.Error(err)
	}

	return string(data)
}
