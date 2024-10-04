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

	configv1 "github.com/openshift/api/config/v1"
	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	fakeclientconfig "github.com/openshift/client-go/config/clientset/versioned/fake"
	fakeclientmachineconfig "github.com/openshift/client-go/machineconfiguration/clientset/versioned/fake"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
)

func TestClusterInfo_Generate(t *testing.T) {
	cases := []struct {
		name                string
		workflow            workflow.AgentWorkflowType
		nodesConfig         AddNodesConfig
		objs                func(t *testing.T) ([]runtime.Object, []runtime.Object, []runtime.Object)
		expectedClusterInfo ClusterInfo
		expectedError       string
	}{
		{
			name:                "skip if not add-nodes workflow",
			workflow:            workflow.AgentWorkflowTypeInstall,
			expectedClusterInfo: ClusterInfo{},
		},
		{
			name:     "default",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			objs:     defaultObjects(),
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
				ClusterName:  "ostest",
				ReleaseImage: "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
				Version:      "4.15.0",
				APIDNSName:   "api.ostest.test.metalkube.org",
				Namespace:    "cluster0",
				PullSecret:   "c3VwZXJzZWNyZXQK", // notsecret
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
				FIPS: true,
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
			objs: func(t *testing.T) ([]runtime.Object, []runtime.Object, []runtime.Object) {
				t.Helper()
				objs, ocObjs, ocMachineConfigObjs := defaultObjects()(t)
				for i, o := range objs {
					if node, ok := o.(*corev1.Node); ok {
						node.Status.NodeInfo.Architecture = "amd64"
						objs[i] = node
						break
					}
				}
				return objs, ocObjs, ocMachineConfigObjs
			},
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
				ClusterName:  "ostest",
				ReleaseImage: "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
				Version:      "4.15.0",
				APIDNSName:   "api.ostest.test.metalkube.org",
				Namespace:    "cluster0",
				PullSecret:   "c3VwZXJzZWNyZXQK", // notsecret
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
				FIPS:            true,
				IgnitionEndpointWorker: &models.IgnitionEndpoint{
					URL:           ptr.To("https://192.168.111.5:22623/config/worker"),
					CaCertificate: ptr.To("LS0tL_FakeCertificate_LS0tCg=="),
				},
			},
		},
		{
			name:     "not supported platform",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			objs: func(t *testing.T) ([]runtime.Object, []runtime.Object, []runtime.Object) {
				t.Helper()
				objs, ocObjs, ocMachineConfigObjs := defaultObjects()(t)
				for i, o := range ocObjs {
					if infra, ok := o.(*configv1.Infrastructure); ok {
						infra.Spec.PlatformSpec.Type = configv1.AWSPlatformType
						ocObjs[i] = infra
						break
					}
				}
				return objs, ocObjs, ocMachineConfigObjs
			},
			expectedError: "Platform: Unsupported value: \"aws\": supported values: \"baremetal\", \"vsphere\", \"none\", \"external\"",
		},
		{
			name:     "sshKey from nodes-config.yaml",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			objs:     defaultObjects(),
			nodesConfig: AddNodesConfig{
				Config: Config{
					SSHKey: "ssh-key-from-config",
				},
			},
			expectedClusterInfo: ClusterInfo{
				ClusterID:    "1b5ba46b-7e56-47b1-a326-a9eebddfb38c",
				ClusterName:  "ostest",
				ReleaseImage: "registry.ci.openshift.org/ocp/release@sha256:65d9b652d0d23084bc45cb66001c22e796d43f5e9e005c2bc2702f94397d596e",
				Version:      "4.15.0",
				APIDNSName:   "api.ostest.test.metalkube.org",
				Namespace:    "cluster0",
				PullSecret:   "c3VwZXJzZWNyZXQK", // notsecret
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
				SSHKey:          "ssh-key-from-config",
				OSImage:         buildStreamData(),
				OSImageLocation: "http://my-coreosimage-url/416.94.202402130130-0",
				IgnitionEndpointWorker: &models.IgnitionEndpoint{
					URL:           ptr.To("https://192.168.111.5:22623/config/worker"),
					CaCertificate: ptr.To("LS0tL_FakeCertificate_LS0tCg=="),
				},
				FIPS: true,
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

			var objects, openshiftObjects, openshiftMachineConfigObjects []runtime.Object
			if tc.objs != nil {
				objects, openshiftObjects, openshiftMachineConfigObjects = tc.objs(t)
			}
			fakeClient := fake.NewSimpleClientset(objects...)
			fakeOCClient := fakeclientconfig.NewSimpleClientset(openshiftObjects...)
			fakeOCMachineConfigClient := fakeclientmachineconfig.NewSimpleClientset(openshiftMachineConfigObjects...)

			clusterInfo := &ClusterInfo{
				Client:                       fakeClient,
				OpenshiftClient:              fakeOCClient,
				OpenshiftMachineConfigClient: fakeOCMachineConfigClient,
			}
			err := clusterInfo.Generate(context.Background(), parents)
			if tc.expectedError == "" {
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
				assert.Equal(t, tc.expectedClusterInfo.FIPS, clusterInfo.FIPS)
			} else {
				assert.Regexp(t, tc.expectedError, err.Error())
			}
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

func defaultObjects() func(t *testing.T) ([]runtime.Object, []runtime.Object, []runtime.Object) {
	return func(t *testing.T) ([]runtime.Object, []runtime.Object, []runtime.Object) {
		t.Helper()
		objects := []runtime.Object{
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
		}

		openshiftObjects := []runtime.Object{
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
			&configv1.Infrastructure{
				ObjectMeta: v1.ObjectMeta{
					Name: "cluster",
				},
				Spec: configv1.InfrastructureSpec{
					PlatformSpec: configv1.PlatformSpec{
						Type: configv1.BareMetalPlatformType,
					},
				},
				Status: configv1.InfrastructureStatus{
					APIServerURL: "https://api.ostest.test.metalkube.org:6443",
				},
			},
			&configv1.ImageDigestMirrorSetList{
				Items: []configv1.ImageDigestMirrorSet{
					{
						Spec: configv1.ImageDigestMirrorSetSpec{
							ImageDigestMirrors: []configv1.ImageDigestMirrors{
								{
									Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
									Mirrors: []configv1.ImageMirror{
										"registry.example.com:5000/ocp4/openshift4",
									},
								},
							},
						},
					},
				},
			},
		}

		openshiftMachineConfigObjects := []runtime.Object{
			&machineconfigv1.MachineConfig{
				ObjectMeta: v1.ObjectMeta{
					Name: "99-worker-ssh",
				},
				Spec: machineconfigv1.MachineConfigSpec{
					Config: runtime.RawExtension{
						Raw: []byte(`
ignition:
  version: 3.2.0
passwd:
  users:
  - name: core
    sshAuthorizedKeys:
    - my-ssh-key`),
					},
				},
			},
			&machineconfigv1.MachineConfig{
				ObjectMeta: v1.ObjectMeta{
					Name: "99-worker-fips",
				},
				Spec: machineconfigv1.MachineConfigSpec{
					FIPS: true,
				},
			},
		}

		return objects, openshiftObjects, openshiftMachineConfigObjects
	}
}
