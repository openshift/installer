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

import (
	time "time"
)

// ReservedResource represents the values of the 'reserved_resource' type.
type ReservedResource struct {
	bitmap_                   uint32
	availabilityZoneType      string
	billingMarketplaceAccount string
	billingModel              BillingModel
	count                     int
	createdAt                 time.Time
	resourceName              string
	resourceType              string
	scope                     string
	updatedAt                 time.Time
	byoc                      bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ReservedResource) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// BYOC returns the value of the 'BYOC' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) BYOC() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.byoc
	}
	return false
}

// GetBYOC returns the value of the 'BYOC' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetBYOC() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.byoc
	}
	return
}

// AvailabilityZoneType returns the value of the 'availability_zone_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) AvailabilityZoneType() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.availabilityZoneType
	}
	return ""
}

// GetAvailabilityZoneType returns the value of the 'availability_zone_type' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetAvailabilityZoneType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.availabilityZoneType
	}
	return
}

// BillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) BillingMarketplaceAccount() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.billingMarketplaceAccount
	}
	return ""
}

// GetBillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetBillingMarketplaceAccount() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.billingMarketplaceAccount
	}
	return
}

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) BillingModel() BillingModel {
	if o != nil && o.bitmap_&8 != 0 {
		return o.billingModel
	}
	return BillingModel("")
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetBillingModel() (value BillingModel, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.billingModel
	}
	return
}

// Count returns the value of the 'count' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) Count() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.count
	}
	return 0
}

// GetCount returns the value of the 'count' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetCount() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.count
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) CreatedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) ResourceName() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetResourceName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.resourceName
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) ResourceType() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetResourceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.resourceType
	}
	return
}

// Scope returns the value of the 'scope' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) Scope() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.scope
	}
	return ""
}

// GetScope returns the value of the 'scope' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetScope() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.scope
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ReservedResource) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&512 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *ReservedResource) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// ReservedResourceListKind is the name of the type used to represent list of objects of
// type 'reserved_resource'.
const ReservedResourceListKind = "ReservedResourceList"

// ReservedResourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'reserved_resource'.
const ReservedResourceListLinkKind = "ReservedResourceListLink"

// ReservedResourceNilKind is the name of the type used to nil lists of objects of
// type 'reserved_resource'.
const ReservedResourceListNilKind = "ReservedResourceListNil"

// ReservedResourceList is a list of values of the 'reserved_resource' type.
type ReservedResourceList struct {
	href  string
	link  bool
	items []*ReservedResource
}

// Len returns the length of the list.
func (l *ReservedResourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ReservedResourceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ReservedResourceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ReservedResourceList) SetItems(items []*ReservedResource) {
	l.items = items
}

// Items returns the items of the list.
func (l *ReservedResourceList) Items() []*ReservedResource {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ReservedResourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ReservedResourceList) Get(i int) *ReservedResource {
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
func (l *ReservedResourceList) Slice() []*ReservedResource {
	var slice []*ReservedResource
	if l == nil {
		slice = make([]*ReservedResource, 0)
	} else {
		slice = make([]*ReservedResource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ReservedResourceList) Each(f func(item *ReservedResource) bool) {
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
func (l *ReservedResourceList) Range(f func(index int, item *ReservedResource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
