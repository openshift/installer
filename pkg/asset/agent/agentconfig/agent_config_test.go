package agentconfig

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/agent"
)

func TestAgentConfig_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedError  string
		expectedFound  bool
		expectedConfig *AgentConfigBuilder
	}{
		{
			name: "valid-config-single-node",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
bootArtifactsBaseURL: http://user-specified-pxe-infra.com`,
			expectedFound:  true,
			expectedConfig: agentConfig().bootArtifactsBaseURL("http://user-specified-pxe-infra.com"),
		},
		{
			name: "not-yaml",
			data: `This is not a yaml file`,

			expectedFound: false,
			expectedError: "failed to unmarshal agent-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type agent.Config",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},

			expectedFound: false,
		},
		{
			name:       "error-fetching-file",
			fetchError: errors.New("fetch failed"),

			expectedFound: false,
			expectedError: "failed to load agent-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-wrong
wrongField: wrongValue`,

			expectedFound: false,
			expectedError: "failed to unmarshal agent-config.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "empty-rendezvousIP",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0`,

			expectedFound:  true,
			expectedConfig: agentConfig().rendezvousIP(""),
		},
		{
			name: "invalid-rendezvousIP",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: not-a-valid-ip`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: rendezvousIP: Invalid value: \"not-a-valid-ip\": \"not-a-valid-ip\" is not a valid IP",
		},
		{
			name: "empty-bootArtifactsBaseURL",
			data: `
apiVersion: v1alpha1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80`,

			expectedFound:  true,
			expectedConfig: agentConfig(),
		},
		{
			name: "invalid-bootArtifactsBaseURL",
			data: `
apiVersion: v1alpha1
metadata:
  name: agent-config-cluster0
bootArtifactsBaseURL: not-a-valid-url`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: bootArtifactsBaseURL: Invalid value: \"not-a-valid-url\": invalid URI \"not-a-valid-url\" (no scheme)",
		},
		{
			name: "invalid-additionalNTPSourceDomain",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
additionalNTPSources:
  - 0.fedora.pool.ntp.org
  - 1.fedora.pool.ntp.org
  - 192.168.111.14
  - fd10:39:192:1::1337
  - invalid_pool.ntp.org
rendezvousIP: 192.168.111.80`,
			expectedFound: false,
			expectedError: "invalid Agent Config configuration: additionalNTPSources: Invalid value: \"invalid_pool.ntp.org\": NTP server is not a valid IP or domain name",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(agentConfigFilename).
				Return(
					&asset.File{
						Filename: agentConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &AgentConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig.build(), asset.Config, "unexpected Config in AgentConfig")
				}
			}
		})
	}
}

func TestAgentConfig_PersistToFile(t *testing.T) {
	validCustomConfig := `# custom comment
apiVersion: v1beta1
kind: AgentConfig
metadata:
  name: test-cluster
rendezvousIP: 192.168.1.100
`

	missingAPIVersionConfig := `# intentionally invalid: missing apiVersion
kind: AgentConfig
metadata:
  name: broken
`

	t.Run("generated asset persists File.Data with apiVersion", func(t *testing.T) {
		dir := t.TempDir()
		a := &AgentConfig{}
		require.NoError(t, a.Generate(context.Background(), asset.Parents{}))
		require.NotNil(t, a.File)
		require.Contains(t, string(a.File.Data), "apiVersion:")

		require.NoError(t, a.PersistToFile(dir))

		got, err := os.ReadFile(filepath.Join(dir, agentConfigFilename))
		require.NoError(t, err)
		assert.Equal(t, a.File.Data, got)
	})

	t.Run("loaded custom config leaves Template empty and persists File.Data", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, agentConfigFilename)
		require.NoError(t, os.WriteFile(path, []byte(validCustomConfig), 0600))

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		fileFetcher := mock.NewMockFileFetcher(mockCtrl)
		fileFetcher.EXPECT().FetchByName(agentConfigFilename).Return(
			&asset.File{Filename: agentConfigFilename, Data: []byte(validCustomConfig)},
			nil,
		)

		a := &AgentConfig{}
		found, err := a.Load(fileFetcher)
		require.NoError(t, err)
		require.True(t, found)
		assert.Empty(t, a.Template)
		require.NotNil(t, a.File)
		assert.Equal(t, []byte(validCustomConfig), a.File.Data)

		require.NoError(t, a.PersistToFile(dir))

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, validCustomConfig, string(got))
	})

	t.Run("nil File does not create or change files", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, agentConfigFilename)
		existing := []byte("existing-content\n")
		require.NoError(t, os.WriteFile(path, existing, 0600))

		a := &AgentConfig{}
		require.NoError(t, a.PersistToFile(dir))

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, existing, got)

		emptyDir := t.TempDir()
		require.NoError(t, a.PersistToFile(emptyDir))
		_, err = os.Stat(filepath.Join(emptyDir, agentConfigFilename))
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("empty config load fails and PersistToFile leaves file unchanged", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, agentConfigFilename)
		emptyContent := []byte("")
		require.NoError(t, os.WriteFile(path, emptyContent, 0600))

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		fileFetcher := mock.NewMockFileFetcher(mockCtrl)
		fileFetcher.EXPECT().FetchByName(agentConfigFilename).Return(
			&asset.File{Filename: agentConfigFilename, Data: emptyContent},
			nil,
		)

		a := &AgentConfig{}
		found, err := a.Load(fileFetcher)
		assert.False(t, found)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "apiVersion")
		assert.Nil(t, a.File)

		require.NoError(t, a.PersistToFile(dir))

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, emptyContent, got)
	})

	t.Run("nonempty missing apiVersion load fails and PersistToFile leaves file unchanged", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, agentConfigFilename)
		require.NoError(t, os.WriteFile(path, []byte(missingAPIVersionConfig), 0600))

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		fileFetcher := mock.NewMockFileFetcher(mockCtrl)
		fileFetcher.EXPECT().FetchByName(agentConfigFilename).Return(
			&asset.File{Filename: agentConfigFilename, Data: []byte(missingAPIVersionConfig)},
			nil,
		)

		a := &AgentConfig{}
		found, err := a.Load(fileFetcher)
		assert.False(t, found)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "apiVersion")
		assert.Nil(t, a.File)

		require.NoError(t, a.PersistToFile(dir))

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, missingAPIVersionConfig, string(got))
	})
}

// AgentConfigBuilder it's a builder class to make it easier creating agent.Config instance
// used in the test cases.
type AgentConfigBuilder struct {
	agent.Config
}

func agentConfig() *AgentConfigBuilder {
	return &AgentConfigBuilder{
		Config: agent.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "agent-config-cluster0",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: agent.AgentConfigVersion,
			},
			RendezvousIP: "192.168.111.80",
		},
	}
}

func (acb *AgentConfigBuilder) build() *agent.Config {
	return &acb.Config
}

func (acb *AgentConfigBuilder) rendezvousIP(ip string) *AgentConfigBuilder {
	acb.Config.RendezvousIP = ip
	return acb
}

func (acb *AgentConfigBuilder) bootArtifactsBaseURL(url string) *AgentConfigBuilder {
	acb.Config.BootArtifactsBaseURL = url
	return acb
}
