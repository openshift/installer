package powervs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const errorBody = `{"errors":[{"code":"err","message":"err"}]}`

// writeBody writes body to w and fails the test if the write fails.
func writeBody(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	_, err := w.Write([]byte(body)) //nolint:gosec // test fixture served by httptest, not a browser response
	assert.NoError(t, err)
}

func TestDeleteRoutingTable(t *testing.T) {
	const (
		vpcID = "vpc-id"
		rtID  = "rt-id"
	)

	cases := []struct {
		name             string
		getStatus        int
		deleteStatus     int
		expectErr        bool
		expectPending    bool
		expectDeleteCall bool
	}{
		{
			name:             "GET 404 means already deleted",
			getStatus:        http.StatusNotFound,
			expectDeleteCall: false,
		},
		{
			name:          "GET 500 is returned and item stays pending",
			getStatus:     http.StatusInternalServerError,
			expectErr:     true,
			expectPending: true,
		},
		{
			name:          "GET 429 is returned and item stays pending",
			getStatus:     http.StatusTooManyRequests,
			expectErr:     true,
			expectPending: true,
		},
		{
			name:             "GET ok and DELETE ok",
			getStatus:        http.StatusOK,
			deleteStatus:     http.StatusNoContent,
			expectDeleteCall: true,
		},
		{
			name:             "GET ok and DELETE fails",
			getStatus:        http.StatusOK,
			deleteStatus:     http.StatusConflict,
			expectErr:        true,
			expectPending:    true,
			expectDeleteCall: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deleteCalled := false
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					w.WriteHeader(tc.getStatus)
					if tc.getStatus == http.StatusOK {
						writeBody(t, w, `{"id":"`+rtID+`","name":"rt"}`)
					} else {
						writeBody(t, w, errorBody)
					}
				case http.MethodDelete:
					deleteCalled = true
					w.WriteHeader(tc.deleteStatus)
					if tc.deleteStatus >= 400 {
						writeBody(t, w, errorBody)
					}
				}
			}))
			defer srv.Close()

			vpcSvc, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
				URL:           srv.URL,
				Authenticator: &core.NoAuthAuthenticator{},
			})
			require.NoError(t, err)

			item := cloudResource{key: "rt", name: "rt", status: vpcID, typeName: routingTableTypeName, id: rtID}
			o := &ClusterUninstaller{
				Logger:             logrus.New(),
				vpcSvc:             vpcSvc,
				pendingItemTracker: newPendingItemTracker(),
			}
			o.insertPendingItems(routingTableTypeName, []cloudResource{item})

			err = o.deleteRoutingTable(item)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectDeleteCall, deleteCalled)
			if tc.expectPending {
				assert.Len(t, o.getPendingItems(routingTableTypeName), 1)
			} else {
				assert.Empty(t, o.getPendingItems(routingTableTypeName))
			}
		})
	}
}

func TestListRoutingTables(t *testing.T) {
	const infraID = "test-abcde"

	const vpcList = `{"limit":50,"total_count":2,"vpcs":[
		{"id":"vpc-cluster","name":"vpc-test-abcde"},
		{"id":"vpc-user","name":"my-vpc"}]}`

	// Routing tables per VPC. Every VPC also has a table the installer did not create.
	tables := map[string]string{
		"vpc-cluster": `{"limit":20,"routing_tables":[{"id":"rt-c1","name":"test-abcde-c-rt"},{"id":"rt-c2","name":"default-rt"}]}`,
		"vpc-user":    `{"limit":20,"routing_tables":[{"id":"rt-u1","name":"test-abcde-u-rt"},{"id":"rt-u2","name":"user-rt"}]}`,
	}

	cases := []struct {
		name      string
		userVPC   string
		expectIDs []string
	}{
		{
			name:      "no user VPC: only the cluster VPC is searched",
			expectIDs: []string{"rt-c1"},
		},
		{
			name:      "user VPC by name: its installer-created table is found, the user's own is not",
			userVPC:   "my-vpc",
			expectIDs: []string{"rt-c1", "rt-u1"},
		},
		{
			name:      "user VPC by ID",
			userVPC:   "vpc-user",
			expectIDs: []string{"rt-c1", "rt-u1"},
		},
		{
			name:      "user VPC that no longer exists is not an error",
			userVPC:   "gone-vpc",
			expectIDs: []string{"rt-c1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				path := r.URL.Path
				switch {
				case path == "/vpcs":
					writeBody(t, w, vpcList)
				case strings.HasSuffix(path, "/routing_tables"):
					vpcID := strings.Split(path, "/")[2]
					writeBody(t, w, tables[vpcID])
				case strings.Contains(path, "/routing_tables/"):
					parts := strings.Split(path, "/")
					id := parts[len(parts)-1]
					name := infraID + "-c-rt"
					if id == "rt-u1" {
						name = infraID + "-u-rt"
					}
					writeBody(t, w, `{"id":"`+id+`","name":"`+name+`"}`)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer srv.Close()

			vpcSvc, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
				URL:           srv.URL,
				Authenticator: &core.NoAuthAuthenticator{},
			})
			require.NoError(t, err)

			o := &ClusterUninstaller{
				Logger:             logrus.New(),
				InfraID:            infraID,
				VPC:                tc.userVPC,
				vpcSvc:             vpcSvc,
				pendingItemTracker: newPendingItemTracker(),
			}

			got, err := o.listRoutingTables()
			require.NoError(t, err)

			var ids []string
			for _, rt := range got.list() {
				ids = append(ids, rt.id)
			}
			assert.ElementsMatch(t, tc.expectIDs, ids)
		})
	}
}
