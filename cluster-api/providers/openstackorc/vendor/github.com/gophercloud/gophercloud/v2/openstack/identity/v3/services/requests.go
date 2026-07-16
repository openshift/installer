package services

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToServiceCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a service.
type CreateOpts struct {
	// Name is the name of the service.
	Name string `json:"name,omitempty"`

	// Description is the description of the service.
	Description string `json:"description,omitempty"`

	// Type is the type of the service.
	Type string `json:"type"`

	// Enabled is whether or not the service is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the service.
	Extra map[string]any `json:"-"`
}

// ToServiceCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToServiceCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "service")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["service"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create adds a new service of the requested type to the catalog.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToServiceCreateMap()
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

// ListOptsBuilder enables extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToServiceListMap() (string, error)
}

// ListOpts provides options for filtering the List results.
type ListOpts struct {
	// ServiceType filter the response by a type of service.
	ServiceType string `q:"type"`

	// Name filters the response by a service name.
	Name string `q:"name"`
}

// ToServiceListMap builds a list query from the list options.
func (opts ListOpts) ToServiceListMap() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the services available to a specific user.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToServiceListMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns additional information about a service, given its ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, serviceID string) (r GetResult) {
	resp, err := client.Get(ctx, serviceURL(client, serviceID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToServiceUpdateMap() (map[string]any, error)
}

// UpdateOpts provides options for updating a service.
type UpdateOpts struct {
	// Name is an updated name for the service.
	Name *string `json:"name,omitempty"`

	// Description is an update description for the service.
	Description *string `json:"description,omitempty"`

	// Type is the type of the service.
	Type string `json:"type,omitempty"`

	// Enabled is whether or not the service is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the service.
	Extra map[string]any `json:"-"`
}

// ToServiceUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToServiceUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "service")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["service"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Service.
func Update(ctx context.Context, client *gophercloud.ServiceClient, serviceID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToServiceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, serviceID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes an existing service.
// It either deletes all associated endpoints, or fails until all endpoints
// are deleted.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, serviceID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, serviceURL(client, serviceID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
