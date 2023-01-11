package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

/*
AccessType maps to OpenStack's Flavor.is_public field. Although the is_public
field is boolean, the request options are ternary, which is why AccessType is
a string. The following values are allowed:

The AccessType arguement is optional, and if it is not supplied, OpenStack
returns the PublicAccess flavors.
*/
type AccessType string

const (
	// PublicAccess returns public flavors and private flavors associated with
	// that project.
	PublicAccess AccessType = "true"

	// PrivateAccess (admin only) returns private flavors, across all projects.
	PrivateAccess AccessType = "false"

	// AllAccess (admin only) returns public and private flavors across all
	// projects.
	AllAccess AccessType = "None"
)

/*
ListOpts filters the results returned by the List() function.
For example, a flavor with a minDisk field of 10 will not be returned if you
specify MinDisk set to 20.

Typically, software will use the last ID of the previous call to List to set
the Marker for the current call.
*/
type ListOpts struct {
	// ChangesSince, if provided, instructs List to return only those things which
	// have changed since the timestamp provided.
	ChangesSince string `q:"changes-since"`

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your
	// criteria.
	MinDisk int `q:"minDisk"`
	MinRAM  int `q:"minRam"`

	// SortDir allows to select sort direction.
	// It can be "asc" or "desc" (default).
	SortDir string `q:"sort_dir"`

	// SortKey allows to sort by one of the flavors attributes.
	// Default is flavorid.
	SortKey string `q:"sort_key"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of
	// flavors.
	Limit int `q:"limit"`

	// AccessType, if provided, instructs List which set of flavors to return.
	// If IsPublic not provided, flavors for the current project are returned.
	AccessType AccessType `q:"is_public"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail instructs OpenStack to provide a list of flavors.
// You may provide criteria by which List curtails its results for easier
// processing.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type CreateOptsBuilder interface {
	ToFlavorCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters used for creating a flavor.
type CreateOpts struct {
	// Name is the name of the flavor.
	Name string `json:"name" required:"true"`

	// RAM is the memory of the flavor, measured in MB.
	RAM int `json:"ram" required:"true"`

	// VCPUs is the number of vcpus for the flavor.
	VCPUs int `json:"vcpus" required:"true"`

	// Disk the amount of root disk space, measured in GB.
	Disk *int `json:"disk" required:"true"`

	// ID is a unique ID for the flavor.
	ID string `json:"id,omitempty"`

	// Swap is the amount of swap space for the flavor, measured in MB.
	Swap *int `json:"swap,omitempty"`

	// RxTxFactor alters the network bandwidth of a flavor.
	RxTxFactor float64 `json:"rxtx_factor,omitempty"`

	// IsPublic flags a flavor as being available to all projects or not.
	IsPublic *bool `json:"os-flavor-access:is_public,omitempty"`

	// Ephemeral is the amount of ephemeral disk space, measured in GB.
	Ephemeral *int `json:"OS-FLV-EXT-DATA:ephemeral,omitempty"`

	// Description is a free form description of the flavor. Limited to
	// 65535 characters in length. Only printable characters are allowed.
	// New in version 2.55
	Description string `json:"description,omitempty"`
}

// ToFlavorCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToFlavorCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "flavor")
}

// Create requests the creation of a new flavor.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFlavorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	ToFlavorUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts specifies parameters used for updating a flavor.
type UpdateOpts struct {
	// Description is a free form description of the flavor. Limited to
	// 65535 characters in length. Only printable characters are allowed.
	// New in version 2.55
	Description string `json:"description,omitempty"`
}

// ToFlavorUpdateMap constructs a request body from UpdateOpts.
func (opts UpdateOpts) ToFlavorUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "flavor")
}

// Update requests the update of a new flavor.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFlavorUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves details of a single flavor. Use Extract to convert its
// result into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes the specified flavor ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAccesses retrieves the tenants which have access to a flavor.
func ListAccesses(client *gophercloud.ServiceClient, id string) pagination.Pager {
	url := accessURL(client, id)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AccessPage{pagination.SinglePageBase(r)}
	})
}

// AddAccessOptsBuilder allows extensions to add additional parameters to the
// AddAccess requests.
type AddAccessOptsBuilder interface {
	ToFlavorAddAccessMap() (map[string]interface{}, error)
}

