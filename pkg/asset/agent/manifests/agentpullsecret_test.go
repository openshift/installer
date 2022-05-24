package manifests

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAgentPullSecret_Generate(t *testing.T) {

	t.Skip("Skipping asset generation test")

	installconfigAsset := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "cluster0",
			},
			PullSecret: "secret-agent",
		},
	}

	parents := asset.Parents{}
	parents.Add(installconfigAsset)

	asset := &AgentPullSecret{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	pullSecretFile := asset.Files()[0]
	assert.Equal(t, "cluster-manifests/pull-secret.yaml", pullSecretFile.Filename)

	secret := asset.Config
	data, err := base64.StdEncoding.DecodeString(secret.StringData[".dockerconfigjson"])
	assert.NoError(t, err)
	assert.Equal(t, installconfigAsset.Config.PullSecret, string(data))
}

func TestAgentPullSecret_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  bool
		expectedConfig *corev1.Secret
	}{
		{
			name: "valid-config-file",
			data: `
apiVersion: v1
kind: Secret
metadata:
  name: pull-secret
  namespace: cluster-0
stringData:
  .dockerconfigjson: c3VwZXItc2VjcmV0Cg==`,
			expectedFound: true,
			expectedConfig: &corev1.Secret{
				TypeMeta: v1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "pull-secret",
					Namespace: "cluster-0",
				},
				StringData: map[string]string{
					".dockerconfigjson": "c3VwZXItc2VjcmV0Cg==",
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
			expectedConfig: &corev1.Secret{},
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
  name: pull-secret
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
			fileFetcher.EXPECT().FetchByName(agentPullSecretFilename).
				Return(
					&asset.File{
						Filename: agentPullSecretFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &AgentPullSecret{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in AgentPullSecret")
			}
		})
	}

}
