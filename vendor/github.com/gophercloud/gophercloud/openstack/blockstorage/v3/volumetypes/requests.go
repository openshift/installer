package volumetypes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVolumeTypeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume Type. This object is passed to
// the volumetypes.Create function. For more information about these parameters,
// see the Volume Type object.
type CreateOpts struct {
	// The name of the volume type
	Name string `json:"name" required:"true"`
	// The volume type description
	Description string `json:"description,omitempty"`
	// the ID of the existing volume snapshot
	IsPublic *bool `json:"os-volume-type-access:is_public,omitempty"`
	// Extra spec key-value pairs defined by the user.
	ExtraSpecs map[string]string `json:"extra_specs,omitempty"`
}

// ToVolumeTypeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToVolumeTypeCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "volume_type")
}

// Create will create a new Volume Type based on the values in CreateOpts. To extract
// the Volume Type object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVolumeTypeCreateMap()
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

// Delete will delete the existing Volume Type with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves the Volume Type with the provided ID. To extract the Volume Type object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeTypeListQuery() (string, error)
}

// ListOpts holds options for listing Volume Types. It is passed to the volumetypes.List
// function.
type ListOpts struct {
	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`
	// Requests a page size of items.
	Limit int `q:"limit"`
	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`
	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ToVolumeTypeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeTypeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns Volume types.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)

	if opts != nil {
		query, err := opts.ToVolumeTypeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VolumeTypePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVolumeTypeUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Volume Type. This object is passed
// to the volumetypes.Update function. For more information about the parameters, see
// the Volume Type object.
type UpdateOpts struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

// ToVolumeTypeUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToVolumeTypeUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "volume_type")
}

// Update will update the Volume Type with provided information. To extract the updated
// Volume Type from the response, call the Extract method on the UpdateResult.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumeTypeUpdateMap()
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

