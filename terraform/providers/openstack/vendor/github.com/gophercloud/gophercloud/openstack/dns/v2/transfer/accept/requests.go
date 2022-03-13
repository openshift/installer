package accept

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToTransferAcceptListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
// https://developer.openstack.org/api-ref/dns/
type ListOpts struct {
	Status string `q:"status"`
}

// ToTransferAcceptListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTransferAcceptListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List implements a transfer accept List request.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToTransferAcceptListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TransferAcceptPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns information about a transfer accept, given its ID.
func Get(client *gophercloud.ServiceClient, transferAcceptID string) (r GetResult) {
	resp, err := client.Get(resourceURL(client, transferAcceptID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional attributes to the
// Create request.
type CreateOptsBuilder interface {
	ToTransferAcceptCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies the attributes used to create a transfer accept.
type CreateOpts struct {
	// Key is used as part of the zone transfer accept process.
	// This is only shown to the creator, and must be communicated out of band.
	Key string `json:"key" required:"true"`

	// ZoneTransferRequestID is ID for this zone transfer request
	ZoneTransferRequestID string `json:"zone_transfer_request_id" required:"true"`
}

// ToTransferAcceptCreateMap formats an CreateOpts structure into a request body.
func (opts CreateOpts) ToTransferAcceptCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Create implements a transfer accept create request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTransferAcceptCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusCreated, http.StatusAccepted},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
