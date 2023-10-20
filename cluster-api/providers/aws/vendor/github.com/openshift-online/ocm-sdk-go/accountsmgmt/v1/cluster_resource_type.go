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

// ClusterResource represents the values of the 'cluster_resource' type.
type ClusterResource struct {
	bitmap_          uint32
	total            *ValueUnit
	updatedTimestamp time.Time
	used             *ValueUnit
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterResource) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Total returns the value of the 'total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterResource) Total() *ValueUnit {
	if o != nil && o.bitmap_&1 != 0 {
		return o.total
	}
	return nil
}

// GetTotal returns the value of the 'total' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterResource) GetTotal() (value *ValueUnit, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.total
	}
	return
}

// UpdatedTimestamp returns the value of the 'updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterResource) UpdatedTimestamp() time.Time {
	if o != nil && o.bitmap_&2 != 0 {
		return o.updatedTimestamp
	}
	return time.Time{}
}

// GetUpdatedTimestamp returns the value of the 'updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterResource) GetUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.updatedTimestamp
	}
	return
}

// Used returns the value of the 'used' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterResource) Used() *ValueUnit {
	if o != nil && o.bitmap_&4 != 0 {
		return o.used
	}
	return nil
}

// GetUsed returns the value of the 'used' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterResource) GetUsed() (value *ValueUnit, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.used
	}
	return
}

// ClusterResourceListKind is the name of the type used to represent list of objects of
// type 'cluster_resource'.
const ClusterResourceListKind = "ClusterResourceList"

// ClusterResourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_resource'.
const ClusterResourceListLinkKind = "ClusterResourceListLink"

// ClusterResourceNilKind is the name of the type used to nil lists of objects of
// type 'cluster_resource'.
const ClusterResourceListNilKind = "ClusterResourceListNil"

// ClusterResourceList is a list of values of the 'cluster_resource' type.
type ClusterResourceList struct {
	href  string
	link  bool
	items []*ClusterResource
}

// Len returns the length of the list.
func (l *ClusterResourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterResourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterResourceList) Get(i int) *ClusterResource {
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
func (l *ClusterResourceList) Slice() []*ClusterResource {
	var slice []*ClusterResource
	if l == nil {
		slice = make([]*ClusterResource, 0)
	} else {
		slice = make([]*ClusterResource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterResourceList) Each(f func(item *ClusterResource) bool) {
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
func (l *ClusterResourceList) Range(f func(index int, item *ClusterResource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