// ListExtraSpecs requests all the extra-specs for the given volume type ID.
func ListExtraSpecs(client *gophercloud.ServiceClient, volumeTypeID string) (r ListExtraSpecsResult) {
	resp, err := client.Get(extraSpecsListURL(client, volumeTypeID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetExtraSpec requests an extra-spec specified by key for the given volume type ID
func GetExtraSpec(client *gophercloud.ServiceClient, volumeTypeID string, key string) (r GetExtraSpecResult) {
	resp, err := client.Get(extraSpecsGetURL(client, volumeTypeID, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateExtraSpecsOptsBuilder allows extensions to add additional parameters to the
// CreateExtraSpecs requests.
type CreateExtraSpecsOptsBuilder interface {
	ToVolumeTypeExtraSpecsCreateMap() (map[string]interface{}, error)
}

// ExtraSpecsOpts is a map that contains key-value pairs.
type ExtraSpecsOpts map[string]string

// ToVolumeTypeExtraSpecsCreateMap assembles a body for a Create request based on
// the contents of ExtraSpecsOpts.
func (opts ExtraSpecsOpts) ToVolumeTypeExtraSpecsCreateMap() (map[string]interface{}, error) {
	return map[string]interface{}{"extra_specs": opts}, nil
}

// CreateExtraSpecs will create or update the extra-specs key-value pairs for
// the specified volume type.
func CreateExtraSpecs(client *gophercloud.ServiceClient, volumeTypeID string, opts CreateExtraSpecsOptsBuilder) (r CreateExtraSpecsResult) {
	b, err := opts.ToVolumeTypeExtraSpecsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(extraSpecsCreateURL(client, volumeTypeID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateExtraSpecOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateExtraSpecOptsBuilder interface {
	ToVolumeTypeExtraSpecUpdateMap() (map[string]string, string, error)
}

// ToVolumeTypeExtraSpecUpdateMap assembles a body for an Update request based on
// the contents of a ExtraSpecOpts.
func (opts ExtraSpecsOpts) ToVolumeTypeExtraSpecUpdateMap() (map[string]string, string, error) {
	if len(opts) != 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "volumetypes.ExtraSpecOpts"
		err.Info = "Must have one and only one key-value pair"
		return nil, "", err
	}

	var key string
	for k := range opts {
		key = k
	}

	return opts, key, nil
}

// UpdateExtraSpec will updates the value of the specified volume type's extra spec
// for the key in opts.
func UpdateExtraSpec(client *gophercloud.ServiceClient, volumeTypeID string, opts UpdateExtraSpecOptsBuilder) (r UpdateExtraSpecResult) {
	b, key, err := opts.ToVolumeTypeExtraSpecUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(extraSpecUpdateURL(client, volumeTypeID, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteExtraSpec will delete the key-value pair with the given key for the given
// volume type ID.
func DeleteExtraSpec(client *gophercloud.ServiceClient, volumeTypeID, key string) (r DeleteExtraSpecResult) {
	resp, err := client.Delete(extraSpecDeleteURL(client, volumeTypeID, key), &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAccesses retrieves the tenants which have access to a volume type.
func ListAccesses(client *gophercloud.ServiceClient, id string) pagination.Pager {
	url := accessURL(client, id)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AccessPage{pagination.SinglePageBase(r)}
	})
}

// AddAccessOptsBuilder allows extensions to add additional parameters to the
// AddAccess requests.
type AddAccessOptsBuilder interface {
	ToVolumeTypeAddAccessMap() (map[string]interface{}, error)
}

// AddAccessOpts represents options for adding access to a volume type.
type AddAccessOpts struct {
	// Project is the project/tenant ID to grant access.
	Project string `json:"project"`
}

// ToVolumeTypeAddAccessMap constructs a request body from AddAccessOpts.
func (opts AddAccessOpts) ToVolumeTypeAddAccessMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "addProjectAccess")
}

// AddAccess grants a tenant/project access to a volume type.
func AddAccess(client *gophercloud.ServiceClient, id string, opts AddAccessOptsBuilder) (r AddAccessResult) {
	b, err := opts.ToVolumeTypeAddAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(accessActionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveAccessOptsBuilder allows extensions to add additional parameters to the
// RemoveAccess requests.
type RemoveAccessOptsBuilder interface {
	ToVolumeTypeRemoveAccessMap() (map[string]interface{}, error)
}

// RemoveAccessOpts represents options for removing access to a volume type.
type RemoveAccessOpts struct {
	// Project is the project/tenant ID to remove access.
	Project string `json:"project"`
}

// ToVolumeTypeRemoveAccessMap constructs a request body from RemoveAccessOpts.
func (opts RemoveAccessOpts) ToVolumeTypeRemoveAccessMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "removeProjectAccess")
}

// RemoveAccess removes/revokes a tenant/project access to a volume type.
func RemoveAccess(client *gophercloud.ServiceClient, id string, opts RemoveAccessOptsBuilder) (r RemoveAccessResult) {
	b, err := opts.ToVolumeTypeRemoveAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(accessActionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateEncryptionOptsBuilder allows extensions to add additional parameters to the
// Create Encryption request.
type CreateEncryptionOptsBuilder interface {
	ToEncryptionCreateMap() (map[string]interface{}, error)
}

// CreateEncryptionOpts contains options for creating an Encryption Type object.
// This object is passed to the volumetypes.CreateEncryption function.
// For more information about these parameters,see the Encryption Type object.
type CreateEncryptionOpts struct {
	// The size of the encryption key.
	KeySize int `json:"key_size"`
	// The class of that provides the encryption support.
	Provider string `json:"provider" required:"true"`
	// Notional service where encryption is performed.
	ControlLocation string `json:"control_location"`
	// The encryption algorithm or mode.
	Cipher string `json:"cipher"`
}

// ToEncryptionCreateMap assembles a request body based on the contents of a
// CreateEncryptionOpts.
func (opts CreateEncryptionOpts) ToEncryptionCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "encryption")
}

// CreateEncryption will creates an Encryption Type object based on the CreateEncryptionOpts.
// To extract the Encryption Type object from the response, call the Extract method on the
// EncryptionCreateResult.
func CreateEncryption(client *gophercloud.ServiceClient, id string, opts CreateEncryptionOptsBuilder) (r CreateEncryptionResult) {
	b, err := opts.ToEncryptionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createEncryptionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete an encryption type for an existing Volume Type with the provided ID.
func DeleteEncryption(client *gophercloud.ServiceClient, id, encryptionID string) (r DeleteEncryptionResult) {
	resp, err := client.Delete(deleteEncryptionURL(client, id, encryptionID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetEncryption retrieves the encryption type for an existing VolumeType with the provided ID.
func GetEncryption(client *gophercloud.ServiceClient, id string) (r GetEncryptionResult) {
	resp, err := client.Get(getEncryptionURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetEncryptionSpecs retrieves the encryption type specs for an existing VolumeType with the provided ID.
func GetEncryptionSpec(client *gophercloud.ServiceClient, id, key string) (r GetEncryptionSpecResult) {
	resp, err := client.Get(getEncryptionSpecURL(client, id, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateEncryptionOptsBuilder allows extensions to add additional parameters to the
// Update encryption request.
type UpdateEncryptionOptsBuilder interface {
	ToUpdateEncryptionMap() (map[string]interface{}, error)
}

// Update Encryption Opts contains options for creating an Update Encryption Type. This object is passed to
// the volumetypes.UpdateEncryption function. For more information about these parameters,
// see the Update Encryption Type object.
type UpdateEncryptionOpts struct {
	// The size of the encryption key.
	KeySize int `json:"key_size"`
	// The class of that provides the encryption support.
	Provider string `json:"provider"`
	// Notional service where encryption is performed.
	ControlLocation string `json:"control_location"`
	// The encryption algorithm or mode.
	Cipher string `json:"cipher"`
}

// ToEncryptionCreateMap assembles a request body based on the contents of a
// UpdateEncryptionOpts.
func (opts UpdateEncryptionOpts) ToUpdateEncryptionMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "encryption")
}

// Update will update an existing encryption for a Volume Type based on the values in UpdateEncryptionOpts.
// To extract the UpdateEncryption Type object from the response, call the Extract method on the
// UpdateEncryptionResult.
func UpdateEncryption(client *gophercloud.ServiceClient, id, encryptionID string, opts UpdateEncryptionOptsBuilder) (r UpdateEncryptionResult) {
	b, err := opts.ToUpdateEncryptionMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(updateEncryptionURL(client, id, encryptionID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
