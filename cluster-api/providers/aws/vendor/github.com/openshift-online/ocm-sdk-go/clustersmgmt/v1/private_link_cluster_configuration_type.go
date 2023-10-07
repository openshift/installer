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

// PrivateLinkClusterConfiguration represents the values of the 'private_link_cluster_configuration' type.
//
// Manages the configuration for the Private Links.
type PrivateLinkClusterConfiguration struct {
	bitmap_    uint32
	principals []*PrivateLinkPrincipal
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *PrivateLinkClusterConfiguration) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Principals returns the value of the 'principals' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of additional principals for the Private Link
func (o *PrivateLinkClusterConfiguration) Principals() []*PrivateLinkPrincipal {
	if o != nil && o.bitmap_&1 != 0 {
		return o.principals
	}
	return nil
}

// GetPrincipals returns the value of the 'principals' attribute and
// a flag indicating if the attribute has a value.
//
// List of additional principals for the Private Link
func (o *PrivateLinkClusterConfiguration) GetPrincipals() (value []*PrivateLinkPrincipal, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.principals
	}
	return
}

// PrivateLinkClusterConfigurationListKind is the name of the type used to represent list of objects of
// type 'private_link_cluster_configuration'.
const PrivateLinkClusterConfigurationListKind = "PrivateLinkClusterConfigurationList"

// PrivateLinkClusterConfigurationListLinkKind is the name of the type used to represent links to list
// of objects of type 'private_link_cluster_configuration'.
const PrivateLinkClusterConfigurationListLinkKind = "PrivateLinkClusterConfigurationListLink"

// PrivateLinkClusterConfigurationNilKind is the name of the type used to nil lists of objects of
// type 'private_link_cluster_configuration'.
const PrivateLinkClusterConfigurationListNilKind = "PrivateLinkClusterConfigurationListNil"

// PrivateLinkClusterConfigurationList is a list of values of the 'private_link_cluster_configuration' type.
type PrivateLinkClusterConfigurationList struct {
	href  string
	link  bool
	items []*PrivateLinkClusterConfiguration
}

// Len returns the length of the list.
func (l *PrivateLinkClusterConfigurationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *PrivateLinkClusterConfigurationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PrivateLinkClusterConfigurationList) Get(i int) *PrivateLinkClusterConfiguration {
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
func (l *PrivateLinkClusterConfigurationList) Slice() []*PrivateLinkClusterConfiguration {
	var slice []*PrivateLinkClusterConfiguration
	if l == nil {
		slice = make([]*PrivateLinkClusterConfiguration, 0)
	} else {
		slice = make([]*PrivateLinkClusterConfiguration, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PrivateLinkClusterConfigurationList) Each(f func(item *PrivateLinkClusterConfiguration) bool) {
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
func (l *PrivateLinkClusterConfigurationList) Range(f func(index int, item *PrivateLinkClusterConfiguration) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
