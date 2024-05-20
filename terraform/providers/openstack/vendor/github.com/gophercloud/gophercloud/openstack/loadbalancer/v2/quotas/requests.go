package quotas

import (
	"github.com/gophercloud/gophercloud"
)

// Get returns load balancer Quotas for a project.
func Get(client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	resp, err := client.Get(getURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToQuotaUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update the load balancer Quotas.
type UpdateOpts struct {
	// Loadbalancer represents the number of load balancers. A "-1" value means no limit.
	Loadbalancer *int `json:"loadbalancer,omitempty"`

	// Listener represents the number of listeners. A "-1" value means no limit.
	Listener *int `json:"listener,omitempty"`

	// Member represents the number of members. A "-1" value means no limit.
	Member *int `json:"member,omitempty"`

	// Poool represents the number of pools. A "-1" value means no limit.
	Pool *int `json:"pool,omitempty"`

	// HealthMonitor represents the number of healthmonitors. A "-1" value means no limit.
	Healthmonitor *int `json:"healthmonitor,omitempty"`

	// L7Policy represents the number of l7policies. A "-1" value means no limit.
	L7Policy *int `json:"l7policy,omitempty"`

	// L7Rule represents the number of l7rules. A "-1" value means no limit.
	L7Rule *int `json:"l7rule,omitempty"`
}

// ToQuotaUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToQuotaUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "quota")
}

// Update accepts a UpdateOpts struct and updates an existing load balancer Quotas using the
// values provided.
func Update(c *gophercloud.ServiceClient, projectID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(updateURL(c, projectID), b, &r.Body, &gophercloud.RequestOpts{
		// allow 200 (neutron/lbaasv2) and 202 (octavia)
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
