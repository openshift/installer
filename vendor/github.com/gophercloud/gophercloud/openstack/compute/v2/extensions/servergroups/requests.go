package servergroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of
// ServerGroups.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return ServerGroupPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies Server Group creation parameters.
type CreateOpts struct {
	// Name is the name of the server group.
	Name string `json:"name" required:"true"`

	// Policies are the server group policies.
	Policies []string `json:"policies,omitempty"`

	// Policy specifies the name of a policy.
	// Requires microversion 2.64 or later.
	Policy string `json:"policy,omitempty"`

	// Rules specifies the set of rules.
	// Requires microversion 2.64 or later.
	Rules *Rules `json:"rules,omitempty"`
}

// ToServerGroupCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToServerGroupCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "server_group")
}

// Create requests the creation of a new Server Group.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToServerGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns data about a previously created ServerGroup.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests the deletion of a previously allocated ServerGroup.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
