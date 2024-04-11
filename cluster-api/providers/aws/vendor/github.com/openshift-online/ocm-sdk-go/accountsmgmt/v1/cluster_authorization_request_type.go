/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// ClusterAuthorizationRequest represents the values of the 'cluster_authorization_request' type.
type ClusterAuthorizationRequest struct {
	bitmap_           uint32
	accountUsername   string
	availabilityZone  string
	cloudAccountID    string
	cloudProviderID   string
	clusterID         string
	displayName       string
	externalClusterID string
	productID         string
	productCategory   string
	quotaVersion      string
	resources         []*ReservedResource
	scope             string
	byoc              bool
	disconnected      bool
	managed           bool
	reserve           bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterAuthorizationRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// BYOC returns the value of the 'BYOC' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) BYOC() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.byoc
	}
	return false
}

// GetBYOC returns the value of the 'BYOC' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetBYOC() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.byoc
	}
	return
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) AccountUsername() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.accountUsername
	}
	return
}

// AvailabilityZone returns the value of the 'availability_zone' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) AvailabilityZone() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.availabilityZone
	}
	return ""
}

// GetAvailabilityZone returns the value of the 'availability_zone' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetAvailabilityZone() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.availabilityZone
	}
	return
}

// CloudAccountID returns the value of the 'cloud_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) CloudAccountID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.cloudAccountID
	}
	return ""
}

// GetCloudAccountID returns the value of the 'cloud_account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetCloudAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.cloudAccountID
	}
	return
}

// CloudProviderID returns the value of the 'cloud_provider_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) CloudProviderID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cloudProviderID
	}
	return ""
}

// GetCloudProviderID returns the value of the 'cloud_provider_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetCloudProviderID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cloudProviderID
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) ClusterID() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetClusterID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.clusterID
	}
	return
}

// Disconnected returns the value of the 'disconnected' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) Disconnected() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.disconnected
	}
	return false
}

// GetDisconnected returns the value of the 'disconnected' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetDisconnected() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.disconnected
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) DisplayName() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.displayName
	}
	return
}

// ExternalClusterID returns the value of the 'external_cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) ExternalClusterID() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.externalClusterID
	}
	return ""
}

// GetExternalClusterID returns the value of the 'external_cluster_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetExternalClusterID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.externalClusterID
	}
	return
}

// Managed returns the value of the 'managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) Managed() bool {
	if o != nil && o.bitmap_&512 != 0 {
		return o.managed
	}
	return false
}

// GetManaged returns the value of the 'managed' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetManaged() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.managed
	}
	return
}

// ProductID returns the value of the 'product_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) ProductID() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.productID
	}
	return ""
}

// GetProductID returns the value of the 'product_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetProductID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.productID
	}
	return
}

// ProductCategory returns the value of the 'product_category' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) ProductCategory() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.productCategory
	}
	return ""
}

// GetProductCategory returns the value of the 'product_category' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetProductCategory() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.productCategory
	}
	return
}

// QuotaVersion returns the value of the 'quota_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) QuotaVersion() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.quotaVersion
	}
	return ""
}

// GetQuotaVersion returns the value of the 'quota_version' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetQuotaVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.quotaVersion
	}
	return
}

// Reserve returns the value of the 'reserve' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) Reserve() bool {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.reserve
	}
	return false
}

// GetReserve returns the value of the 'reserve' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetReserve() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.reserve
	}
	return
}

// Resources returns the value of the 'resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) Resources() []*ReservedResource {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.resources
	}
	return nil
}

// GetResources returns the value of the 'resources' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetResources() (value []*ReservedResource, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.resources
	}
	return
}

// Scope returns the value of the 'scope' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterAuthorizationRequest) Scope() string {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.scope
	}
	return ""
}

// GetScope returns the value of the 'scope' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterAuthorizationRequest) GetScope() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.scope
	}
	return
}

// ClusterAuthorizationRequestListKind is the name of the type used to represent list of objects of
// type 'cluster_authorization_request'.
const ClusterAuthorizationRequestListKind = "ClusterAuthorizationRequestList"

// ClusterAuthorizationRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_authorization_request'.
const ClusterAuthorizationRequestListLinkKind = "ClusterAuthorizationRequestListLink"

// ClusterAuthorizationRequestNilKind is the name of the type used to nil lists of objects of
// type 'cluster_authorization_request'.
const ClusterAuthorizationRequestListNilKind = "ClusterAuthorizationRequestListNil"

// ClusterAuthorizationRequestList is a list of values of the 'cluster_authorization_request' type.
type ClusterAuthorizationRequestList struct {
	href  string
	link  bool
	items []*ClusterAuthorizationRequest
}

// Len returns the length of the list.
func (l *ClusterAuthorizationRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterAuthorizationRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterAuthorizationRequestList) Get(i int) *ClusterAuthorizationRequest {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *ClusterAuthorizationRequestList) Slice() []*ClusterAuthorizationRequest {
	var slice []*ClusterAuthorizationRequest
	if l == nil {
		slice = make([]*ClusterAuthorizationRequest, 0)
	} else {
		slice = make([]*ClusterAuthorizationRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterAuthorizationRequestList) Each(f func(item *ClusterAuthorizationRequest) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *ClusterAuthorizationRequestList) Range(f func(index int, item *ClusterAuthorizationRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
