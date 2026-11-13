package powervs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

// newTestVPCClient returns a Client whose VPC API talks to a fake server. The fake
// serves the region lookup (pointing the endpoint back at itself) plus whatever
// handlers the test registers.
func newTestVPCClient(t *testing.T, handlers func(baseURL string) map[string]http.HandlerFunc) *Client {
	t.Helper()

	mux := http.NewServeMux()
	var server *httptest.Server
	mux.HandleFunc("/v1/regions/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{"name":"us-south","endpoint":%q}`, server.URL)
		assert.NoError(t, err)
	})
	server = httptest.NewServer(mux)
	t.Cleanup(server.Close)
	for path, h := range handlers(server.URL) {
		mux.HandleFunc(path, h)
	}

	vpcAPI, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		URL:           server.URL + "/v1",
		Authenticator: &core.NoAuthAuthenticator{},
	})
	assert.NoError(t, err)
	return &Client{vpcAPI: vpcAPI}
}

func writeJSON(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(body))
	assert.NoError(t, err)
}

func TestListPublicGatewaysByVPC(t *testing.T) {
	c := newTestVPCClient(t, func(baseURL string) map[string]http.HandlerFunc {
		return map[string]http.HandlerFunc{
			"/v1/public_gateways": func(w http.ResponseWriter, r *http.Request) {
				// Two pages: the second is requested with start=page2.
				if r.URL.Query().Get("start") == "page2" {
					writeJSON(t, w, `{"limit":2,"total_count":3,"public_gateways":[
					{"id":"pgw-3","name":"c","vpc":{"id":"vpc-a"},"zone":{"name":"us-south-3"}}]}`)
					return
				}
				writeJSON(t, w, fmt.Sprintf(`{"limit":2,"total_count":3,
				"next":{"href":"%s/v1/public_gateways?start=page2"},
				"public_gateways":[
					{"id":"pgw-1","name":"a","vpc":{"id":"vpc-a"},"zone":{"name":"us-south-1"}},
					{"id":"pgw-2","name":"b","vpc":{"id":"vpc-other"},"zone":{"name":"us-south-2"}}]}`, baseURL))
			},
		}
	})

	got, err := c.ListPublicGatewaysByVPC(context.Background(), "vpc-a", "us-south")
	assert.NoError(t, err)

	var ids []string
	for _, gw := range got {
		ids = append(ids, *gw.ID)
	}
	// Both pages are read, and the gateway of the other VPC is filtered out.
	assert.Equal(t, []string{"pgw-1", "pgw-3"}, ids)
}

func TestListPublicGatewaysByVPCError(t *testing.T) {
	c := newTestVPCClient(t, func(string) map[string]http.HandlerFunc {
		return map[string]http.HandlerFunc{
			"/v1/public_gateways": func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		}
	})

	_, err := c.ListPublicGatewaysByVPC(context.Background(), "vpc-a", "us-south")
	assert.Error(t, err)
}

func TestListVPCRoutingTables(t *testing.T) {
	c := newTestVPCClient(t, func(baseURL string) map[string]http.HandlerFunc {
		return map[string]http.HandlerFunc{
			"/v1/vpcs/vpc-a/routing_tables": func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("start") == "page2" {
					writeJSON(t, w, `{"limit":1,"total_count":2,"routing_tables":[
					{"id":"rt-2","name":"two","route_transit_gateway_ingress":true}]}`)
					return
				}
				writeJSON(t, w, fmt.Sprintf(`{"limit":1,"total_count":2,
				"next":{"href":"%s/v1/vpcs/vpc-a/routing_tables?start=page2"},
				"routing_tables":[{"id":"rt-1","name":"one","route_transit_gateway_ingress":false}]}`, baseURL))
			},
		}
	})

	got, err := c.ListVPCRoutingTables(context.Background(), "vpc-a", "us-south")
	assert.NoError(t, err)

	if assert.Len(t, got, 2) {
		assert.Equal(t, "rt-1", *got[0].ID)
		assert.False(t, *got[0].RouteTransitGatewayIngress)
		assert.Equal(t, "rt-2", *got[1].ID)
		assert.True(t, *got[1].RouteTransitGatewayIngress)
	}
}

func TestSetVPCRoutingTableRouteNextHop(t *testing.T) {
	for _, tc := range []struct {
		name        string
		currentHop  string
		routeName   string
		secondPage  bool
		listErr     bool
		patchErr    bool
		wantChanged bool
		wantPatch   bool
		wantErr     bool
	}{
		{name: "next hop differs", currentHop: "10.0.64.4", routeName: "rt-default", wantChanged: true, wantPatch: true},
		{name: "next hop already correct", currentHop: "10.0.0.5", routeName: "rt-default"},
		{name: "route missing", currentHop: "10.0.0.5", routeName: "other", wantErr: true},
		{name: "route on second page", currentHop: "10.0.64.4", routeName: "rt-default", secondPage: true, wantChanged: true, wantPatch: true},
		{name: "route list error", currentHop: "10.0.64.4", routeName: "rt-default", listErr: true, wantErr: true},
		{name: "patch error", currentHop: "10.0.64.4", routeName: "rt-default", patchErr: true, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			patched := false
			c := newTestVPCClient(t, func(baseURL string) map[string]http.HandlerFunc {
				return map[string]http.HandlerFunc{
					"/v1/vpcs/vpc-a/routing_tables/rt-1/routes": func(w http.ResponseWriter, r *http.Request) {
						if tc.listErr {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						if tc.secondPage {
							if r.URL.Query().Get("start") == "page2" {
								writeJSON(t, w, fmt.Sprintf(`{"limit":1,"total_count":2,"routes":[
								{"id":"route-1","name":"rt-default","action":"deliver","destination":"0.0.0.0/0","next_hop":{"address":%q}}]}`, tc.currentHop))
								return
							}
							writeJSON(t, w, fmt.Sprintf(`{"limit":1,"total_count":2,
							"next":{"href":"%s/v1/vpcs/vpc-a/routing_tables/rt-1/routes?start=page2"},
							"routes":[{"id":"route-0","name":"rt-other","action":"deliver","destination":"1.1.1.0/24","next_hop":{"address":"10.0.64.1"}}]}`, baseURL))
							return
						}
						writeJSON(t, w, fmt.Sprintf(`{"limit":50,"total_count":1,"routes":[
						{"id":"route-1","name":"rt-default","action":"deliver","destination":"0.0.0.0/0","next_hop":{"address":%q}}]}`, tc.currentHop))
					},
					"/v1/vpcs/vpc-a/routing_tables/rt-1/routes/route-1": func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, http.MethodPatch, r.Method)
						if tc.patchErr {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						var body struct {
							NextHop struct {
								Address string `json:"address"`
							} `json:"next_hop"`
						}
						assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
						assert.Equal(t, "10.0.0.5", body.NextHop.Address)
						patched = true
						writeJSON(t, w, `{"id":"route-1","name":"rt-default"}`)
					},
				}
			})

			changed, err := c.SetVPCRoutingTableRouteNextHop(context.Background(), "vpc-a", "us-south", "rt-1", tc.routeName, "10.0.0.5")
			if tc.wantErr {
				assert.Error(t, err)
				assert.False(t, changed)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantChanged, changed)
			assert.Equal(t, tc.wantPatch, patched)
		})
	}
}
