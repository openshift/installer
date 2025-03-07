package manifests

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestInfraEnvFile_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *aiv1beta1.InfraEnv
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: infraEnv
  namespace: cluster0
spec:
  clusterRef:
    name: ocp-edge-cluster-0
    namespace: cluster0
  nmStateConfigLabelSelector: 
    matchLabels:
      cluster0-nmstate-label-name: cluster0-nmstate-label-value
  pullSecretRef:
    name: pull-secret
  sshAuthorizedKey: |
    ssh-rsa AAAAmyKey`,
			expectedFound: true,
			expectedConfig: &aiv1beta1.InfraEnv{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "infraEnv",
					Namespace: "cluster0",
				},
				Spec: aiv1beta1.InfraEnvSpec{
					ClusterRef: &aiv1beta1.ClusterReference{
						Name:      "ocp-edge-cluster-0",
						Namespace: "cluster0",
					},
					NMStateConfigLabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"cluster0-nmstate-label-name": "cluster0-nmstate-label-value",
						},
					},
					PullSecretRef: &corev1.LocalObjectReference{
						Name: "pull-secret",
					},
					SSHAuthorizedKey: "ssh-rsa AAAAmyKey",
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "failed to unmarshal cluster-manifests/infraenv.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.InfraEnv",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-manifests/infraenv.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
		metadata:
		  name: infraEnv
		  namespace: cluster0
		spec:
		  wrongField: wrongValue`,
			expectedError: "failed to unmarshal cluster-manifests/infraenv.yaml: error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(infraEnvFilename).
				Return(
					&asset.File{
						Filename: infraEnvFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &InfraEnvFile{}
			found, err := asset.Load(fileFetcher)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in InfraEnv")
			}
		})
	}
}
