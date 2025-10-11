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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// ManagedService represents the values of the 'managed_service' type.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type ManagedService struct {
	bitmap_ uint32
	enabled bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ManagedService) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the cluster belongs to a managed service
// This should only be set by the "Managed Service" service.
// clusters with this set can only be modified by the "Managed Service" service.
func (o *ManagedService) Enabled() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the cluster belongs to a managed service
// This should only be set by the "Managed Service" service.
// clusters with this set can only be modified by the "Managed Service" service.
func (o *ManagedService) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.enabled
	}
	return
}

// ManagedServiceListKind is the name of the type used to represent list of objects of
// type 'managed_service'.
const ManagedServiceListKind = "ManagedServiceList"

// ManagedServiceListLinkKind is the name of the type used to represent links to list
// of objects of type 'managed_service'.
const ManagedServiceListLinkKind = "ManagedServiceListLink"

// ManagedServiceNilKind is the name of the type used to nil lists of objects of
// type 'managed_service'.
const ManagedServiceListNilKind = "ManagedServiceListNil"

// ManagedServiceList is a list of values of the 'managed_service' type.
type ManagedServiceList struct {
	href  string
	link  bool
	items []*ManagedService
}

// Len returns the length of the list.
func (l *ManagedServiceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ManagedServiceList) SetItems(items []*ManagedService) {
	l.items = items
}

// Items returns the items of the list.
func (l *ManagedServiceList) Items() []*ManagedService {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ManagedServiceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ManagedServiceList) Get(i int) *ManagedService {
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
func (l *ManagedServiceList) Slice() []*ManagedService {
	var slice []*ManagedService
	if l == nil {
		slice = make([]*ManagedService, 0)
	} else {
		slice = make([]*ManagedService, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ManagedServiceList) Each(f func(item *ManagedService) bool) {
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
func (l *ManagedServiceList) Range(f func(index int, item *ManagedService) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
