package manifests

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	hivev1agent "github.com/openshift/hive/apis/hive/v1/agent"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestClusterDeployment_Generate(t *testing.T) {

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *hivev1.ClusterDeployment
	}{
		{
			name: "missing config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&agent.OptionalInstallConfig{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid configurations",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				getValidOptionalInstallConfig(),
			},
			expectedConfig: &hivev1.ClusterDeployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterDeployment",
					APIVersion: "hive.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      getClusterDeploymentName(getValidOptionalInstallConfig()),
					Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
				},
				Spec: hivev1.ClusterDeploymentSpec{
					ClusterName: getClusterDeploymentName(getValidOptionalInstallConfig()),
					BaseDomain:  "testing.com",
					PullSecretRef: &corev1.LocalObjectReference{
						Name: getPullSecretName(getValidOptionalInstallConfig().ClusterName()),
					},
					ClusterInstallRef: &hivev1.ClusterInstallLocalReference{
						Group:   "extensions.hive.openshift.io",
						Version: "v1beta1",
						Kind:    "AgentClusterInstall",
						Name:    getAgentClusterInstallName(getValidOptionalInstallConfig()),
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ClusterDeployment{}
			err := asset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/cluster-deployment.yaml", configFile.Filename)

				var actualConfig hivev1.ClusterDeployment
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}
		})
	}

}

func TestClusterDeployment_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  bool
		expectedConfig *hivev1.ClusterDeployment
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: compact-cluster
  namespace: cluster0
spec:
  baseDomain: agent.example.com
  clusterInstallRef:
    group: extensions.hive.openshift.io
    kind: AgentClusterInstall
    name: test-agent-cluster-install
    version: v1beta1
  clusterName: compact-cluster
  controlPlaneConfig:
    servingCertificates: {}
  platform:
    agentBareMetal:
      agentSelector:
        matchLabels:
          bla: aaa
  pullSecretRef:
    name: pull-secret`,
			expectedFound: true,
			expectedConfig: &hivev1.ClusterDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "compact-cluster",
					Namespace: "cluster0",
				},
				Spec: hivev1.ClusterDeploymentSpec{
					BaseDomain: "agent.example.com",
					ClusterInstallRef: &hivev1.ClusterInstallLocalReference{
						Group:   "extensions.hive.openshift.io",
						Kind:    "AgentClusterInstall",
						Name:    "test-agent-cluster-install",
						Version: "v1beta1",
					},
					ClusterName:        "compact-cluster",
					ControlPlaneConfig: hivev1.ControlPlaneConfigSpec{},
					Platform: hivev1.Platform{
						AgentBareMetal: &hivev1agent.BareMetalPlatform{
							AgentSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{
									"bla": "aaa",
								},
							},
						},
					},
					PullSecretRef: &corev1.LocalObjectReference{
						Name: "pull-secret",
					},
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: true,
		},
		{
			name:           "empty",
			data:           "",
			expectedFound:  true,
			expectedConfig: &hivev1.ClusterDeployment{},
			expectedError:  false,
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: true,
		},
		{
			name: "unknown-field",
			data: `
metadata:
  name: cluster-deployment-bad
  namespace: cluster0
spec:
  wrongField: wrongValue`,
			expectedError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(clusterDeploymentFilename).
				Return(
					&asset.File{
						Filename: clusterDeploymentFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ClusterDeployment{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ClusterDeployment")
			}
		})
	}

}
