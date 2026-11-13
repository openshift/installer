package providers

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToProviderListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Provider attributes you want to see returned.
type ListOpts struct {
	Fields []string `q:"fields"`
}

// ToProviderListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToProviderListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// providers.
//
// Default policy settings return only those providers that are owned by
// the project who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToProviderListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ProviderPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListFlavorCapabilitiesOptsBuilder allows extensions to add additional
// parameters to the ListFlavorCapabilities request.
type ListFlavorCapabilitiesOptsBuilder interface {
	ToFlavorCapabilitiesListQuery() (string, error)
}

// ListFlavorCapabilitiesOpts allows the result of a ListFlavorCapabilities
// request to be filtered.
type ListFlavorCapabilitiesOpts struct {
	// Fields is a list of attributes to return for each capability.
	Fields []string `q:"fields"`
}

// ToFlavorCapabilitiesListQuery formats a ListFlavorCapabilitiesOpts into a
// query string.
func (opts ListFlavorCapabilitiesOpts) ToFlavorCapabilitiesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListFlavorCapabilities lists the flavor capabilities of a provider.
// This operation requires Octavia API version 2.6 or later.
func ListFlavorCapabilities(ctx context.Context, c *gophercloud.ServiceClient, provider string, opts ListFlavorCapabilitiesOptsBuilder) (r ListFlavorCapabilitiesResult) {
	url := flavorCapabilitiesURL(c, provider)
	if opts != nil {
		query, err := opts.ToFlavorCapabilitiesListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := c.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAvailabilityZoneCapabilitiesOptsBuilder allows extensions to add
// additional parameters to the ListAvailabilityZoneCapabilities request.
type ListAvailabilityZoneCapabilitiesOptsBuilder interface {
	ToAvailabilityZoneCapabilitiesListQuery() (string, error)
}

// ListAvailabilityZoneCapabilitiesOpts allows the result of a
// ListAvailabilityZoneCapabilities request to be filtered.
type ListAvailabilityZoneCapabilitiesOpts struct {
	// Fields is a list of attributes to return for each capability.
	Fields []string `q:"fields"`
}

// ToAvailabilityZoneCapabilitiesListQuery formats a
// ListAvailabilityZoneCapabilitiesOpts into a query string.
func (opts ListAvailabilityZoneCapabilitiesOpts) ToAvailabilityZoneCapabilitiesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListAvailabilityZoneCapabilities lists the availability zone capabilities
// of a provider. This operation requires Octavia API version 2.14 or later.
func ListAvailabilityZoneCapabilities(ctx context.Context, c *gophercloud.ServiceClient, provider string, opts ListAvailabilityZoneCapabilitiesOptsBuilder) (r ListAvailabilityZoneCapabilitiesResult) {
	url := availabilityZoneCapabilitiesURL(c, provider)
	if opts != nil {
		query, err := opts.ToAvailabilityZoneCapabilitiesListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := c.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
