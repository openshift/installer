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

// WifConfigKind is the name of the type used to represent objects
// of type 'wif_config'.
const WifConfigKind = "WifConfig"

// WifConfigLinkKind is the name of the type used to represent links
// to objects of type 'wif_config'.
const WifConfigLinkKind = "WifConfigLink"

// WifConfigNilKind is the name of the type used to nil references
// to objects of type 'wif_config'.
const WifConfigNilKind = "WifConfigNil"

// WifConfig represents the values of the 'wif_config' type.
//
// Definition of an wif_config resource.
type WifConfig struct {
	bitmap_      uint32
	id           string
	href         string
	displayName  string
	gcp          *WifGcp
	organization *OrganizationLink
}

// Kind returns the name of the type of the object.
func (o *WifConfig) Kind() string {
	if o == nil {
		return WifConfigNilKind
	}
	if o.bitmap_&1 != 0 {
		return WifConfigLinkKind
	}
	return WifConfigKind
}

// Link returns true iif this is a link.
func (o *WifConfig) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *WifConfig) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *WifConfig) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *WifConfig) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *WifConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifConfig) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name OCM clients will display for this wif_config.
func (o *WifConfig) DisplayName() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name OCM clients will display for this wif_config.
func (o *WifConfig) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.displayName
	}
	return
}

// Gcp returns the value of the 'gcp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Holds GCP related data.
func (o *WifConfig) Gcp() *WifGcp {
	if o != nil && o.bitmap_&16 != 0 {
		return o.gcp
	}
	return nil
}

// GetGcp returns the value of the 'gcp' attribute and
// a flag indicating if the attribute has a value.
//
// Holds GCP related data.
func (o *WifConfig) GetGcp() (value *WifGcp, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.gcp
	}
	return
}

// Organization returns the value of the 'organization' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The OCM organization that this wif_config resource belongs to.
func (o *WifConfig) Organization() *OrganizationLink {
	if o != nil && o.bitmap_&32 != 0 {
		return o.organization
	}
	return nil
}

// GetOrganization returns the value of the 'organization' attribute and
// a flag indicating if the attribute has a value.
//
// The OCM organization that this wif_config resource belongs to.
func (o *WifConfig) GetOrganization() (value *OrganizationLink, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.organization
	}
	return
}

// WifConfigListKind is the name of the type used to represent list of objects of
// type 'wif_config'.
const WifConfigListKind = "WifConfigList"

// WifConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_config'.
const WifConfigListLinkKind = "WifConfigListLink"

// WifConfigNilKind is the name of the type used to nil lists of objects of
// type 'wif_config'.
const WifConfigListNilKind = "WifConfigListNil"

// WifConfigList is a list of values of the 'wif_config' type.
type WifConfigList struct {
	href  string
	link  bool
	items []*WifConfig
}

// Kind returns the name of the type of the object.
func (l *WifConfigList) Kind() string {
	if l == nil {
		return WifConfigListNilKind
	}
	if l.link {
		return WifConfigListLinkKind
	}
	return WifConfigListKind
}

// Link returns true iif this is a link.
func (l *WifConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *WifConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *WifConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *WifConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *WifConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifConfigList) Get(i int) *WifConfig {
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
func (l *WifConfigList) Slice() []*WifConfig {
	var slice []*WifConfig
	if l == nil {
		slice = make([]*WifConfig, 0)
	} else {
		slice = make([]*WifConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifConfigList) Each(f func(item *WifConfig) bool) {
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
func (l *WifConfigList) Range(f func(index int, item *WifConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
