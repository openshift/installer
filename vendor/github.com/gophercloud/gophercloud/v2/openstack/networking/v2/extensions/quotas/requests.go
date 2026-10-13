package quotas

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get returns Networking Quotas for a project.
func Get(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetDetail returns detailed Networking Quotas for a project.
func GetDetail(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (r GetDetailResult) {
	resp, err := client.Get(ctx, getDetailURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToQuotaUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update the Networking Quotas.
type UpdateOpts struct {
	// FloatingIP represents a number of floating IPs. A "-1" value means no limit.
	FloatingIP *int `json:"floatingip,omitempty"`

	// Network represents a number of networks. A "-1" value means no limit.
	Network *int `json:"network,omitempty"`

	// Port represents a number of ports. A "-1" value means no limit.
	Port *int `json:"port,omitempty"`

	// RBACPolicy represents a number of RBAC policies. A "-1" value means no limit.
	RBACPolicy *int `json:"rbac_policy,omitempty"`

	// Router represents a number of routers. A "-1" value means no limit.
	Router *int `json:"router,omitempty"`

	// SecurityGroup represents a number of security groups. A "-1" value means no limit.
	SecurityGroup *int `json:"security_group,omitempty"`

	// SecurityGroupRule represents a number of security group rules. A "-1" value means no limit.
	SecurityGroupRule *int `json:"security_group_rule,omitempty"`

	// Subnet represents a number of subnets. A "-1" value means no limit.
	Subnet *int `json:"subnet,omitempty"`

	// SubnetPool represents a number of subnet pools. A "-1" value means no limit.
	SubnetPool *int `json:"subnetpool,omitempty"`

	// Trunk represents a number of trunks. A "-1" value means no limit.
	Trunk *int `json:"trunk,omitempty"`
}

// ToQuotaUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToQuotaUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "quota")
}

// Update accepts a UpdateOpts struct and updates an existing Networking Quotas using the
// values provided.
func Update(ctx context.Context, c *gophercloud.ServiceClient, projectID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, projectID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
