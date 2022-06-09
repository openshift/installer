package manifests

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestInfraEnv_Generate(t *testing.T) {

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *aiv1beta1.InfraEnv
	}{
		{
			name:          "missing-config",
			expectedError: "missing configuration or manifest file",
		},
		// {
		// 	name: "default",
		// 	dependencies: []asset.Asset{
		// 		&installconfig.InstallConfig{
		// 			Config: &types.InstallConfig{
		// 				ObjectMeta: v1.ObjectMeta{
		// 					Name:      "ocp-edge-cluster-0",
		// 					Namespace: "cluster-0",
		// 				},
		// 				PullSecret: "secret-agent",
		// 				SSHKey:     "ssh-key",
		// 			},
		// 		},
		// 		&AgentPullSecret{},
		// 	},
		// 	expectedConfig: &aiv1beta1.InfraEnv{
		// 		ObjectMeta: v1.ObjectMeta{
		// 			Name:      "infraEnv",
		// 			Namespace: "cluster-0",
		// 		},
		// 		Spec: aiv1beta1.InfraEnvSpec{
		// 			ClusterRef: &aiv1beta1.ClusterReference{
		// 				Name:      "ocp-edge-cluster-0",
		// 				Namespace: "cluster-0",
		// 			},
		// 			SSHAuthorizedKey: "ssh-key",
		// 			PullSecretRef: &corev1.LocalObjectReference{
		// 				Name: "pull-secret",
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &InfraEnv{}
			err := asset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/infraenv.yaml", configFile.Filename)

				var actualConfig aiv1beta1.InfraEnv
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestInfraEnv_LoadedFromDisk(t *testing.T) {

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
				ObjectMeta: v1.ObjectMeta{
					Name:      "infraEnv",
					Namespace: "cluster0",
				},
				Spec: aiv1beta1.InfraEnvSpec{
					ClusterRef: &aiv1beta1.ClusterReference{
						Name:      "ocp-edge-cluster-0",
						Namespace: "cluster0",
					},
					NMStateConfigLabelSelector: v1.LabelSelector{
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
		{
			name: "empty-NMStateLabelSelector",
			data: `
metadata:
  name: infraEnv
  namespace: cluster0
spec:
  clusterRef:
    name: ocp-edge-cluster-0
    namespace: cluster0
  nmStateConfigLabelSelector: 
  pullSecretRef:
    name: pull-secret
  sshAuthorizedKey: |
    ssh-rsa AAAAmyKey`,
			expectedError: "invalid InfraEnv configuration: Spec.NMStateConfigLabelSelector.MatchLabels: Required value: at least one label must be set",
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

			asset := &InfraEnv{}
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
