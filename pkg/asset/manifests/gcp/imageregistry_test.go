package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestBuildGCSConfig(t *testing.T) {
	cases := []struct {
		name           string
		infraID        string
		installConfig  *types.InstallConfig
		expectedBucket string
		expectedKeyID  string
		expectedRegion string
	}{
		{
			name:    "Without KMS encryption",
			infraID: "test-cluster-abc123",
			installConfig: &types.InstallConfig{
				Platform: types.Platform{
					GCP: &gcp.Platform{
						ProjectID: "test-project",
						Region:    "us-central1",
					},
				},
			},
			expectedBucket: "test-cluster-abc123-image-registry",
			expectedKeyID:  "",
			expectedRegion: "us-central1",
		},
		{
			name:    "With KMS encryption",
			infraID: "test-cluster-xyz789",
			installConfig: &types.InstallConfig{
				Platform: types.Platform{
					GCP: &gcp.Platform{
						ProjectID: "test-project",
						Region:    "us-east1",
						DefaultMachinePlatform: &gcp.MachinePool{
							OSDisk: gcp.OSDisk{
								EncryptionKey: &gcp.EncryptionKeyReference{
									KMSKey: &gcp.KMSKeyReference{
										Name:     "registry-key",
										KeyRing:  "openshift-keyring",
										Location: "us-east1",
									},
								},
							},
						},
					},
				},
			},
			expectedBucket: "test-cluster-xyz789-image-registry",
			expectedKeyID:  "projects/test-project/locations/us-east1/keyRings/openshift-keyring/cryptoKeys/registry-key",
			expectedRegion: "us-east1",
		},
		{
			name:    "With KMS encryption and custom project ID",
			infraID: "test-cluster-custom",
			installConfig: &types.InstallConfig{
				Platform: types.Platform{
					GCP: &gcp.Platform{
						ProjectID: "default-project",
						Region:    "europe-west1",
						DefaultMachinePlatform: &gcp.MachinePool{
							OSDisk: gcp.OSDisk{
								EncryptionKey: &gcp.EncryptionKeyReference{
									KMSKey: &gcp.KMSKeyReference{
										Name:      "registry-key",
										KeyRing:   "custom-keyring",
										Location:  "europe-west1",
										ProjectID: "kms-project",
									},
								},
							},
						},
					},
				},
			},
			expectedBucket: "test-cluster-custom-image-registry",
			expectedKeyID:  "projects/kms-project/locations/europe-west1/keyRings/custom-keyring/cryptoKeys/registry-key",
			expectedRegion: "europe-west1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gcsConfig := buildGCSConfig(tc.infraID, tc.installConfig)

			assert.Equal(t, tc.expectedBucket, gcsConfig.Bucket)
			assert.Equal(t, tc.expectedRegion, gcsConfig.Region)
			assert.Equal(t, tc.installConfig.Platform.GCP.ProjectID, gcsConfig.ProjectID)

			if tc.expectedKeyID == "" {
				assert.Empty(t, gcsConfig.KeyID)
			} else {
				assert.Equal(t, tc.expectedKeyID, gcsConfig.KeyID)
			}
		})
	}
}

func TestFormatKMSKeyResourcePath(t *testing.T) {
	cases := []struct {
		name         string
		kmsKey       *gcp.KMSKeyReference
		projectID    string
		expectedPath string
	}{
		{
			name:         "Nil KMS key returns empty string",
			kmsKey:       nil,
			projectID:    "test-project",
			expectedPath: "",
		},
		{
			name: "KMS key without project ID uses default project",
			kmsKey: &gcp.KMSKeyReference{
				Name:     "test-key",
				KeyRing:  "test-keyring",
				Location: "us-central1",
			},
			projectID:    "default-project",
			expectedPath: "projects/default-project/locations/us-central1/keyRings/test-keyring/cryptoKeys/test-key",
		},
		{
			name: "KMS key with project ID uses its own project",
			kmsKey: &gcp.KMSKeyReference{
				Name:      "test-key",
				KeyRing:   "test-keyring",
				Location:  "us-east1",
				ProjectID: "custom-project",
			},
			projectID:    "default-project",
			expectedPath: "projects/custom-project/locations/us-east1/keyRings/test-keyring/cryptoKeys/test-key",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := gcp.FormatKMSKeyResourcePath(tc.kmsKey, tc.projectID)
			assert.Equal(t, tc.expectedPath, path)
		})
	}
}

