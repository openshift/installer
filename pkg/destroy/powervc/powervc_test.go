package powervc

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/stretchr/testify/assert"
)

func TestRemainingServers(t *testing.T) {
	filter := map[string]string{"openshiftClusterID": "infra-abcde"}
	allServers := []servers.Server{
		{ID: "1", Status: "ACTIVE", Metadata: map[string]string{"openshiftClusterID": "infra-abcde"}},
		{ID: "2", Status: "ERROR", Metadata: map[string]string{"openshiftClusterID": "infra-abcde"}},
		{ID: "3", Status: "DELETED", Metadata: map[string]string{"openshiftClusterID": "infra-abcde", "Name": "x"}},
		{ID: "4", Status: "ERROR", Metadata: map[string]string{"openshiftClusterID": "other"}},
		{ID: "5", Status: "ACTIVE"},
		{ID: "6", Status: "ERROR", TaskState: "deleting", Metadata: map[string]string{"openshiftClusterID": "infra-abcde"}},
	}

	tests := []struct {
		name            string
		servers         []servers.Server
		filter          map[string]string
		expectedPending []string
		expectedFailed  []string
	}{
		{
			name:            "splits the cluster's servers into pending and failed",
			servers:         allServers,
			filter:          filter,
			expectedPending: []string{"1", "3", "6"},
			expectedFailed:  []string{"2"},
		},
		{
			name:    "no servers left",
			servers: nil,
			filter:  filter,
		},
		{
			name:    "empty filter matches nothing",
			servers: allServers,
			filter:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending, failed := remainingServers(tt.servers, tt.filter)
			assert.Equal(t, tt.expectedPending, serverIDs(pending))
			assert.Equal(t, tt.expectedFailed, serverIDs(failed))
		})
	}
}

func TestDescribeServers(t *testing.T) {
	list := []servers.Server{
		{ID: "1", Name: "infra-abcde-master-0", Fault: servers.Fault{Message: "Not authorized."}},
		{ID: "2", Name: "infra-abcde-master-1"},
	}
	assert.Equal(t, `"infra-abcde-master-0" (1): Not authorized., "infra-abcde-master-1" (2)`, describeServers(list, true))
	assert.Equal(t, `"infra-abcde-master-0" (1), "infra-abcde-master-1" (2)`, describeServers(list, false))
}

func serverIDs(list []servers.Server) []string {
	var ids []string
	for _, server := range list {
		ids = append(ids, server.ID)
	}
	return ids
}

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
