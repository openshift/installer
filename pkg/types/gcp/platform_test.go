package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
