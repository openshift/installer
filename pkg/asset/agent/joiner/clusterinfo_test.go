package joiner

import (
	"testing"

	"github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	fakeclientconfig "github.com/openshift/client-go/config/clientset/versioned/fake"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/stretchr/testify/assert"
)

func TestClusterInfo_Generate(t *testing.T) {
	cases := []struct {
		name                string
		workflow            workflow.AgentWorkflowType
		objects             []runtime.Object
		openshiftOjects     []runtime.Object
		configHost          string
		expectedClusterInfo ClusterInfo
	}{
		{
			name:                "skip if not add-nodes workflow",
			workflow:            workflow.AgentWorkflowTypeInstall,
			expectedClusterInfo: ClusterInfo{},
		},
		{
			name:       "default",
			workflow:   workflow.AgentWorkflowTypeAddNodes,
			configHost: "https://api.ostest.test.metalkube.org:6443",
			openshiftOjects: []runtime.Object{
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
			},
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
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
				PlatformType: v1beta1.BareMetalPlatformType,
				SSHKey:       "my-ssh-key",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agentWorkflow := &workflow.AgentWorkflow{Workflow: tc.workflow}
			addNodesConfig := &AddNodesConfig{}
			parents := asset.Parents{}
			parents.Add(agentWorkflow)
			parents.Add(addNodesConfig)

			fakeClient := fake.NewSimpleClientset(tc.objects...)
			fakeOCClient := fakeclientconfig.NewSimpleClientset(tc.openshiftOjects...)
			fakeConfig := &rest.Config{
				Host: tc.configHost,
			}

			clusterInfo := &ClusterInfo{
				Config:          fakeConfig,
				Client:          fakeClient,
				OpenshiftClient: fakeOCClient,
			}
			err := clusterInfo.Generate(parents)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedClusterInfo.ClusterID, clusterInfo.ClusterID)
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
		})
	}
}

func makeInstallConfig(t *testing.T) string {
	ic := &types.InstallConfig{
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
