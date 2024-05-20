package tenantnetworks

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of Networks.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.SinglePageBase(r)}
	})
}

// Get returns data about a previously created Network.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
