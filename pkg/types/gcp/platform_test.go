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
