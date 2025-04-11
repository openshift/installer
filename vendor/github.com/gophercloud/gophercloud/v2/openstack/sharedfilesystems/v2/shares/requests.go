package shares

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// SchedulerHints contains options for providing scheduler hints when creating
// a Share.
type SchedulerHints struct {
	// DifferentHost will place the share on a different back-end that does not
	// host the given shares.
	DifferentHost string `json:"different_host,omitempty"`

	// SameHost will place the share on a back-end that hosts the given shares.
	SameHost string `json:"same_host,omitempty"`

	// OnlyHost value must be a manage-share service host in
	// host@backend#POOL format (admin only). Only available in and beyond
	// API version 2.67
	OnlyHost string `json:"only_host,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]any, error)
}

// CreateOpts contains the options for create a Share. This object is
// passed to shares.Create(). For more information about these parameters,
// please refer to the Share object, or the shared file systems API v2
// documentation
type CreateOpts struct {
	// Defines the share protocol to use
	ShareProto string `json:"share_proto" required:"true"`
	// Size in GB
	Size int `json:"size" required:"true"`
	// Defines the share name
	Name string `json:"name,omitempty"`
	// Share description
	Description string `json:"description,omitempty"`
	// DisplayName is equivalent to Name. The API supports using both
	// This is an inherited attribute from the block storage API
	DisplayName string `json:"display_name,omitempty"`
	// DisplayDescription is equivalent to Description. The API supports using both
	// This is an inherited attribute from the block storage API
	DisplayDescription string `json:"display_description,omitempty"`
	// ShareType defines the sharetype. If omitted, a default share type is used
	ShareType string `json:"share_type,omitempty"`
	// VolumeType is deprecated but supported. Either ShareType or VolumeType can be used
	VolumeType string `json:"volume_type,omitempty"`
	// The UUID from which to create a share
	SnapshotID string `json:"snapshot_id,omitempty"`
	// Determines whether or not the share is public
	IsPublic *bool `json:"is_public,omitempty"`
	// The UUID of the share group. Available starting from the microversion 2.31
	ShareGroupID string `json:"share_group_id,omitempty"`
	// Key value pairs of user defined metadata
	Metadata map[string]string `json:"metadata,omitempty"`
	// The UUID of the share network to which the share belongs to
	ShareNetworkID string `json:"share_network_id,omitempty"`
	// The UUID of the consistency group to which the share belongs to
	ConsistencyGroupID string `json:"consistency_group_id,omitempty"`
	// The availability zone of the share
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// SchedulerHints are hints for the scheduler to select the share backend
	// Only available in and beyond API version 2.65
	SchedulerHints *SchedulerHints `json:"scheduler_hints,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "share")
}

// Create will create a new Share based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOpts holds options for listing Shares. It is passed to the
// shares.List function.
type ListOpts struct {
	// (Admin only). Defines whether to list the requested resources for all projects.
	AllTenants bool `q:"all_tenants"`
	// The share name.
	Name string `q:"name"`
	// Filters by a share status.
	Status string `q:"status"`
	// The UUID of the share server.
	ShareServerID string `q:"share_server_id"`
	// One or more metadata key and value pairs as a dictionary of strings.
	Metadata map[string]string `q:"metadata"`
	// The extra specifications for the share type.
	ExtraSpecs map[string]string `q:"extra_specs"`
	// The UUID of the share type.
	ShareTypeID string `q:"share_type_id"`
	// The maximum number of shares to return.
	Limit int `q:"limit"`
	// The offset to define start point of share or share group listing.
	Offset int `q:"offset"`
	// The key to sort a list of shares.
	SortKey string `q:"sort_key"`
	// The direction to sort a list of shares.
	SortDir string `q:"sort_dir"`
	// The UUID of the share’s base snapshot to filter the request based on.
	SnapshotID string `q:"snapshot_id"`
	// The share host name.
	Host string `q:"host"`
	// The share network ID.
	ShareNetworkID string `q:"share_network_id"`
	// The UUID of the project in which the share was created. Useful with all_tenants parameter.
	ProjectID string `q:"project_id"`
	// The level of visibility for the share.
	IsPublic *bool `q:"is_public"`
	// The UUID of a share group to filter resource.
	ShareGroupID string `q:"share_group_id"`
	// The export location UUID that can be used to filter shares or share instances.
	ExportLocationID string `q:"export_location_id"`
	// The export location path that can be used to filter shares or share instances.
	ExportLocationPath string `q:"export_location_path"`
	// The name pattern that can be used to filter shares, share snapshots, share networks or share groups.
	NamePattern string `q:"name~"`
	// The description pattern that can be used to filter shares, share snapshots, share networks or share groups.
	DescriptionPattern string `q:"description~"`
	// Whether to show count in API response or not, default is False.
	WithCount bool `q:"with_count"`
	// DisplayName is equivalent to Name. The API supports using both
	// This is an inherited attribute from the block storage API
	DisplayName string `q:"display_name"`
	// Equivalent to NamePattern.
	DisplayNamePattern string `q:"display_name~"`
	// VolumeTypeID is deprecated but supported. Either ShareTypeID or VolumeTypeID can be used
	VolumeTypeID string `q:"volume_type_id"`
	// The UUID of the share group snapshot.
	ShareGroupSnapshotID string `q:"share_group_snapshot_id"`
	// DisplayDescription is equivalent to Description. The API supports using both
	// This is an inherited attribute from the block storage API
	DisplayDescription string `q:"display_description"`
	// Equivalent to DescriptionPattern
	DisplayDescriptionPattern string `q:"display_description~"`
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToShareListQuery() (string, error)
}

// ToShareListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToShareListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail returns []Share optionally limited by the conditions provided in ListOpts.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToShareListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := SharePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// Delete will delete an existing Share with the given UUID.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get will get a single share with given UUID
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListExportLocations will list shareID's export locations.
// Client must have Microversion set; minimum supported microversion for ListExportLocations is 2.9.
func ListExportLocations(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ListExportLocationsResult) {
	resp, err := client.Get(ctx, listExportLocationsURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetExportLocation will get shareID's export location by an ID.
// Client must have Microversion set; minimum supported microversion for GetExportLocation is 2.9.
func GetExportLocation(ctx context.Context, client *gophercloud.ServiceClient, shareID string, id string) (r GetExportLocationResult) {
	resp, err := client.Get(ctx, getExportLocationURL(client, shareID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GrantAccessOptsBuilder allows extensions to add additional parameters to the
// GrantAccess request.
type GrantAccessOptsBuilder interface {
	ToGrantAccessMap() (map[string]any, error)
}

// GrantAccessOpts contains the options for creation of an GrantAccess request.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Grant Access documentation
type GrantAccessOpts struct {
	// The access rule type that can be "ip", "cert" or "user".
	AccessType string `json:"access_type"`
	// The value that defines the access that can be a valid format of IP, cert or user.
	AccessTo string `json:"access_to"`
	// The access level to the share is either "rw" or "ro".
	AccessLevel string `json:"access_level"`
}

// ToGrantAccessMap assembles a request body based on the contents of a
// GrantAccessOpts.
func (opts GrantAccessOpts) ToGrantAccessMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "allow_access")
}

// GrantAccess will grant access to a Share based on the values in GrantAccessOpts. To extract
// the GrantAccess object from the response, call the Extract method on the GrantAccessResult.
// Client must have Microversion set; minimum supported microversion for GrantAccess is 2.7.
func GrantAccess(ctx context.Context, client *gophercloud.ServiceClient, id string, opts GrantAccessOptsBuilder) (r GrantAccessResult) {
	b, err := opts.ToGrantAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, grantAccessURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RevokeAccessOptsBuilder allows extensions to add additional parameters to the
// RevokeAccess request.
type RevokeAccessOptsBuilder interface {
	ToRevokeAccessMap() (map[string]any, error)
}

// RevokeAccessOpts contains the options for creation of a RevokeAccess request.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Revoke Access documentation
type RevokeAccessOpts struct {
	AccessID string `json:"access_id"`
}

// ToRevokeAccessMap assembles a request body based on the contents of a
// RevokeAccessOpts.
func (opts RevokeAccessOpts) ToRevokeAccessMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "deny_access")
}

// RevokeAccess will revoke an existing access to a Share based on the values in RevokeAccessOpts.
// RevokeAccessResult contains only the error. To extract it, call the ExtractErr method on
// the RevokeAccessResult. Client must have Microversion set; minimum supported microversion
// for RevokeAccess is 2.7.
func RevokeAccess(ctx context.Context, client *gophercloud.ServiceClient, id string, opts RevokeAccessOptsBuilder) (r RevokeAccessResult) {
	b, err := opts.ToRevokeAccessMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, revokeAccessURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAccessRights lists all access rules assigned to a Share based on its id. To extract
// the AccessRight slice from the response, call the Extract method on the ListAccessRightsResult.
// Client must have Microversion set; minimum supported microversion for ListAccessRights is 2.7.
func ListAccessRights(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ListAccessRightsResult) {
	requestBody := map[string]any{"access_list": nil}
	resp, err := client.Post(ctx, listAccessRightsURL(client, id), requestBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ExtendOptsBuilder allows extensions to add additional parameters to the
// Extend request.
type ExtendOptsBuilder interface {
	ToShareExtendMap() (map[string]any, error)
}

// ExtendOpts contains options for extending a Share.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Extend share documentation
type ExtendOpts struct {
	// New size in GBs.
	NewSize int `json:"new_size"`
}

// ToShareExtendMap assembles a request body based on the contents of a
// ExtendOpts.
func (opts ExtendOpts) ToShareExtendMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "extend")
}

// Extend will extend the capacity of an existing share. ExtendResult contains only the error.
// To extract it, call the ExtractErr method on the ExtendResult.
// Client must have Microversion set; minimum supported microversion for Extend is 2.7.
func Extend(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ExtendOptsBuilder) (r ExtendResult) {
	b, err := opts.ToShareExtendMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, extendURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ShrinkOptsBuilder allows extensions to add additional parameters to the
// Shrink request.
type ShrinkOptsBuilder interface {
	ToShareShrinkMap() (map[string]any, error)
}

// ShrinkOpts contains options for shrinking a Share.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Shrink share documentation
type ShrinkOpts struct {
	// New size in GBs.
	NewSize int `json:"new_size"`
}

// ToShareShrinkMap assembles a request body based on the contents of a
// ShrinkOpts.
func (opts ShrinkOpts) ToShareShrinkMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "shrink")
}

// Shrink will shrink the capacity of an existing share. ShrinkResult contains only the error.
// To extract it, call the ExtractErr method on the ShrinkResult.
// Client must have Microversion set; minimum supported microversion for Shrink is 2.7.
func Shrink(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ShrinkOptsBuilder) (r ShrinkResult) {
	b, err := opts.ToShareShrinkMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, shrinkURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToShareUpdateMap() (map[string]any, error)
}

// UpdateOpts contain options for updating an existing Share. This object is passed
// to the share.Update function. For more information about the parameters, see
// the Share object.
type UpdateOpts struct {
	// Share name. Manila share update logic doesn't have a "name" alias.
	DisplayName *string `json:"display_name,omitempty"`
	// Share description. Manila share update logic doesn't have a "description" alias.
	DisplayDescription *string `json:"display_description,omitempty"`
	// Determines whether or not the share is public
	IsPublic *bool `json:"is_public,omitempty"`
}

// ToShareUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToShareUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "share")
}

// Update will update the Share with provided information. To extract the updated
// Share from the response, call the Extract method on the UpdateResult.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToShareUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMetadata retrieves metadata of the specified share. To extract the retrieved
// metadata from the response, call the Extract method on the MetadataResult.
func GetMetadata(ctx context.Context, client *gophercloud.ServiceClient, id string) (r MetadataResult) {
	resp, err := client.Get(ctx, getMetadataURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMetadatum retrieves a single metadata item of the specified share. To extract the retrieved
// metadata from the response, call the Extract method on the GetMetadatumResult.
func GetMetadatum(ctx context.Context, client *gophercloud.ServiceClient, id, key string) (r GetMetadatumResult) {
	resp, err := client.Get(ctx, getMetadatumURL(client, id, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetMetadataOpts contains options for setting share metadata.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Metadata, Show share metadata documentation.
type SetMetadataOpts struct {
	Metadata map[string]string `json:"metadata"`
}

// ToSetMetadataMap assembles a request body based on the contents of an
// SetMetadataOpts.
func (opts SetMetadataOpts) ToSetMetadataMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// SetMetadataOptsBuilder allows extensions to add additional parameters to the
// SetMetadata request.
type SetMetadataOptsBuilder interface {
	ToSetMetadataMap() (map[string]any, error)
}

// SetMetadata sets metadata of the specified share.
// Existing metadata items are either kept or overwritten by the metadata from the request.
// To extract the updated metadata from the response, call the Extract
// method on the MetadataResult.
func SetMetadata(ctx context.Context, client *gophercloud.ServiceClient, id string, opts SetMetadataOptsBuilder) (r MetadataResult) {
	b, err := opts.ToSetMetadataMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, setMetadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMetadataOpts contains options for updating share metadata.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Metadata, Update share metadata documentation.
type UpdateMetadataOpts struct {
	Metadata map[string]string `json:"metadata"`
}

// ToUpdateMetadataMap assembles a request body based on the contents of an
// UpdateMetadataOpts.
func (opts UpdateMetadataOpts) ToUpdateMetadataMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateMetadataOptsBuilder allows extensions to add additional parameters to the
// UpdateMetadata request.
type UpdateMetadataOptsBuilder interface {
	ToUpdateMetadataMap() (map[string]any, error)
}

// UpdateMetadata updates metadata of the specified share.
// All existing metadata items are discarded and replaced by the metadata from the request.
// To extract the updated metadata from the response, call the Extract
// method on the MetadataResult.
func UpdateMetadata(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateMetadataOptsBuilder) (r MetadataResult) {
	b, err := opts.ToUpdateMetadataMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, updateMetadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMetadatum deletes a single key-value pair from the metadata of the specified share.
func DeleteMetadatum(ctx context.Context, client *gophercloud.ServiceClient, id, key string) (r DeleteMetadatumResult) {
	resp, err := client.Delete(ctx, deleteMetadatumURL(client, id, key), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RevertOptsBuilder allows extensions to add additional parameters to the
// Revert request.
type RevertOptsBuilder interface {
	ToShareRevertMap() (map[string]any, error)
}

// RevertOpts contains options for reverting a Share to a snapshot.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Revert share documentation.
// Available only since Manila Microversion 2.27
type RevertOpts struct {
	// SnapshotID is a Snapshot ID to revert a Share to
	SnapshotID string `json:"snapshot_id"`
}

// ToShareRevertMap assembles a request body based on the contents of a
// RevertOpts.
func (opts RevertOpts) ToShareRevertMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "revert")
}

// Revert will revert the existing share to a Snapshot. RevertResult contains only the error.
// To extract it, call the ExtractErr method on the RevertResult.
// Client must have Microversion set; minimum supported microversion for Revert is 2.27.
func Revert(ctx context.Context, client *gophercloud.ServiceClient, id string, opts RevertOptsBuilder) (r RevertResult) {
	b, err := opts.ToShareRevertMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, revertURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetStatusOptsBuilder allows extensions to add additional parameters to the
// ResetStatus request.
type ResetStatusOptsBuilder interface {
	ToShareResetStatusMap() (map[string]any, error)
}

// ResetStatusOpts contains options for resetting a Share status.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, ResetStatus share documentation.
type ResetStatusOpts struct {
	// Status is a share status to reset to. Must be "new", "error" or "active".
	Status string `json:"status"`
}

// ToShareResetStatusMap assembles a request body based on the contents of a
// ResetStatusOpts.
func (opts ResetStatusOpts) ToShareResetStatusMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "reset_status")
}

// ResetStatus will reset the existing share status. ResetStatusResult contains only the error.
// To extract it, call the ExtractErr method on the ResetStatusResult.
// Client must have Microversion set; minimum supported microversion for ResetStatus is 2.7.
func ResetStatus(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResetStatusOptsBuilder) (r ResetStatusResult) {
	b, err := opts.ToShareResetStatusMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, resetStatusURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ForceDelete will delete the existing share in any state. ForceDeleteResult contains only the error.
// To extract it, call the ExtractErr method on the ForceDeleteResult.
// Client must have Microversion set; minimum supported microversion for ForceDelete is 2.7.
func ForceDelete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ForceDeleteResult) {
	b := map[string]any{
		"force_delete": nil,
	}
	resp, err := client.Post(ctx, forceDeleteURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unmanage will remove a share from the management of the Shared File System
// service without deleting the share. UnmanageResult contains only the error.
// To extract it, call the ExtractErr method on the UnmanageResult.
// Client must have Microversion set; minimum supported microversion for Unmanage is 2.7.
func Unmanage(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnmanageResult) {
	b := map[string]any{
		"unmanage": nil,
	}
	resp, err := client.Post(ctx, unmanageURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
