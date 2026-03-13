package vsphere

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestLoadCredentialsFile(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string
		fileMode      os.FileMode
		expectError   bool
		expectNil     bool
		expectedCreds map[string]*vsphere.VCenterComponentCredentials
	}{
		{
			name:        "file does not exist",
			expectNil:   true,
			expectError: false,
		},
		{
			name:        "insecure file permissions",
			fileContent: "[vcenter1.example.com]\nuser = admin\n",
			fileMode:    0644,
			expectError: true,
		},
		{
			name: "valid single vCenter with all components",
			fileContent: `[vcenter1.example.com]
user = installer@vsphere.local
password = installer-pass
machine-api.user = machine-api@vsphere.local
machine-api.password = machine-api-pass
csi-driver.user = csi@vsphere.local
csi-driver.password = csi-pass
cloud-controller.user = ccm@vsphere.local
cloud-controller.password = ccm-pass
diagnostics.user = diagnostics@vsphere.local
diagnostics.password = diagnostics-pass
`,
			fileMode:    0600,
			expectError: false,
			expectedCreds: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "machine-api@vsphere.local",
						Password: "machine-api-pass",
					},
					CSIDriver: &vsphere.VCenterCredential{
						User:     "csi@vsphere.local",
						Password: "csi-pass",
					},
					CloudController: &vsphere.VCenterCredential{
						User:     "ccm@vsphere.local",
						Password: "ccm-pass",
					},
					Diagnostics: &vsphere.VCenterCredential{
						User:     "diagnostics@vsphere.local",
						Password: "diagnostics-pass",
					},
				},
			},
		},
		{
			name: "multiple vCenters",
			fileContent: `[vcenter1.example.com]
machine-api.user = vc1-machine-api@vsphere.local
machine-api.password = vc1-machine-api-pass

[vcenter2.example.com]
machine-api.user = vc2-machine-api@vsphere.local
machine-api.password = vc2-machine-api-pass
csi-driver.user = vc2-csi@vsphere.local
csi-driver.password = vc2-csi-pass
`,
			fileMode:    0600,
			expectError: false,
			expectedCreds: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "vc1-machine-api@vsphere.local",
						Password: "vc1-machine-api-pass",
					},
				},
				"vcenter2.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "vc2-machine-api@vsphere.local",
						Password: "vc2-machine-api-pass",
					},
					CSIDriver: &vsphere.VCenterCredential{
						User:     "vc2-csi@vsphere.local",
						Password: "vc2-csi-pass",
					},
				},
			},
		},
		{
			name: "partial component credentials",
			fileContent: `[vcenter1.example.com]
machine-api.user = machine-api@vsphere.local
machine-api.password = machine-api-pass
`,
			fileMode:    0600,
			expectError: false,
			expectedCreds: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "machine-api@vsphere.local",
						Password: "machine-api-pass",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory and file if fileContent is provided
			if tt.fileContent != "" {
				tmpDir := t.TempDir()
				credFile := filepath.Join(tmpDir, "credentials")

				err := os.WriteFile(credFile, []byte(tt.fileContent), tt.fileMode)
				require.NoError(t, err)

				// Set environment variable to use this file
				t.Setenv(CredentialsFileEnvVar, credFile)
			}

			// Load credentials
			creds, err := LoadCredentialsFile()

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expectNil {
				assert.Nil(t, creds)
				return
			}

			require.NotNil(t, creds)
			assert.Equal(t, len(tt.expectedCreds), len(creds))

			for vcenter, expectedCompCreds := range tt.expectedCreds {
				actualCompCreds, exists := creds[vcenter]
				require.True(t, exists, "vCenter %s not found in loaded credentials", vcenter)

				if expectedCompCreds.MachineAPI != nil {
					require.NotNil(t, actualCompCreds.MachineAPI)
					assert.Equal(t, expectedCompCreds.MachineAPI.User, actualCompCreds.MachineAPI.User)
					assert.Equal(t, expectedCompCreds.MachineAPI.Password, actualCompCreds.MachineAPI.Password)
				} else {
					assert.Nil(t, actualCompCreds.MachineAPI)
				}

				if expectedCompCreds.CSIDriver != nil {
					require.NotNil(t, actualCompCreds.CSIDriver)
					assert.Equal(t, expectedCompCreds.CSIDriver.User, actualCompCreds.CSIDriver.User)
					assert.Equal(t, expectedCompCreds.CSIDriver.Password, actualCompCreds.CSIDriver.Password)
				} else {
					assert.Nil(t, actualCompCreds.CSIDriver)
				}

				if expectedCompCreds.CloudController != nil {
					require.NotNil(t, actualCompCreds.CloudController)
					assert.Equal(t, expectedCompCreds.CloudController.User, actualCompCreds.CloudController.User)
					assert.Equal(t, expectedCompCreds.CloudController.Password, actualCompCreds.CloudController.Password)
				} else {
					assert.Nil(t, actualCompCreds.CloudController)
				}

				if expectedCompCreds.Diagnostics != nil {
					require.NotNil(t, actualCompCreds.Diagnostics)
					assert.Equal(t, expectedCompCreds.Diagnostics.User, actualCompCreds.Diagnostics.User)
					assert.Equal(t, expectedCompCreds.Diagnostics.Password, actualCompCreds.Diagnostics.Password)
				} else {
					assert.Nil(t, actualCompCreds.Diagnostics)
				}
			}
		})
	}
}

