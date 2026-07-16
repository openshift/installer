package domains

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToDomainListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Enabled filters the response by enabled domains.
	Enabled *bool `q:"enabled"`

	// Name filters the response by domain name.
	Name string `q:"name"`

	// Limit limits the number of projects returned per page.
	Limit int `q:"limit"`
}

// ToDomainListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDomainListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the domains to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDomainListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DomainPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAvailable enumerates the domains which are available to a specific user.
func ListAvailable(client *gophercloud.ServiceClient) pagination.Pager {
	url := listAvailableURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DomainPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single domain, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToDomainCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a domain.
type CreateOpts struct {
	// Name is the name of the new domain.
	Name string `json:"name" required:"true"`

	// Description is a description of the domain.
	Description string `json:"description,omitempty"`

	// Enabled sets the domain status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ToDomainCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToDomainCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "domain")
}

// Create creates a new Domain.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a domain.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, domainID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, domainID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]any, error)
}

// UpdateOpts represents parameters to update a domain.
type UpdateOpts struct {
	// Name is the name of the domain.
	Name string `json:"name,omitempty"`

	// Description is the description of the domain.
	Description *string `json:"description,omitempty"`

	// Enabled sets the domain status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ToUpdateCreateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToDomainUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "domain")
}

// Update modifies the attributes of a domain.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDomainUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