func TestImageRegistryConfigGeneration(t *testing.T) {
	cases := []struct {
		name                  string
		installConfig         *types.InstallConfig
		infraID               string
		expectKeyID           bool
		expectedKeyIDContains string
	}{
		{
			name: "Generate ImageRegistry config without KMS",
			installConfig: &types.InstallConfig{
				Platform: types.Platform{
					GCP: &gcp.Platform{
						ProjectID: "test-project",
						Region:    "us-central1",
					},
				},
			},
			infraID:     "test-cluster-abc",
			expectKeyID: false,
		},
		{
			name: "Generate ImageRegistry config with KMS",
			installConfig: &types.InstallConfig{
				Platform: types.Platform{
					GCP: &gcp.Platform{
						ProjectID: "test-project",
						Region:    "us-east1",
						DefaultMachinePlatform: &gcp.MachinePool{
							OSDisk: gcp.OSDisk{
								EncryptionKey: &gcp.EncryptionKeyReference{
									KMSKey: &gcp.KMSKeyReference{
										Name:     "registry-key",
										KeyRing:  "openshift-keyring",
										Location: "us-east1",
									},
								},
							},
						},
					},
				},
			},
			infraID:               "test-cluster-xyz",
			expectKeyID:           true,
			expectedKeyIDContains: "openshift-keyring/cryptoKeys/registry-key",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gcsConfig := buildGCSConfig(tc.infraID, tc.installConfig)

			// Verify the GCS config structure
			assert.NotNil(t, gcsConfig)
			assert.NotEmpty(t, gcsConfig.Bucket)
			assert.NotEmpty(t, gcsConfig.Region)
			assert.NotEmpty(t, gcsConfig.ProjectID)

			// Verify keyID presence
			if tc.expectKeyID {
				assert.NotEmpty(t, gcsConfig.KeyID, "Expected KeyID to be set")
				assert.Contains(t, gcsConfig.KeyID, tc.expectedKeyIDContains)
			} else {
				assert.Empty(t, gcsConfig.KeyID, "Expected KeyID to be empty")
			}

			// Create a full ImageRegistry CR structure to verify it can be marshaled
			configMap := map[string]interface{}{
				"apiVersion": "imageregistry.operator.openshift.io/v1",
				"kind":       "Config",
				"metadata": map[string]interface{}{
					"name": "cluster",
				},
				"spec": map[string]interface{}{
					"managementState": "Managed",
					"storage": map[string]interface{}{
						"gcs": gcsConfig,
					},
					"replicas": 2,
				},
			}

			// Verify we can marshal it to YAML without errors
			data, err := yaml.Marshal(configMap)
			assert.NoError(t, err)
			assert.NotEmpty(t, data)

			// Unmarshal and verify structure
			var unmarshaled map[string]interface{}
			err = yaml.Unmarshal(data, &unmarshaled)
			assert.NoError(t, err)

			if assert.NotNil(t, unmarshaled) {
				spec, ok := unmarshaled["spec"].(map[string]interface{})
				if assert.True(t, ok) {
					storage, ok := spec["storage"].(map[string]interface{})
					if assert.True(t, ok) {
						gcs, ok := storage["gcs"].(map[string]interface{})
						if assert.True(t, ok) && tc.expectKeyID {
							_, ok := gcs["keyID"]
							assert.True(t, ok, "Unmarshaled YAML should contain keyID")
						}
					}
				}
			}
		})
	}
}
