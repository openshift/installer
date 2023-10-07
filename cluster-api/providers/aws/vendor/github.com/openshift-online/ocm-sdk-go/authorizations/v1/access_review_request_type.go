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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// AccessReviewRequest represents the values of the 'access_review_request' type.
//
// Representation of an access review
type AccessReviewRequest struct {
	bitmap_         uint32
	accountUsername string
	action          string
	clusterID       string
	clusterUUID     string
	organizationID  string
	resourceType    string
	subscriptionID  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccessReviewRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the username of the account of which access is being reviewed
func (o *AccessReviewRequest) AccountUsername() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the username of the account of which access is being reviewed
func (o *AccessReviewRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accountUsername
	}
	return
}

// Action returns the value of the 'action' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the action, one of: [get,list,create,delete,update]
func (o *AccessReviewRequest) Action() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.action
	}
	return ""
}

// GetAction returns the value of the 'action' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the action, one of: [get,list,create,delete,update]
func (o *AccessReviewRequest) GetAction() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.action
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Cluster (internal id) the resource type belongs to
func (o *AccessReviewRequest) ClusterID() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Cluster (internal id) the resource type belongs to
func (o *AccessReviewRequest) GetClusterID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.clusterID
	}
	return
}

// ClusterUUID returns the value of the 'cluster_UUID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Cluster (external id) the resource type belongs to
func (o *AccessReviewRequest) ClusterUUID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clusterUUID
	}
	return ""
}

// GetClusterUUID returns the value of the 'cluster_UUID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Cluster (external id) the resource type belongs to
func (o *AccessReviewRequest) GetClusterUUID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clusterUUID
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Organization the resource type belongs to
func (o *AccessReviewRequest) OrganizationID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Organization the resource type belongs to
func (o *AccessReviewRequest) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.organizationID
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the type of the resource an action would be taken on.
// See uhc-account-manager/openapi/openapi.yaml for a list of possible values
func (o *AccessReviewRequest) ResourceType() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the type of the resource an action would be taken on.
// See uhc-account-manager/openapi/openapi.yaml for a list of possible values
func (o *AccessReviewRequest) GetResourceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.resourceType
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Subscription the resource type belongs to
func (o *AccessReviewRequest) SubscriptionID() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Subscription the resource type belongs to
func (o *AccessReviewRequest) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.subscriptionID
	}
	return
}

// AccessReviewRequestListKind is the name of the type used to represent list of objects of
// type 'access_review_request'.
const AccessReviewRequestListKind = "AccessReviewRequestList"

// AccessReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'access_review_request'.
const AccessReviewRequestListLinkKind = "AccessReviewRequestListLink"

// AccessReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'access_review_request'.
const AccessReviewRequestListNilKind = "AccessReviewRequestListNil"

// AccessReviewRequestList is a list of values of the 'access_review_request' type.
type AccessReviewRequestList struct {
	href  string
	link  bool
	items []*AccessReviewRequest
}

// Len returns the length of the list.
func (l *AccessReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AccessReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccessReviewRequestList) Get(i int) *AccessReviewRequest {
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
func (l *AccessReviewRequestList) Slice() []*AccessReviewRequest {
	var slice []*AccessReviewRequest
	if l == nil {
		slice = make([]*AccessReviewRequest, 0)
	} else {
		slice = make([]*AccessReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccessReviewRequestList) Each(f func(item *AccessReviewRequest) bool) {
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
func (l *AccessReviewRequestList) Range(f func(index int, item *AccessReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
