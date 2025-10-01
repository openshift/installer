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

// PrivateLinkPrincipalKind is the name of the type used to represent objects
// of type 'private_link_principal'.
const PrivateLinkPrincipalKind = "PrivateLinkPrincipal"

// PrivateLinkPrincipalLinkKind is the name of the type used to represent links
// to objects of type 'private_link_principal'.
const PrivateLinkPrincipalLinkKind = "PrivateLinkPrincipalLink"

// PrivateLinkPrincipalNilKind is the name of the type used to nil references
// to objects of type 'private_link_principal'.
const PrivateLinkPrincipalNilKind = "PrivateLinkPrincipalNil"

// PrivateLinkPrincipal represents the values of the 'private_link_principal' type.
type PrivateLinkPrincipal struct {
	bitmap_   uint32
	id        string
	href      string
	principal string
}

// Kind returns the name of the type of the object.
func (o *PrivateLinkPrincipal) Kind() string {
	if o == nil {
		return PrivateLinkPrincipalNilKind
	}
	if o.bitmap_&1 != 0 {
		return PrivateLinkPrincipalLinkKind
	}
	return PrivateLinkPrincipalKind
}

// Link returns true if this is a link.
func (o *PrivateLinkPrincipal) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *PrivateLinkPrincipal) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *PrivateLinkPrincipal) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *PrivateLinkPrincipal) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *PrivateLinkPrincipal) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *PrivateLinkPrincipal) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Principal returns the value of the 'principal' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ARN for a principal that is allowed for this Private Link.
func (o *PrivateLinkPrincipal) Principal() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.principal
	}
	return ""
}

// GetPrincipal returns the value of the 'principal' attribute and
// a flag indicating if the attribute has a value.
//
// ARN for a principal that is allowed for this Private Link.
func (o *PrivateLinkPrincipal) GetPrincipal() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.principal
	}
	return
}

// PrivateLinkPrincipalListKind is the name of the type used to represent list of objects of
// type 'private_link_principal'.
const PrivateLinkPrincipalListKind = "PrivateLinkPrincipalList"

// PrivateLinkPrincipalListLinkKind is the name of the type used to represent links to list
// of objects of type 'private_link_principal'.
const PrivateLinkPrincipalListLinkKind = "PrivateLinkPrincipalListLink"

// PrivateLinkPrincipalNilKind is the name of the type used to nil lists of objects of
// type 'private_link_principal'.
const PrivateLinkPrincipalListNilKind = "PrivateLinkPrincipalListNil"

// PrivateLinkPrincipalList is a list of values of the 'private_link_principal' type.
type PrivateLinkPrincipalList struct {
	href  string
	link  bool
	items []*PrivateLinkPrincipal
}

// Kind returns the name of the type of the object.
func (l *PrivateLinkPrincipalList) Kind() string {
	if l == nil {
		return PrivateLinkPrincipalListNilKind
	}
	if l.link {
		return PrivateLinkPrincipalListLinkKind
	}
	return PrivateLinkPrincipalListKind
}

// Link returns true iif this is a link.
func (l *PrivateLinkPrincipalList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *PrivateLinkPrincipalList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *PrivateLinkPrincipalList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *PrivateLinkPrincipalList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalList) SetItems(items []*PrivateLinkPrincipal) {
	l.items = items
}

// Items returns the items of the list.
func (l *PrivateLinkPrincipalList) Items() []*PrivateLinkPrincipal {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *PrivateLinkPrincipalList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PrivateLinkPrincipalList) Get(i int) *PrivateLinkPrincipal {
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
func (l *PrivateLinkPrincipalList) Slice() []*PrivateLinkPrincipal {
	var slice []*PrivateLinkPrincipal
	if l == nil {
		slice = make([]*PrivateLinkPrincipal, 0)
	} else {
		slice = make([]*PrivateLinkPrincipal, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PrivateLinkPrincipalList) Each(f func(item *PrivateLinkPrincipal) bool) {
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
func (l *PrivateLinkPrincipalList) Range(f func(index int, item *PrivateLinkPrincipal) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
