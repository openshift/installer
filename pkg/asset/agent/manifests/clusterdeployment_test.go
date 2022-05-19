package manifests

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	hivev1agent "github.com/openshift/hive/apis/hive/v1/agent"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
				ObjectMeta: v1.ObjectMeta{
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
