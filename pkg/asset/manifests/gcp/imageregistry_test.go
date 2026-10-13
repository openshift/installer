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
