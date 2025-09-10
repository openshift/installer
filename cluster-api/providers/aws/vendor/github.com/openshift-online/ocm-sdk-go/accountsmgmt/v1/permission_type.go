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

// PermissionKind is the name of the type used to represent objects
// of type 'permission'.
const PermissionKind = "Permission"

// PermissionLinkKind is the name of the type used to represent links
// to objects of type 'permission'.
const PermissionLinkKind = "PermissionLink"

// PermissionNilKind is the name of the type used to nil references
// to objects of type 'permission'.
const PermissionNilKind = "PermissionNil"

// Permission represents the values of the 'permission' type.
type Permission struct {
	bitmap_  uint32
	id       string
	href     string
	action   Action
	resource string
}

// Kind returns the name of the type of the object.
func (o *Permission) Kind() string {
	if o == nil {
		return PermissionNilKind
	}
	if o.bitmap_&1 != 0 {
		return PermissionLinkKind
	}
	return PermissionKind
}

// Link returns true if this is a link.
func (o *Permission) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Permission) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Permission) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Permission) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Permission) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Permission) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Action returns the value of the 'action' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Permission) Action() Action {
	if o != nil && o.bitmap_&8 != 0 {
		return o.action
	}
	return Action("")
}

// GetAction returns the value of the 'action' attribute and
// a flag indicating if the attribute has a value.
func (o *Permission) GetAction() (value Action, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.action
	}
	return
}

// Resource returns the value of the 'resource' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Permission) Resource() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.resource
	}
	return ""
}

// GetResource returns the value of the 'resource' attribute and
// a flag indicating if the attribute has a value.
func (o *Permission) GetResource() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.resource
	}
	return
}

// PermissionListKind is the name of the type used to represent list of objects of
// type 'permission'.
const PermissionListKind = "PermissionList"

// PermissionListLinkKind is the name of the type used to represent links to list
// of objects of type 'permission'.
const PermissionListLinkKind = "PermissionListLink"

// PermissionNilKind is the name of the type used to nil lists of objects of
// type 'permission'.
const PermissionListNilKind = "PermissionListNil"

// PermissionList is a list of values of the 'permission' type.
type PermissionList struct {
	href  string
	link  bool
	items []*Permission
}

// Kind returns the name of the type of the object.
func (l *PermissionList) Kind() string {
	if l == nil {
		return PermissionListNilKind
	}
	if l.link {
		return PermissionListLinkKind
	}
	return PermissionListKind
}

// Link returns true iif this is a link.
func (l *PermissionList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *PermissionList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *PermissionList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *PermissionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *PermissionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *PermissionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *PermissionList) SetItems(items []*Permission) {
	l.items = items
}

// Items returns the items of the list.
func (l *PermissionList) Items() []*Permission {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *PermissionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PermissionList) Get(i int) *Permission {
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
func (l *PermissionList) Slice() []*Permission {
	var slice []*Permission
	if l == nil {
		slice = make([]*Permission, 0)
	} else {
		slice = make([]*Permission, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PermissionList) Each(f func(item *Permission) bool) {
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
func (l *PermissionList) Range(f func(index int, item *Permission) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
