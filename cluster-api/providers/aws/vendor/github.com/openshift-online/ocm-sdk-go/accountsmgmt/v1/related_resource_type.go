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

// RelatedResource represents the values of the 'related_resource' type.
//
// Resource which can be provisioned using the allowed quota.
type RelatedResource struct {
	bitmap_              uint32
	byoc                 string
	availabilityZoneType string
	billingModel         string
	cloudProvider        string
	cost                 int
	product              string
	resourceName         string
	resourceType         string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RelatedResource) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// BYOC returns the value of the 'BYOC' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) BYOC() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.byoc
	}
	return ""
}

// GetBYOC returns the value of the 'BYOC' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetBYOC() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.byoc
	}
	return
}

// AvailabilityZoneType returns the value of the 'availability_zone_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) AvailabilityZoneType() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.availabilityZoneType
	}
	return ""
}

// GetAvailabilityZoneType returns the value of the 'availability_zone_type' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetAvailabilityZoneType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.availabilityZoneType
	}
	return
}

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) BillingModel() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.billingModel
	}
	return ""
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetBillingModel() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.billingModel
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) CloudProvider() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.cloudProvider
	}
	return ""
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetCloudProvider() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// Cost returns the value of the 'cost' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) Cost() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cost
	}
	return 0
}

// GetCost returns the value of the 'cost' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetCost() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cost
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) Product() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.product
	}
	return ""
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetProduct() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.product
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) ResourceName() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetResourceName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.resourceName
	}
	return
}

// ResourceType returns the value of the 'resource_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RelatedResource) ResourceType() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.resourceType
	}
	return ""
}

// GetResourceType returns the value of the 'resource_type' attribute and
// a flag indicating if the attribute has a value.
func (o *RelatedResource) GetResourceType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.resourceType
	}
	return
}

// RelatedResourceListKind is the name of the type used to represent list of objects of
// type 'related_resource'.
const RelatedResourceListKind = "RelatedResourceList"

// RelatedResourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'related_resource'.
const RelatedResourceListLinkKind = "RelatedResourceListLink"

// RelatedResourceNilKind is the name of the type used to nil lists of objects of
// type 'related_resource'.
const RelatedResourceListNilKind = "RelatedResourceListNil"

// RelatedResourceList is a list of values of the 'related_resource' type.
type RelatedResourceList struct {
	href  string
	link  bool
	items []*RelatedResource
}

// Len returns the length of the list.
func (l *RelatedResourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *RelatedResourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RelatedResourceList) Get(i int) *RelatedResource {
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
func (l *RelatedResourceList) Slice() []*RelatedResource {
	var slice []*RelatedResource
	if l == nil {
		slice = make([]*RelatedResource, 0)
	} else {
		slice = make([]*RelatedResource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RelatedResourceList) Each(f func(item *RelatedResource) bool) {
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
func (l *RelatedResourceList) Range(f func(index int, item *RelatedResource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
