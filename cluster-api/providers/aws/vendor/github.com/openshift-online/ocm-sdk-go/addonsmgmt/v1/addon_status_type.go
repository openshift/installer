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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonStatusKind is the name of the type used to represent objects
// of type 'addon_status'.
const AddonStatusKind = "AddonStatus"

// AddonStatusLinkKind is the name of the type used to represent links
// to objects of type 'addon_status'.
const AddonStatusLinkKind = "AddonStatusLink"

// AddonStatusNilKind is the name of the type used to nil references
// to objects of type 'addon_status'.
const AddonStatusNilKind = "AddonStatusNil"

// AddonStatus represents the values of the 'addon_status' type.
//
// Representation of an addon status.
type AddonStatus struct {
	bitmap_          uint32
	id               string
	href             string
	addonId          string
	correlationID    string
	statusConditions []*AddonStatusCondition
	version          string
}

// Kind returns the name of the type of the object.
func (o *AddonStatus) Kind() string {
	if o == nil {
		return AddonStatusNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddonStatusLinkKind
	}
	return AddonStatusKind
}

// Link returns true if this is a link.
func (o *AddonStatus) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddonStatus) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddonStatus) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddonStatus) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddonStatus) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonStatus) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AddonId returns the value of the 'addon_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the addon whose status belongs to.
func (o *AddonStatus) AddonId() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.addonId
	}
	return ""
}

// GetAddonId returns the value of the 'addon_id' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the addon whose status belongs to.
func (o *AddonStatus) GetAddonId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.addonId
	}
	return
}

// CorrelationID returns the value of the 'correlation_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifier for co-relating current AddonCR revision and reported status.
func (o *AddonStatus) CorrelationID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.correlationID
	}
	return ""
}

// GetCorrelationID returns the value of the 'correlation_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Identifier for co-relating current AddonCR revision and reported status.
func (o *AddonStatus) GetCorrelationID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.correlationID
	}
	return
}

// StatusConditions returns the value of the 'status_conditions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of reported addon status conditions
func (o *AddonStatus) StatusConditions() []*AddonStatusCondition {
	if o != nil && o.bitmap_&32 != 0 {
		return o.statusConditions
	}
	return nil
}

// GetStatusConditions returns the value of the 'status_conditions' attribute and
// a flag indicating if the attribute has a value.
//
// List of reported addon status conditions
func (o *AddonStatus) GetStatusConditions() (value []*AddonStatusCondition, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.statusConditions
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version of the addon reporting the status
func (o *AddonStatus) Version() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Version of the addon reporting the status
func (o *AddonStatus) GetVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.version
	}
	return
}

// AddonStatusListKind is the name of the type used to represent list of objects of
// type 'addon_status'.
const AddonStatusListKind = "AddonStatusList"

// AddonStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_status'.
const AddonStatusListLinkKind = "AddonStatusListLink"

// AddonStatusNilKind is the name of the type used to nil lists of objects of
// type 'addon_status'.
const AddonStatusListNilKind = "AddonStatusListNil"

// AddonStatusList is a list of values of the 'addon_status' type.
type AddonStatusList struct {
	href  string
	link  bool
	items []*AddonStatus
}

// Kind returns the name of the type of the object.
func (l *AddonStatusList) Kind() string {
	if l == nil {
		return AddonStatusListNilKind
	}
	if l.link {
		return AddonStatusListLinkKind
	}
	return AddonStatusListKind
}

// Link returns true iif this is a link.
func (l *AddonStatusList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddonStatusList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddonStatusList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddonStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddonStatusList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddonStatusList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddonStatusList) SetItems(items []*AddonStatus) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddonStatusList) Items() []*AddonStatus {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddonStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonStatusList) Get(i int) *AddonStatus {
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
func (l *AddonStatusList) Slice() []*AddonStatus {
	var slice []*AddonStatus
	if l == nil {
		slice = make([]*AddonStatus, 0)
	} else {
		slice = make([]*AddonStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonStatusList) Each(f func(item *AddonStatus) bool) {
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
func (l *AddonStatusList) Range(f func(index int, item *AddonStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
