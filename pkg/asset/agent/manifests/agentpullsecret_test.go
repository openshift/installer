package manifests

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestAgentPullSecret_Generate(t *testing.T) {
	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *corev1.Secret
	}{
		{
			name: "retrieved from cluster",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					ClusterName: "ostest",
					Namespace:   "cluster0",
					PullSecret:  "{\n  \"auths\": {\n    \"cloud.openshift.com\": {\n      \"auth\": \"b3BlUTA=\",\n      \"email\": \"test@redhat.com\"\n    }\n  }\n}",
				},
			},
			expectedConfig: &corev1.Secret{
				TypeMeta: v1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "ostest-pull-secret",
					Namespace: "cluster0",
				},
				StringData: map[string]string{
					".dockerconfigjson": "{\n  \"auths\": {\n    \"cloud.openshift.com\": {\n      \"auth\": \"b3BlUTA=\",\n      \"email\": \"test@redhat.com\"\n    }\n  }\n}",
				},
			},
		},
		{
			name: "missing install config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{}, &workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall}, &joiner.ClusterInfo{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(), &workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall}, &joiner.ClusterInfo{},
			},
			expectedConfig: &corev1.Secret{
				TypeMeta: v1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      getPullSecretName(getValidOptionalInstallConfig().ClusterName()),
					Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
				},
				StringData: map[string]string{
					".dockerconfigjson": "{\n  \"auths\": {\n    \"cloud.openshift.com\": {\n      \"auth\": \"b3BlUTA=\",\n      \"email\": \"test@redhat.com\"\n    }\n  }\n}",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &AgentPullSecret{}
			err := asset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/pull-secret.yaml", configFile.Filename)

				var actualConfig corev1.Secret
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestAgentPullSecret_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
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
  .dockerconfigjson: '{"auths":{"cloud.test":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}'`,
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
					".dockerconfigjson": "{\n  \"auths\": {\n    \"cloud.test\": {\n      \"auth\": \"c3VwZXItc2VjcmV0Cg==\"\n    }\n  }\n}",
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "failed to unmarshal cluster-manifests/pull-secret.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1.Secret",
		},
		{
			name:          "empty",
			data:          "",
			expectedError: "invalid PullSecret configuration: stringData: Required value: the pull secret is empty",
		},
		{
			name: "missing-string-data",
			data: `
apiVersion: v1
kind: Secret
metadata:
  name: pull-secret
  namespace: cluster-0`,
			expectedError: "invalid PullSecret configuration: stringData: Required value: the pull secret is empty",
		},
		{
			name: "missing-secret-key",
			data: `
apiVersion: v1
kind: Secret
metadata:
  name: pull-secret
  namespace: cluster-0
stringData:
  .dockerconfigjson:`,
			expectedError: "invalid PullSecret configuration: stringData: Required value: the pull secret does not contain any data",
		},
		{
			name: "bad-secret-format",
			data: `
apiVersion: v1
kind: Secret
metadata:
  name: pull-secret
  namespace: cluster-0
stringData:
  .dockerconfigjson: 'foo'`,
			expectedError: "invalid PullSecret configuration: stringData: Invalid value: \"foo\": invalid character 'o' in literal false (expecting 'a')",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-manifests/pull-secret.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
		metadata:
		  name: pull-secret
		  namespace: cluster0
		spec:
		  wrongField: wrongValue`,
			expectedError: "failed to unmarshal cluster-manifests/pull-secret.yaml: error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
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
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in AgentPullSecret")
			}
		})
	}

}