// AddAccessOpts represents options for adding access to a flavor.
type AddAccessOpts struct {
	// Tenant is the project/tenant ID to grant access.
	Tenant string `json:"tenant"`
}

// ToFlavorAddAccessMap constructs a request body from AddAccessOpts.
func (opts AddAccessOpts) ToFlavorAddAccessMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "addTenantAccess")
}

// AddAccess grants a tenant/project access to a flavor.
func AddAccess(client *gophercloud.ServiceClient, id string, opts AddAccessOptsBuilder) (r AddAccessResult) {
	b, err := opts.ToFlavorAddAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(accessActionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveAccessOptsBuilder allows extensions to add additional parameters to the
// RemoveAccess requests.
type RemoveAccessOptsBuilder interface {
	ToFlavorRemoveAccessMap() (map[string]interface{}, error)
}

// RemoveAccessOpts represents options for removing access to a flavor.
type RemoveAccessOpts struct {
	// Tenant is the project/tenant ID to grant access.
	Tenant string `json:"tenant"`
}

// ToFlavorRemoveAccessMap constructs a request body from RemoveAccessOpts.
func (opts RemoveAccessOpts) ToFlavorRemoveAccessMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "removeTenantAccess")
}

// RemoveAccess removes/revokes a tenant/project access to a flavor.
func RemoveAccess(client *gophercloud.ServiceClient, id string, opts RemoveAccessOptsBuilder) (r RemoveAccessResult) {
	b, err := opts.ToFlavorRemoveAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(accessActionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ExtraSpecs requests all the extra-specs for the given flavor ID.
func ListExtraSpecs(client *gophercloud.ServiceClient, flavorID string) (r ListExtraSpecsResult) {
	resp, err := client.Get(extraSpecsListURL(client, flavorID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetExtraSpec(client *gophercloud.ServiceClient, flavorID string, key string) (r GetExtraSpecResult) {
	resp, err := client.Get(extraSpecsGetURL(client, flavorID, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateExtraSpecsOptsBuilder allows extensions to add additional parameters to the
// CreateExtraSpecs requests.
type CreateExtraSpecsOptsBuilder interface {
	ToFlavorExtraSpecsCreateMap() (map[string]interface{}, error)
}

// ExtraSpecsOpts is a map that contains key-value pairs.
type ExtraSpecsOpts map[string]string

// ToFlavorExtraSpecsCreateMap assembles a body for a Create request based on
// the contents of ExtraSpecsOpts.
func (opts ExtraSpecsOpts) ToFlavorExtraSpecsCreateMap() (map[string]interface{}, error) {
	return map[string]interface{}{"extra_specs": opts}, nil
}

// CreateExtraSpecs will create or update the extra-specs key-value pairs for
// the specified Flavor.
func CreateExtraSpecs(client *gophercloud.ServiceClient, flavorID string, opts CreateExtraSpecsOptsBuilder) (r CreateExtraSpecsResult) {
	b, err := opts.ToFlavorExtraSpecsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(extraSpecsCreateURL(client, flavorID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateExtraSpecOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateExtraSpecOptsBuilder interface {
	ToFlavorExtraSpecUpdateMap() (map[string]string, string, error)
}

// ToFlavorExtraSpecUpdateMap assembles a body for an Update request based on
// the contents of a ExtraSpecOpts.
func (opts ExtraSpecsOpts) ToFlavorExtraSpecUpdateMap() (map[string]string, string, error) {
	if len(opts) != 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "flavors.ExtraSpecOpts"
		err.Info = "Must have 1 and only one key-value pair"
		return nil, "", err
	}

	var key string
	for k := range opts {
		key = k
	}

	return opts, key, nil
}

// UpdateExtraSpec will updates the value of the specified flavor's extra spec
// for the key in opts.
func UpdateExtraSpec(client *gophercloud.ServiceClient, flavorID string, opts UpdateExtraSpecOptsBuilder) (r UpdateExtraSpecResult) {
	b, key, err := opts.ToFlavorExtraSpecUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(extraSpecUpdateURL(client, flavorID, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteExtraSpec will delete the key-value pair with the given key for the given
// flavor ID.
func DeleteExtraSpec(client *gophercloud.ServiceClient, flavorID, key string) (r DeleteExtraSpecResult) {
	resp, err := client.Delete(extraSpecDeleteURL(client, flavorID, key), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
