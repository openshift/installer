package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/gcp"
)

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
				Name:     "bootstrap-key",
				KeyRing:  "bootstrap-keyring",
				Location: "us-central1",
			},
			projectID:    "default-project",
			expectedPath: "projects/default-project/locations/us-central1/keyRings/bootstrap-keyring/cryptoKeys/bootstrap-key",
		},
		{
			name: "KMS key with project ID overrides default project",
			kmsKey: &gcp.KMSKeyReference{
				Name:      "bootstrap-key",
				KeyRing:   "bootstrap-keyring",
				Location:  "us-east1",
				ProjectID: "custom-project",
			},
			projectID:    "default-project",
			expectedPath: "projects/custom-project/locations/us-east1/keyRings/bootstrap-keyring/cryptoKeys/bootstrap-key",
		},
		{
			name: "KMS key with different location",
			kmsKey: &gcp.KMSKeyReference{
				Name:     "europe-key",
				KeyRing:  "europe-keyring",
				Location: "europe-west1",
			},
			projectID:    "europe-project",
			expectedPath: "projects/europe-project/locations/europe-west1/keyRings/europe-keyring/cryptoKeys/europe-key",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := formatKMSKeyResourcePath(tc.kmsKey, tc.projectID)
			assert.Equal(t, tc.expectedPath, path)
		})
	}
}

func TestGetBootstrapStorageName(t *testing.T) {
	cases := []struct {
		name         string
		clusterID    string
		expectedName string
	}{
		{
			name:         "Standard cluster ID",
			clusterID:    "test-cluster-abc123",
			expectedName: "test-cluster-abc123-bootstrap-ignition",
		},
		{
			name:         "Short cluster ID",
			clusterID:    "test",
			expectedName: "test-bootstrap-ignition",
		},
		{
			name:         "Long cluster ID",
			clusterID:    "very-long-cluster-name-with-lots-of-characters",
			expectedName: "very-long-cluster-name-with-lots-of-characters-bootstrap-ignition",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name := GetBootstrapStorageName(tc.clusterID)
			assert.Equal(t, tc.expectedName, name)
		})
	}
}
