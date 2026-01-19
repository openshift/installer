package manifests

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
)

// getValidOptionalInstallConfigWithFencing returns an install config with fencing credentials
// for testing TNF (Two-Node Fencing) cluster configurations.
func getValidOptionalInstallConfigWithFencing() *agent.OptionalInstallConfig {
	installConfig := getValidOptionalInstallConfig()
	// Set to 2 replicas for TNF
	installConfig.Config.ControlPlane.Replicas = ptr.To(int64(2))
	// Add fencing credentials
	installConfig.Config.ControlPlane.Fencing = &types.Fencing{
		Credentials: []*types.Credential{
			{
				HostName:                "master-0",
				Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
				Username:                "admin",
				Password:                "password123",
				CertificateVerification: types.CertificateVerificationDisabled,
			},
			{
				HostName:                "master-1",
				Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/def",
				Username:                "admin",
				Password:                "password456",
				CertificateVerification: types.CertificateVerificationEnabled,
			},
		},
	}
	return installConfig
}

// getValidOptionalInstallConfigWithFencingNoCertVerification returns an install config
// with fencing credentials where certificateVerification is not set.
func getValidOptionalInstallConfigWithFencingNoCertVerification() *agent.OptionalInstallConfig {
	installConfig := getValidOptionalInstallConfig()
	installConfig.Config.ControlPlane.Replicas = ptr.To(int64(2))
	installConfig.Config.ControlPlane.Fencing = &types.Fencing{
		Credentials: []*types.Credential{
			{
				HostName: "master-0",
				Address:  "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
				Username: "admin",
				Password: "password123",
				// CertificateVerification omitted
			},
		},
	}
	return installConfig
}

// getValidOptionalInstallConfigWithEmptyFencingCredentials returns an install config
// with an empty fencing credentials array.
func getValidOptionalInstallConfigWithEmptyFencingCredentials() *agent.OptionalInstallConfig {
	installConfig := getValidOptionalInstallConfig()
	installConfig.Config.ControlPlane.Replicas = ptr.To(int64(2))
	installConfig.Config.ControlPlane.Fencing = &types.Fencing{
		Credentials: []*types.Credential{},
	}
	return installConfig
}

func TestFencingCredentials_Generate(t *testing.T) {
	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *FencingCredentialsConfig
		expectedFiles  int
	}{
		{
			name: "valid fencing credentials - install workflow",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigWithFencing(),
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
			},
			expectedConfig: &FencingCredentialsConfig{
				Credentials: []*types.Credential{
					{
						HostName:                "master-0",
						Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
						Username:                "admin",
						Password:                "password123",
						CertificateVerification: types.CertificateVerificationDisabled,
					},
					{
						HostName:                "master-1",
						Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/def",
						Username:                "admin",
						Password:                "password456",
						CertificateVerification: types.CertificateVerificationEnabled,
					},
				},
			},
			expectedFiles: 1,
		},
		{
			name: "no fencing credentials - non-TNF cluster",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(), // No fencing configured
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
			},
			expectedConfig: nil,
			expectedFiles:  0,
		},
		{
			name: "add-nodes workflow - no fencing",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigWithFencing(),
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
			},
			expectedConfig: nil,
			expectedFiles:  0,
		},
		{
			name: "nil install config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
			},
			expectedConfig: nil,
			expectedFiles:  0,
		},
		{
			name: "empty fencing credentials array",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigWithEmptyFencingCredentials(),
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
			},
			expectedConfig: nil,
			expectedFiles:  0,
		},
		{
			name: "fencing without certificateVerification field",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfigWithFencingNoCertVerification(),
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
			},
			expectedConfig: &FencingCredentialsConfig{
				Credentials: []*types.Credential{
					{
						HostName: "master-0",
						Address:  "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
						Username: "admin",
						Password: "password123",
					},
				},
			},
			expectedFiles: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			fencingAsset := &FencingCredentials{}
			err := fencingAsset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, fencingAsset.Config)
				assert.Equal(t, tc.expectedFiles, len(fencingAsset.Files()))

				if tc.expectedFiles > 0 {
					configFile := fencingAsset.Files()[0]
					assert.Equal(t, "cluster-manifests/fencing-credentials.yaml", configFile.Filename)

					// Verify YAML contains expected structure
					assert.Contains(t, string(configFile.Data), "credentials:")
					assert.Contains(t, string(configFile.Data), "hostname:")
					assert.Contains(t, string(configFile.Data), "address:")
					assert.Contains(t, string(configFile.Data), "username:")
					assert.Contains(t, string(configFile.Data), "password:")
				}
			}
		})
	}
}

