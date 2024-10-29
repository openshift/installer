package flavors

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts allows to manage the output of the request.
type ListOpts struct {
	// The name of the flavor to filter by.
	Name string `q:"name"`
	// The flavor profile id to filter by.
	FlavorProfileID string `q:"flavor_profile_id"`
	// The enabled status of the flavor to filter by.
	Enabled *bool `q:"enabled"`
	// The fields that you want the server to return
	Fields []string `q:"fields"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// Flavor. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToFlavorCreateMap() (map[string]any, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name" required:"true"`

	// Human-readable description for the Flavor.
	Description string `json:"description,omitempty"`

	// The ID of the FlavorProfile which give the metadata for the creation of
	// a LoadBalancer.
	FlavorProfileId string `json:"flavor_profile_id" required:"true"`

	// If the resource is available for use. The default is True.
	Enabled bool `json:"enabled,omitempty"`
}

// ToFlavorCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFlavorCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "flavor")
}

// Create is and operation which add a new Flavor into the database.
// CreateResult will be returned.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFlavorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular Flavor based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToFlavorUpdateMap() (map[string]any, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the Flavor.
	Description string `json:"description,omitempty"`

	// If the resource is available for use.
	Enabled bool `json:"enabled,omitempty"`
}

// ToFlavorUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToFlavorUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "flavor")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update is an operation which modifies the attributes of the specified
// Flavor.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToFlavorUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular Flavor based on its
// unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