func TestMergeCredentials(t *testing.T) {
	tests := []struct {
		name           string
		vcenters       []vsphere.VCenter
		fileCredentials map[string]*vsphere.VCenterComponentCredentials
		expected       []vsphere.VCenter
	}{
		{
			name: "install-config has no component credentials, file has credentials",
			vcenters: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
				},
			},
			fileCredentials: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "machine-api@vsphere.local",
						Password: "machine-api-pass",
					},
				},
			},
			expected: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
					ComponentCredentials: &vsphere.VCenterComponentCredentials{
						MachineAPI: &vsphere.VCenterCredential{
							User:     "machine-api@vsphere.local",
							Password: "machine-api-pass",
						},
					},
				},
			},
		},
		{
			name: "install-config has component credentials, file has credentials (install-config wins)",
			vcenters: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
					ComponentCredentials: &vsphere.VCenterComponentCredentials{
						MachineAPI: &vsphere.VCenterCredential{
							User:     "config-machine-api@vsphere.local",
							Password: "config-machine-api-pass",
						},
					},
				},
			},
			fileCredentials: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "file-machine-api@vsphere.local",
						Password: "file-machine-api-pass",
					},
				},
			},
			expected: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
					ComponentCredentials: &vsphere.VCenterComponentCredentials{
						MachineAPI: &vsphere.VCenterCredential{
							User:     "config-machine-api@vsphere.local",
							Password: "config-machine-api-pass",
						},
					},
				},
			},
		},
		{
			name: "partial merge - install-config has some components, file has others",
			vcenters: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
					ComponentCredentials: &vsphere.VCenterComponentCredentials{
						MachineAPI: &vsphere.VCenterCredential{
							User:     "config-machine-api@vsphere.local",
							Password: "config-machine-api-pass",
						},
					},
				},
			},
			fileCredentials: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter1.example.com": {
					CSIDriver: &vsphere.VCenterCredential{
						User:     "file-csi@vsphere.local",
						Password: "file-csi-pass",
					},
				},
			},
			expected: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
					ComponentCredentials: &vsphere.VCenterComponentCredentials{
						MachineAPI: &vsphere.VCenterCredential{
							User:     "config-machine-api@vsphere.local",
							Password: "config-machine-api-pass",
						},
						CSIDriver: &vsphere.VCenterCredential{
							User:     "file-csi@vsphere.local",
							Password: "file-csi-pass",
						},
					},
				},
			},
		},
		{
			name: "nil file credentials",
			vcenters: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
				},
			},
			fileCredentials: nil,
			expected: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
				},
			},
		},
		{
			name: "file has credentials for non-existent vCenter",
			vcenters: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
				},
			},
			fileCredentials: map[string]*vsphere.VCenterComponentCredentials{
				"vcenter2.example.com": {
					MachineAPI: &vsphere.VCenterCredential{
						User:     "machine-api@vsphere.local",
						Password: "machine-api-pass",
					},
				},
			},
			expected: []vsphere.VCenter{
				{
					Server:   "vcenter1.example.com",
					Username: "admin@vsphere.local",
					Password: "admin-pass",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test data
			vcenters := make([]vsphere.VCenter, len(tt.vcenters))
			copy(vcenters, tt.vcenters)

			// Merge credentials
			MergeCredentials(vcenters, tt.fileCredentials)

			// Verify results
			require.Equal(t, len(tt.expected), len(vcenters))

			for i, expected := range tt.expected {
				actual := vcenters[i]
				assert.Equal(t, expected.Server, actual.Server)
				assert.Equal(t, expected.Username, actual.Username)
				assert.Equal(t, expected.Password, actual.Password)

				if expected.ComponentCredentials == nil {
					assert.Nil(t, actual.ComponentCredentials)
				} else {
					require.NotNil(t, actual.ComponentCredentials)

					if expected.ComponentCredentials.MachineAPI != nil {
						require.NotNil(t, actual.ComponentCredentials.MachineAPI)
						assert.Equal(t, expected.ComponentCredentials.MachineAPI.User, actual.ComponentCredentials.MachineAPI.User)
						assert.Equal(t, expected.ComponentCredentials.MachineAPI.Password, actual.ComponentCredentials.MachineAPI.Password)
					} else {
						assert.Nil(t, actual.ComponentCredentials.MachineAPI)
					}

					if expected.ComponentCredentials.CSIDriver != nil {
						require.NotNil(t, actual.ComponentCredentials.CSIDriver)
						assert.Equal(t, expected.ComponentCredentials.CSIDriver.User, actual.ComponentCredentials.CSIDriver.User)
						assert.Equal(t, expected.ComponentCredentials.CSIDriver.Password, actual.ComponentCredentials.CSIDriver.Password)
					} else {
						assert.Nil(t, actual.ComponentCredentials.CSIDriver)
					}

					if expected.ComponentCredentials.CloudController != nil {
						require.NotNil(t, actual.ComponentCredentials.CloudController)
						assert.Equal(t, expected.ComponentCredentials.CloudController.User, actual.ComponentCredentials.CloudController.User)
						assert.Equal(t, expected.ComponentCredentials.CloudController.Password, actual.ComponentCredentials.CloudController.Password)
					} else {
						assert.Nil(t, actual.ComponentCredentials.CloudController)
					}

					if expected.ComponentCredentials.Diagnostics != nil {
						require.NotNil(t, actual.ComponentCredentials.Diagnostics)
						assert.Equal(t, expected.ComponentCredentials.Diagnostics.User, actual.ComponentCredentials.Diagnostics.User)
						assert.Equal(t, expected.ComponentCredentials.Diagnostics.Password, actual.ComponentCredentials.Diagnostics.Password)
					} else {
						assert.Nil(t, actual.ComponentCredentials.Diagnostics)
					}
				}
			}
		})
	}
}
