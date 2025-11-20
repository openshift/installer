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

// ResourceReview represents the values of the 'resource_review' type.
//
// Contains the result of performing a resource access review.
type ResourceReview struct {
	fieldSet_       []bool
	accountUsername string
	action          string
	clusterIDs      []string
	clusterUUIDs    []string
	organizationIDs []string
	resourceType    string
	subscriptionIDs []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ResourceReview) Empty() bool {
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
// Name of the account that is trying to perform the access.
func (o *ResourceReview) AccountUsername() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the account that is trying to perform the access.
func (o *ResourceReview) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.accountUsername
	}
	return
}

// Action returns the value of the 'action' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Action that will the user is trying to perform.
func (o *ResourceReview) Action() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.action
	}
	return ""
}

// GetAction returns the value of the 'action' attribute and
// a flag indicating if the attribute has a value.
//
// Action that will the user is trying to perform.
func (o *ResourceReview) GetAction() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.action
	}
	return
}

// ClusterIDs returns the value of the 'cluster_IDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifiers of the Clusters (internal ids) that the user has permission to perform the action upon.
func (o *ResourceReview) ClusterIDs() []string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.clusterIDs
	}
	return nil
}

// GetClusterIDs returns the value of the 'cluster_IDs' attribute and
// a flag indicating if the attribute has a value.
//
// Identifiers of the Clusters (internal ids) that the user has permission to perform the action upon.
func (o *ResourceReview) GetClusterIDs() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.clusterIDs
	}
	return
}

// ClusterUUIDs returns the value of the 'cluster_UUIDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifiers which Clusters (external ids) that the user has permission to perform the action upon.
func (o *ResourceReview) ClusterUUIDs() []string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.clusterUUIDs
	}
	return nil
}

// GetClusterUUIDs returns the value of the 'cluster_UUIDs' attribute and
// a flag indicating if the attribute has a value.
//
// Identifiers which Clusters (external ids) that the user has permission to perform the action upon.
func (o *ResourceReview) GetClusterUUIDs() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.clusterUUIDs
	}
	return
}

// OrganizationIDs returns the value of the 'organization_IDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifiers of the organizations that the user has permissions to perform the action
// upon.
func (o *ResourceReview) OrganizationIDs() []string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.organizationIDs
	}
	return nil
}

// GetOrganizationIDs returns the value of the 'organization_IDs' attribute and
// a flag indicating if the attribute has a value.
//
// Identifiers of the organizations that the user has permissions to perform the action
// upon.
func (o *ResourceReview) GetOrganizationIDs() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.organizationIDs
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of resource.
func (o *ResourceReview) ResourceType() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of resource.
func (o *ResourceReview) GetResourceType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.resourceType
	}
	return
}

// SubscriptionIDs returns the value of the 'subscription_IDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifiers of the subscriptions that the user has permission to perform the action upon.
func (o *ResourceReview) SubscriptionIDs() []string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.subscriptionIDs
	}
	return nil
}

// GetSubscriptionIDs returns the value of the 'subscription_IDs' attribute and
// a flag indicating if the attribute has a value.
//
// Identifiers of the subscriptions that the user has permission to perform the action upon.
func (o *ResourceReview) GetSubscriptionIDs() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.subscriptionIDs
	}
	return
}

// ResourceReviewListKind is the name of the type used to represent list of objects of
// type 'resource_review'.
const ResourceReviewListKind = "ResourceReviewList"

// ResourceReviewListLinkKind is the name of the type used to represent links to list
// of objects of type 'resource_review'.
const ResourceReviewListLinkKind = "ResourceReviewListLink"

// ResourceReviewNilKind is the name of the type used to nil lists of objects of
// type 'resource_review'.
const ResourceReviewListNilKind = "ResourceReviewListNil"

// ResourceReviewList is a list of values of the 'resource_review' type.
type ResourceReviewList struct {
	href  string
	link  bool
	items []*ResourceReview
}

// Len returns the length of the list.
func (l *ResourceReviewList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ResourceReviewList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ResourceReviewList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ResourceReviewList) SetItems(items []*ResourceReview) {
	l.items = items
}

// Items returns the items of the list.
func (l *ResourceReviewList) Items() []*ResourceReview {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ResourceReviewList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ResourceReviewList) Get(i int) *ResourceReview {
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
func (l *ResourceReviewList) Slice() []*ResourceReview {
	var slice []*ResourceReview
	if l == nil {
		slice = make([]*ResourceReview, 0)
	} else {
		slice = make([]*ResourceReview, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ResourceReviewList) Each(f func(item *ResourceReview) bool) {
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
func (l *ResourceReviewList) Range(f func(index int, item *ResourceReview) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
