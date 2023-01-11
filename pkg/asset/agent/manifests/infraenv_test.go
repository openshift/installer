package manifests

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
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
			name: "missing-config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
				&agentconfig.AgentConfig{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(),
				getValidAgentConfig(),
			},
			expectedConfig: &aiv1beta1.InfraEnv{
				ObjectMeta: metav1.ObjectMeta{
					Name:      getInfraEnvName(getValidOptionalInstallConfig()),
					Namespace: getObjectMetaNamespace(getValidOptionalInstallConfig()),
				},
				Spec: aiv1beta1.InfraEnvSpec{
					ClusterRef: &aiv1beta1.ClusterReference{
						Name:      getClusterDeploymentName(getValidOptionalInstallConfig()),
						Namespace: getObjectMetaNamespace(getValidOptionalInstallConfig()),
					},
					SSHAuthorizedKey: strings.Trim(testSSHKey, "|\n\t"),
					PullSecretRef: &corev1.LocalObjectReference{
						Name: getPullSecretName(getValidOptionalInstallConfig()),
					},
					NMStateConfigLabelSelector: metav1.LabelSelector{
						MatchLabels: getNMStateConfigLabels(getValidOptionalInstallConfig()),
					},
				},
			},
		},
		{
			name: "proxy valid configuration",
			dependencies: []asset.Asset{
				getProxyValidOptionalInstallConfig(),
				getValidAgentConfig(),
			},
			expectedConfig: &aiv1beta1.InfraEnv{
				ObjectMeta: metav1.ObjectMeta{
					Name:      getClusterDeploymentName(getProxyValidOptionalInstallConfig()),
					Namespace: getObjectMetaNamespace(getProxyValidOptionalInstallConfig()),
				},
				Spec: aiv1beta1.InfraEnvSpec{
					Proxy:            getProxy(getProxyValidOptionalInstallConfig()),
					SSHAuthorizedKey: strings.Trim(testSSHKey, "|\n\t"),
					PullSecretRef: &corev1.LocalObjectReference{
						Name: getPullSecretName(getProxyValidOptionalInstallConfig()),
					},
					NMStateConfigLabelSelector: metav1.LabelSelector{
						MatchLabels: getNMStateConfigLabels(getProxyValidOptionalInstallConfig()),
					},
					ClusterRef: &aiv1beta1.ClusterReference{
						Name:      getClusterDeploymentName(getProxyValidOptionalInstallConfig()),
						Namespace: getObjectMetaNamespace(getProxyValidOptionalInstallConfig()),
					},
				},
			},
		},
		{
			name: "Additional NTP sources",
			dependencies: []asset.Asset{
				getProxyValidOptionalInstallConfig(),
				getValidAgentConfigWithAdditionalNTPSources(),
			},
			expectedConfig: &aiv1beta1.InfraEnv{
				ObjectMeta: metav1.ObjectMeta{
					Name:      getClusterDeploymentName(getProxyValidOptionalInstallConfig()),
					Namespace: getObjectMetaNamespace(getProxyValidOptionalInstallConfig()),
				},
				Spec: aiv1beta1.InfraEnvSpec{
					Proxy:            getProxy(getProxyValidOptionalInstallConfig()),
					SSHAuthorizedKey: strings.Trim(testSSHKey, "|\n\t"),
					PullSecretRef: &corev1.LocalObjectReference{
						Name: getPullSecretName(getProxyValidOptionalInstallConfig()),
					},
					NMStateConfigLabelSelector: metav1.LabelSelector{
						MatchLabels: getNMStateConfigLabels(getProxyValidOptionalInstallConfig()),
					},
					ClusterRef: &aiv1beta1.ClusterReference{
						Name:      getClusterDeploymentName(getProxyValidOptionalInstallConfig()),
						Namespace: getObjectMetaNamespace(getProxyValidOptionalInstallConfig()),
					},
					AdditionalNTPSources: getValidAgentConfigWithAdditionalNTPSources().Config.AdditionalNTPSources,
				},
			},
		},
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
