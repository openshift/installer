package powervc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
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

func TestServerTimeoutError(t *testing.T) {
	pending := []servers.Server{
		{ID: "1", Name: "infra-abcde-master-0"},
	}
	failed := []servers.Server{
		{ID: "2", Name: "infra-abcde-master-1", Fault: servers.Fault{Message: "Not authorized."}},
	}

	tests := []struct {
		name          string
		pending       []servers.Server
		failed        []servers.Server
		lastServerErr error
		expected      string
	}{
		{
			name:     "pending only",
			pending:  pending,
			expected: `timed out waiting for servers to be deleted (pending: "infra-abcde-master-0" (1))`,
		},
		{
			name:     "pending and failed",
			pending:  pending,
			failed:   failed,
			expected: `timed out waiting for servers to be deleted (pending: "infra-abcde-master-0" (1); 1 server(s) in ERROR state: "infra-abcde-master-1" (2): Not authorized. (PowerVC does not let non-admin users delete volume-backed servers; delete them as an admin user and run destroy again))`,
		},
		{
			name:          "with lastServerErr",
			pending:       pending,
			failed:        failed,
			lastServerErr: errors.New("network unreachable"),
			expected:      `timed out waiting for servers to be deleted (pending: "infra-abcde-master-0" (1); 1 server(s) in ERROR state: "infra-abcde-master-1" (2): Not authorized. (PowerVC does not let non-admin users delete volume-backed servers; delete them as an admin user and run destroy again)): network unreachable`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := serverTimeoutError(tt.pending, tt.failed, tt.lastServerErr)
			assert.EqualError(t, err, tt.expected)
		})
	}
}

func TestContainerMatchesFilter(t *testing.T) {
	filter := map[string]string{"openshiftClusterID": "infra-abcde"}
	tests := []struct {
		name     string
		metadata map[string]string
		filter   map[string]string
		expected bool
	}{
		{
			name:     "key case changed by Swift",
			metadata: map[string]string{"Openshiftclusterid": "infra-abcde"},
			filter:   filter,
			expected: true,
		},
		{
			name:     "exact key",
			metadata: map[string]string{"openshiftClusterID": "infra-abcde", "Other": "x"},
			filter:   filter,
			expected: true,
		},
		{
			name:     "other cluster",
			metadata: map[string]string{"Openshiftclusterid": "infra-zzzzz"},
			filter:   filter,
		},
		{
			name:   "no metadata",
			filter: filter,
		},
		{
			name:     "empty filter matches nothing",
			metadata: map[string]string{"Openshiftclusterid": "infra-abcde"},
			filter:   map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, containerMatchesFilter(tt.metadata, tt.filter))
		})
	}
}

func TestObjectURL(t *testing.T) {
	conn := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "https://swift.example.com/v1/AUTH_project/",
	}
	tests := []struct {
		name      string
		container string
		object    string
		expected  string
	}{
		{
			name:      "slashes stay path separators",
			container: "infra-abcde-image-registry-x",
			object:    "docker/registry/v2/blobs/data",
			expected:  "https://swift.example.com/v1/AUTH_project/infra-abcde-image-registry-x/docker/registry/v2/blobs/data",
		},
		{
			name:      "segments are escaped",
			container: "c",
			object:    "a b/c?d",
			expected:  "https://swift.example.com/v1/AUTH_project/c/a%20b/c%3Fd",
		},
		{
			name:      "no slash",
			container: "infra-abcde-ignition",
			object:    "infra-abcde-ignition",
			expected:  "https://swift.example.com/v1/AUTH_project/infra-abcde-ignition/infra-abcde-ignition",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, objectURL(conn, tt.container, tt.object))
		})
	}
}

func TestIgnitionContainerName(t *testing.T) {
	assert.Equal(t, "infra-abcde-ignition", ignitionContainerName("infra-abcde"))
	assert.Equal(t, "", ignitionContainerName(""))
}

func TestDescribeRetryError(t *testing.T) {
	respErr := gophercloud.ErrUnexpectedResponseCode{
		URL:    "https://powervc.example.com:8774/v2.1/servers/detail",
		Method: "GET",
		Actual: 500,
		Body:   []byte("secret-body"),
	}

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "unexpected response",
			err:      respErr,
			expected: "unexpected HTTP status 500",
		},
		{
			name:     "wrapped unexpected response",
			err:      fmt.Errorf("wrapped: %w", respErr),
			expected: "unexpected HTTP status 500",
		},
		// gophercloud's re-authentication errors embed the response error
		// without implementing Unwrap, so errors.As does not reach it.
		{
			name:     "error after re-authentication",
			err:      &gophercloud.ErrErrorAfterReauthentication{ErrOriginal: respErr},
			expected: "after re-authenticating: unexpected HTTP status 500",
		},
		{
			name:     "failed re-authentication",
			err:      &gophercloud.ErrUnableToReauthenticate{ErrOriginal: respErr, ErrReauth: respErr},
			expected: "unable to re-authenticate",
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: "connection refused",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := describeRetryError(tt.err)
			assert.Equal(t, tt.expected, got)
			assert.NotContains(t, got, respErr.URL)
			assert.NotContains(t, got, "secret-body")
		})
	}
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
		{ID: "6", Name: "infra-abcde-master-1-0", DeviceID: "attached-server-id"},
		{ID: "7", Name: "infra-abcde-worker-0-xyz12-0"},
		{ID: "8", Name: "infra-abcde-bastion-0"},
		{ID: "9", Name: "infra-abcde-bastion", DeviceID: "bastion-server-id"},
		{ID: "10", Name: "infra-abcde-port"},
		{ID: "11", Name: "infra-abcde-masterless-0"},
	}

	tests := []struct {
		name     string
		infraID  string
		expected []string
	}{
		{
			name:     "matches only unattached master, worker and bootstrap ports of the cluster",
			infraID:  "infra-abcde",
			expected: []string{"1", "2", "7"},
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
