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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

// SelfCapabilityReviewRequest represents the values of the 'self_capability_review_request' type.
//
// Representation of a capability review.
type SelfCapabilityReviewRequest struct {
	fieldSet_       []bool
	accountUsername string
	capability      string
	clusterID       string
	organizationID  string
	resourceType    string
	subscriptionID  string
	type_           string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SelfCapabilityReviewRequest) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the username of the account of which capability is being reviewed.
func (o *SelfCapabilityReviewRequest) AccountUsername() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the username of the account of which capability is being reviewed.
func (o *SelfCapabilityReviewRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.accountUsername
	}
	return
}

// Capability returns the value of the 'capability' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Capability to review [manage_cluster_admin].
func (o *SelfCapabilityReviewRequest) Capability() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.capability
	}
	return ""
}

// GetCapability returns the value of the 'capability' attribute and
// a flag indicating if the attribute has a value.
//
// Capability to review [manage_cluster_admin].
func (o *SelfCapabilityReviewRequest) GetCapability() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.capability
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Cluster (internal id) the resource type belongs to.
func (o *SelfCapabilityReviewRequest) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Cluster (internal id) the resource type belongs to.
func (o *SelfCapabilityReviewRequest) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.clusterID
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Organization the resource type belongs to.
func (o *SelfCapabilityReviewRequest) OrganizationID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Organization the resource type belongs to.
func (o *SelfCapabilityReviewRequest) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.organizationID
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the type of the resource.
// See uhc-account-manager/openapi/openapi.yaml for a list of possible values.
func (o *SelfCapabilityReviewRequest) ResourceType() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the type of the resource.
// See uhc-account-manager/openapi/openapi.yaml for a list of possible values.
func (o *SelfCapabilityReviewRequest) GetResourceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.resourceType
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Subscription the resource type belongs to.
func (o *SelfCapabilityReviewRequest) SubscriptionID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Subscription the resource type belongs to.
func (o *SelfCapabilityReviewRequest) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.subscriptionID
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of capability [Cluster].
func (o *SelfCapabilityReviewRequest) Type() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of capability [Cluster].
func (o *SelfCapabilityReviewRequest) GetType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.type_
	}
	return
}

// SelfCapabilityReviewRequestListKind is the name of the type used to represent list of objects of
// type 'self_capability_review_request'.
const SelfCapabilityReviewRequestListKind = "SelfCapabilityReviewRequestList"

// SelfCapabilityReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'self_capability_review_request'.
const SelfCapabilityReviewRequestListLinkKind = "SelfCapabilityReviewRequestListLink"

// SelfCapabilityReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'self_capability_review_request'.
const SelfCapabilityReviewRequestListNilKind = "SelfCapabilityReviewRequestListNil"

// SelfCapabilityReviewRequestList is a list of values of the 'self_capability_review_request' type.
type SelfCapabilityReviewRequestList struct {
	href  string
	link  bool
	items []*SelfCapabilityReviewRequest
}

// Len returns the length of the list.
func (l *SelfCapabilityReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SelfCapabilityReviewRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SelfCapabilityReviewRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SelfCapabilityReviewRequestList) SetItems(items []*SelfCapabilityReviewRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *SelfCapabilityReviewRequestList) Items() []*SelfCapabilityReviewRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SelfCapabilityReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SelfCapabilityReviewRequestList) Get(i int) *SelfCapabilityReviewRequest {
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
func (l *SelfCapabilityReviewRequestList) Slice() []*SelfCapabilityReviewRequest {
	var slice []*SelfCapabilityReviewRequest
	if l == nil {
		slice = make([]*SelfCapabilityReviewRequest, 0)
	} else {
		slice = make([]*SelfCapabilityReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SelfCapabilityReviewRequestList) Each(f func(item *SelfCapabilityReviewRequest) bool) {
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
func (l *SelfCapabilityReviewRequestList) Range(f func(index int, item *SelfCapabilityReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
