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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	time "time"
)

// RegistryAllowlistKind is the name of the type used to represent objects
// of type 'registry_allowlist'.
const RegistryAllowlistKind = "RegistryAllowlist"

// RegistryAllowlistLinkKind is the name of the type used to represent links
// to objects of type 'registry_allowlist'.
const RegistryAllowlistLinkKind = "RegistryAllowlistLink"

// RegistryAllowlistNilKind is the name of the type used to nil references
// to objects of type 'registry_allowlist'.
const RegistryAllowlistNilKind = "RegistryAllowlistNil"

// RegistryAllowlist represents the values of the 'registry_allowlist' type.
//
// RegistryAllowlist represents a single registry allowlist.
type RegistryAllowlist struct {
	bitmap_           uint32
	id                string
	href              string
	cloudProvider     *CloudProvider
	creationTimestamp time.Time
	registries        []string
}

// Kind returns the name of the type of the object.
func (o *RegistryAllowlist) Kind() string {
	if o == nil {
		return RegistryAllowlistNilKind
	}
	if o.bitmap_&1 != 0 {
		return RegistryAllowlistLinkKind
	}
	return RegistryAllowlistKind
}

// Link returns true if this is a link.
func (o *RegistryAllowlist) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *RegistryAllowlist) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *RegistryAllowlist) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *RegistryAllowlist) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *RegistryAllowlist) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RegistryAllowlist) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CloudProvider is the cloud provider for which this allowlist is valid.
func (o *RegistryAllowlist) CloudProvider() *CloudProvider {
	if o != nil && o.bitmap_&8 != 0 {
		return o.cloudProvider
	}
	return nil
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// CloudProvider is the cloud provider for which this allowlist is valid.
func (o *RegistryAllowlist) GetCloudProvider() (value *CloudProvider, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.cloudProvider
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CreationTimestamp is the date and time when the allow list has been created.
func (o *RegistryAllowlist) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// CreationTimestamp is the date and time when the allow list has been created.
func (o *RegistryAllowlist) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// Registries returns the value of the 'registries' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Registries is the list of registries contained in this Allowlist.
func (o *RegistryAllowlist) Registries() []string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.registries
	}
	return nil
}

// GetRegistries returns the value of the 'registries' attribute and
// a flag indicating if the attribute has a value.
//
// Registries is the list of registries contained in this Allowlist.
func (o *RegistryAllowlist) GetRegistries() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.registries
	}
	return
}

// RegistryAllowlistListKind is the name of the type used to represent list of objects of
// type 'registry_allowlist'.
const RegistryAllowlistListKind = "RegistryAllowlistList"

// RegistryAllowlistListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry_allowlist'.
const RegistryAllowlistListLinkKind = "RegistryAllowlistListLink"

// RegistryAllowlistNilKind is the name of the type used to nil lists of objects of
// type 'registry_allowlist'.
const RegistryAllowlistListNilKind = "RegistryAllowlistListNil"

// RegistryAllowlistList is a list of values of the 'registry_allowlist' type.
type RegistryAllowlistList struct {
	href  string
	link  bool
	items []*RegistryAllowlist
}

// Kind returns the name of the type of the object.
func (l *RegistryAllowlistList) Kind() string {
	if l == nil {
		return RegistryAllowlistListNilKind
	}
	if l.link {
		return RegistryAllowlistListLinkKind
	}
	return RegistryAllowlistListKind
}

// Link returns true iif this is a link.
func (l *RegistryAllowlistList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *RegistryAllowlistList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *RegistryAllowlistList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *RegistryAllowlistList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *RegistryAllowlistList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *RegistryAllowlistList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *RegistryAllowlistList) SetItems(items []*RegistryAllowlist) {
	l.items = items
}

// Items returns the items of the list.
func (l *RegistryAllowlistList) Items() []*RegistryAllowlist {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *RegistryAllowlistList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RegistryAllowlistList) Get(i int) *RegistryAllowlist {
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
func (l *RegistryAllowlistList) Slice() []*RegistryAllowlist {
	var slice []*RegistryAllowlist
	if l == nil {
		slice = make([]*RegistryAllowlist, 0)
	} else {
		slice = make([]*RegistryAllowlist, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RegistryAllowlistList) Each(f func(item *RegistryAllowlist) bool) {
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
func (l *RegistryAllowlistList) Range(f func(index int, item *RegistryAllowlist) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