func TestFencingCredentials_Load(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *FencingCredentialsConfig
	}{
		{
			name: "valid fencing credentials file",
			data: `credentials:
- hostname: master-0
  address: redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc
  username: admin
  password: password123
  certificateVerification: Disabled
- hostname: master-1
  address: redfish+https://192.168.111.1:8000/redfish/v1/Systems/def
  username: admin
  password: password456
  certificateVerification: Enabled`,
			expectedFound: true,
			expectedConfig: &FencingCredentialsConfig{
				Credentials: []*types.Credential{
					{
						HostName:                "master-0",
						Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
						Username:                "admin",
						Password:                "password123",
						CertificateVerification: types.CertificateVerificationDisabled,
					},
					{
						HostName:                "master-1",
						Address:                 "redfish+https://192.168.111.1:8000/redfish/v1/Systems/def",
						Username:                "admin",
						Password:                "password456",
						CertificateVerification: types.CertificateVerificationEnabled,
					},
				},
			},
		},
		{
			name: "fencing credentials without certificateVerification",
			data: `credentials:
- hostname: master-0
  address: redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc
  username: admin
  password: password123`,
			expectedFound: true,
			expectedConfig: &FencingCredentialsConfig{
				Credentials: []*types.Credential{
					{
						HostName: "master-0",
						Address:  "redfish+https://192.168.111.1:8000/redfish/v1/Systems/abc",
						Username: "admin",
						Password: "password123",
					},
				},
			},
		},
		{
			name:          "invalid yaml",
			data:          `this is not valid yaml: [`,
			expectedFound: false,
			expectedError: "failed to unmarshal cluster-manifests/fencing-credentials.yaml",
		},
		{
			name:          "file not found",
			fetchError:    &os.PathError{Err: os.ErrNotExist},
			expectedFound: false,
		},
		{
			name:          "error fetching file",
			fetchError:    errors.New("fetch failed"),
			expectedFound: false,
			expectedError: "failed to load cluster-manifests/fencing-credentials.yaml file: fetch failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(fencingCredentialsFilename).
				Return(
					&asset.File{
						Filename: fencingCredentialsFilename,
						Data:     []byte(tc.data),
					},
					tc.fetchError,
				)

			fencingAsset := &FencingCredentials{}
			found, err := fencingAsset.Load(fileFetcher)

			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, fencingAsset.Config)
			}
		})
	}
}

func TestFencingCredentials_YAMLFormat(t *testing.T) {
	// This test verifies that the YAML format produced by the installer
	// matches what assisted-service expects to consume.
	parents := asset.Parents{}
	parents.Add(
		getValidOptionalInstallConfigWithFencing(),
		&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
	)

	fencingAsset := &FencingCredentials{}
	err := fencingAsset.Generate(context.Background(), parents)
	assert.NoError(t, err)
	assert.NotEmpty(t, fencingAsset.Files())

	yamlContent := string(fencingAsset.Files()[0].Data)

	// Verify YAML uses lowercase field names (yaml tags, not JSON tags)
	assert.Contains(t, yamlContent, "hostname: master-0")
	assert.Contains(t, yamlContent, "hostname: master-1")
	assert.Contains(t, yamlContent, "username: admin")
	assert.Contains(t, yamlContent, "password: password123")
	assert.Contains(t, yamlContent, "password: password456")
	assert.Contains(t, yamlContent, "certificateVerification: Disabled")
	assert.Contains(t, yamlContent, "certificateVerification: Enabled")

	// Verify it does NOT use camelCase (JSON tag format)
	assert.NotContains(t, yamlContent, "hostName:")
}

func TestFencingCredentials_Name(t *testing.T) {
	fencingAsset := &FencingCredentials{}
	assert.Equal(t, "Fencing Credentials", fencingAsset.Name())
}

func TestFencingCredentials_Dependencies(t *testing.T) {
	fencingAsset := &FencingCredentials{}
	deps := fencingAsset.Dependencies()

	assert.Len(t, deps, 2)

	// Check for expected dependency types
	hasWorkflow := false
	hasInstallConfig := false
	for _, dep := range deps {
		switch dep.(type) {
		case *workflow.AgentWorkflow:
			hasWorkflow = true
		case *agent.OptionalInstallConfig:
			hasInstallConfig = true
		}
	}
	assert.True(t, hasWorkflow, "should depend on AgentWorkflow")
	assert.True(t, hasInstallConfig, "should depend on OptionalInstallConfig")
}
