package powervc

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/stretchr/testify/assert"
)

func TestClusterPorts(t *testing.T) {
	allPorts := []ports.Port{
		{ID: "1", Name: "infra-abcde-master-0-0"},
		{ID: "2", Name: "infra-abcde-bootstrap-0"},
		{ID: "3", Name: "infra-abcdef-master-0-0"},
		{ID: "4", Name: "other-cluster-master-0-0"},
		{ID: "5", Name: ""},
	}

	tests := []struct {
		name     string
		infraID  string
		expected []string
	}{
		{
			name:     "matches only ports with the infraID prefix",
			infraID:  "infra-abcde",
			expected: []string{"1", "2"},
		},
		{
			name:     "empty infraID matches nothing",
			infraID:  "",
			expected: nil,
		},
		{
			name:     "no matching ports",
			infraID:  "missing",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ids []string
			for _, port := range clusterPorts(allPorts, tt.infraID) {
				ids = append(ids, port.ID)
			}
			assert.Equal(t, tt.expected, ids)
		})
	}
}
