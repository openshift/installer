package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCloudEnvironment(t *testing.T) {
	cases := []struct {
		name      string
		projectID string
		region    string
		expected  string
	}{
		{
			name:      "standard project",
			projectID: "my-project",
			region:    "us-central1",
			expected:  "",
		},
		{
			name:      "eu0 sovereign cloud",
			projectID: "eu0:my-project",
			region:    "u-de-1",
			expected:  CloudEnvironmentSovereign,
		},
		{
			name:      "s3ns sovereign cloud",
			projectID: "s3ns:my-project",
			region:    "u-fr-1",
			expected:  CloudEnvironmentSovereign,
		},
		{
			name:      "domain-scoped project without sovereign region is not sovereign",
			projectID: "other-prefix:my-project",
			region:    "us-central1",
			expected:  "",
		},
		{
			name:      "sovereign region without domain-scoped project is not sovereign",
			projectID: "my-project",
			region:    "u-de-1",
			expected:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetCloudEnvironment(tc.projectID, tc.region)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormatKMSKeyResourcePath(t *testing.T) {
	cases := []struct {
		name         string
		kmsKey       *KMSKeyReference
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
			kmsKey: &KMSKeyReference{
				Name:     "bootstrap-key",
				KeyRing:  "bootstrap-keyring",
				Location: "us-central1",
			},
			projectID:    "default-project",
			expectedPath: "projects/default-project/locations/us-central1/keyRings/bootstrap-keyring/cryptoKeys/bootstrap-key",
		},
		{
			name: "KMS key with project ID overrides default project",
			kmsKey: &KMSKeyReference{
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
			kmsKey: &KMSKeyReference{
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
			path := FormatKMSKeyResourcePath(tc.kmsKey, tc.projectID)
			assert.Equal(t, tc.expectedPath, path)
		})
	}
}

func TestGetDefaultServiceAccount(t *testing.T) {
	cases := []struct {
		name      string
		projectID string
		clusterID string
		role      string
		expected  string
	}{
		{
			name:      "standard GCP project",
			projectID: "my-project",
			clusterID: "012345678",
			role:      "master",
			expected:  "012345678-m@my-project.iam.gserviceaccount.com",
		},
		{
			name:      "domain-scoped project",
			projectID: "eu0:my-project",
			clusterID: "012345678",
			role:      "master",
			expected:  "012345678-m@my-project.eu0.iam.gserviceaccount.com",
		},
		{
			name:      "domain-scoped project worker role",
			projectID: "eu0:my-project",
			clusterID: "012345678",
			role:      "worker",
			expected:  "012345678-w@my-project.eu0.iam.gserviceaccount.com",
		},
		{
			name:      "domain-scoped project with different prefix",
			projectID: "other-prefix:my-project",
			clusterID: "012345678",
			role:      "master",
			expected:  "012345678-m@my-project.other-prefix.iam.gserviceaccount.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &Platform{ProjectID: tc.projectID}
			result := GetDefaultServiceAccount(platform, tc.clusterID, tc.role)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetConfiguredOSImage(t *testing.T) {
	poolImage := &OSImage{Name: "pool-image", Project: "pool-project"}
	defaultImage := &OSImage{Name: "default-image", Project: "default-project"}

	cases := []struct {
		name                   string
		defaultMachinePlatform *MachinePool
		mpool                  *MachinePool
		expected               *OSImage
	}{
		{
			name:     "no image anywhere",
			mpool:    &MachinePool{},
			expected: nil,
		},
		{
			name:     "nil pool and no default machine platform",
			mpool:    nil,
			expected: nil,
		},
		{
			name:     "image on pool",
			mpool:    &MachinePool{OSImage: poolImage},
			expected: poolImage,
		},
		{
			name:                   "image only on default machine platform",
			defaultMachinePlatform: &MachinePool{OSImage: defaultImage},
			mpool:                  &MachinePool{},
			expected:               defaultImage,
		},
		{
			name:                   "nil pool falls back to default machine platform",
			defaultMachinePlatform: &MachinePool{OSImage: defaultImage},
			mpool:                  nil,
			expected:               defaultImage,
		},
		{
			name:                   "pool image overrides default machine platform",
			defaultMachinePlatform: &MachinePool{OSImage: defaultImage},
			mpool:                  &MachinePool{OSImage: poolImage},
			expected:               poolImage,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &Platform{DefaultMachinePlatform: tc.defaultMachinePlatform}
			assert.Equal(t, tc.expected, GetConfiguredOSImage(platform, tc.mpool))
		})
	}
}

func TestNeedsRHCOSUpload(t *testing.T) {
	osImage := &OSImage{Name: "my-image", Project: "my-project"}

	cases := []struct {
		name                   string
		projectID              string
		region                 string
		defaultMachinePlatform *MachinePool
		mpool                  *MachinePool
		expected               bool
	}{
		{
			name:      "public gcp without image uses the public rhcos project",
			projectID: "my-project",
			region:    "us-central1",
			mpool:     &MachinePool{},
			expected:  false,
		},
		{
			name:      "sovereign without image requires an upload",
			projectID: "eu0:my-project",
			region:    "u-northeast1",
			mpool:     &MachinePool{},
			expected:  true,
		},
		{
			name:      "sovereign with image on pool uses the custom image",
			projectID: "eu0:my-project",
			region:    "u-northeast1",
			mpool:     &MachinePool{OSImage: osImage},
			expected:  false,
		},
		{
			// A pool that inherits its image from the default machine platform
			// must not trigger an upload it will never reference.
			name:                   "sovereign with image only on default machine platform",
			projectID:              "eu0:my-project",
			region:                 "u-northeast1",
			defaultMachinePlatform: &MachinePool{OSImage: osImage},
			mpool:                  &MachinePool{},
			expected:               false,
		},
		{
			// Compute pools may omit the platform stanza entirely.
			name:      "sovereign with nil pool requires an upload",
			projectID: "eu0:my-project",
			region:    "u-northeast1",
			mpool:     nil,
			expected:  true,
		},
		{
			name:                   "sovereign with nil pool inheriting a default image",
			projectID:              "eu0:my-project",
			region:                 "u-northeast1",
			defaultMachinePlatform: &MachinePool{OSImage: osImage},
			mpool:                  nil,
			expected:               false,
		},
		{
			name:      "domain-scoped project outside a sovereign region",
			projectID: "eu0:my-project",
			region:    "us-central1",
			mpool:     &MachinePool{},
			expected:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &Platform{
				ProjectID:              tc.projectID,
				Region:                 tc.region,
				DefaultMachinePlatform: tc.defaultMachinePlatform,
			}
			assert.Equal(t, tc.expected, NeedsRHCOSUpload(platform, tc.mpool))
		})
	}
}
