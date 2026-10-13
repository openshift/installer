package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
