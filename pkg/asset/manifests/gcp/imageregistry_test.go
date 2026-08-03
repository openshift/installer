package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/gcp"
)

func TestBuildGCSConfig(t *testing.T) {
	cases := []struct {
		name           string
		platform       *gcp.Platform
		expectedKeyID  string
		expectedRegion string
	}{
		{
			name: "Without KMS encryption",
			platform: &gcp.Platform{
				ProjectID: "test-project",
				Region:    "us-central1",
			},
			expectedKeyID:  "",
			expectedRegion: "us-central1",
		},
		{
			name: "With KMS encryption",
			platform: &gcp.Platform{
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
			expectedKeyID:  "projects/test-project/locations/us-east1/keyRings/openshift-keyring/cryptoKeys/registry-key",
			expectedRegion: "us-east1",
		},
		{
			name: "With KMS encryption and custom project ID",
			platform: &gcp.Platform{
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
			expectedKeyID:  "projects/kms-project/locations/europe-west1/keyRings/custom-keyring/cryptoKeys/registry-key",
			expectedRegion: "europe-west1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gcsConfig := BuildGCSConfig(tc.platform)

			assert.Empty(t, gcsConfig.Bucket, "Bucket should be empty to let the operator generate it")
			assert.Equal(t, tc.expectedRegion, gcsConfig.Region)
			assert.Equal(t, tc.platform.ProjectID, gcsConfig.ProjectID)

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
