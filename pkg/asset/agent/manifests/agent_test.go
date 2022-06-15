package manifests

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestAgent_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig []*aiv1beta1.Agent
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: host-1
  namespace: cluster0
  annotations:
    macAddress: "11:22:33:44:55:66"
spec:
  clusterDeploymentName:
    name: cluster-0
    namespace: cluster0
  hostname: mydomain.org
  InstallerArgs: "--debug"
  IgnitionConfigOverrides: "my-ignition-override"`,
			expectedFound: true,
			expectedConfig: []*aiv1beta1.Agent{
				{
					ObjectMeta: v1.ObjectMeta{
						Name:      "host-1",
						Namespace: "cluster0",
						Annotations: map[string]string{
							"macAddress": "11:22:33:44:55:66",
						},
					},
					Spec: aiv1beta1.AgentSpec{
						ClusterDeploymentName: &aiv1beta1.ClusterReference{
							Name:      "cluster-0",
							Namespace: "cluster0",
						},
						Hostname:                "mydomain.org",
						InstallerArgs:           "--debug",
						IgnitionConfigOverrides: "my-ignition-override",
					},
				},
			},
		},
		{
			name: "valid-config-multiple-agents",
			data: `
metadata:
  name: host-1
  namespace: cluster0
  annotations:
    macAddress: "11:22:33:44:55:01"
spec:
  clusterDeploymentName:
    name: cluster-0
    namespace: cluster0
  hostname: 1-mydomain.org
  InstallerArgs: "--debug"
  IgnitionConfigOverrides: "my-ignition-override
---
metadata:
  name: host-2
  namespace: cluster0
  annotations:
    macAddress: "11:22:33:44:55:02"
spec:
  clusterDeploymentName:
    name: cluster-0
    namespace: cluster0
  hostname: 2-mydomain.org
  InstallerArgs: "--debug"
  IgnitionConfigOverrides: "my-ignition-override"`,
			expectedFound: true,
			expectedConfig: []*aiv1beta1.Agent{
				{
					ObjectMeta: v1.ObjectMeta{
						Name:      "host-1",
						Namespace: "cluster0",
						Annotations: map[string]string{
							"macAddress": "11:22:33:44:55:01",
						},
					},
					Spec: aiv1beta1.AgentSpec{
						ClusterDeploymentName: &aiv1beta1.ClusterReference{
							Name:      "cluster-0",
							Namespace: "cluster0",
						},
						Hostname:                "1-mydomain.org",
						InstallerArgs:           "--debug",
						IgnitionConfigOverrides: "my-ignition-override",
					},
				},
				{
					ObjectMeta: v1.ObjectMeta{
						Name:      "host-2",
						Namespace: "cluster0",
						Annotations: map[string]string{
							"macAddress": "11:22:33:44:55:02",
						},
					},
					Spec: aiv1beta1.AgentSpec{
						ClusterDeploymentName: &aiv1beta1.ClusterReference{
							Name:      "cluster-0",
							Namespace: "cluster0",
						},
						Hostname:                "2-mydomain.org",
						InstallerArgs:           "--debug",
						IgnitionConfigOverrides: "my-ignition-override",
					},
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "could not decode YAML for cluster-manifests/agent.yaml: Error reading multiple YAMLs: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.Agent",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-manifests/agent.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
		metadata:
		  name: host-1
		  namespace: cluster0
		spec:
		  wrongField: wrongValue`,
			expectedError: "could not decode YAML for cluster-manifests/agent.yaml: Error reading multiple YAMLs: error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
		},
		{
			name: "missing-annotation",
			data: `
metadata:
  name: host-1
  namespace: cluster0
spec:
  clusterDeploymentName:
    name: cluster-0
    namespace: cluster0
  hostname: mydomain.org
  InstallerArgs: "--debug"
  IgnitionConfigOverrides: "my-ignition-override"`,
			expectedError: "invalid Agent configuration: metadata.annotations.macAddress: Required value: macAddress annotation is missing",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(agentFilename).
				Return(
					&asset.File{
						Filename: agentFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &Agent{}
			found, err := asset.Load(fileFetcher)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in Agent")
			}
		})
	}

}
