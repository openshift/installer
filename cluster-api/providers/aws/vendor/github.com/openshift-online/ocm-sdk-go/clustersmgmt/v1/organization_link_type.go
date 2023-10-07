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

// OrganizationLink represents the values of the 'organization_link' type.
//
// Definition of an organization link.
type OrganizationLink struct {
	bitmap_ uint32
	href    string
	id      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OrganizationLink) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HREF returns the value of the 'HREF' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HREF for the Organization, filled in response.
func (o *OrganizationLink) HREF() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the value of the 'HREF' attribute and
// a flag indicating if the attribute has a value.
//
// HREF for the Organization, filled in response.
func (o *OrganizationLink) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.href
	}
	return
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The organization's ID.
func (o *OrganizationLink) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// The organization's ID.
func (o *OrganizationLink) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// OrganizationLinkListKind is the name of the type used to represent list of objects of
// type 'organization_link'.
const OrganizationLinkListKind = "OrganizationLinkList"

// OrganizationLinkListLinkKind is the name of the type used to represent links to list
// of objects of type 'organization_link'.
const OrganizationLinkListLinkKind = "OrganizationLinkListLink"

// OrganizationLinkNilKind is the name of the type used to nil lists of objects of
// type 'organization_link'.
const OrganizationLinkListNilKind = "OrganizationLinkListNil"

// OrganizationLinkList is a list of values of the 'organization_link' type.
type OrganizationLinkList struct {
	href  string
	link  bool
	items []*OrganizationLink
}

// Len returns the length of the list.
func (l *OrganizationLinkList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *OrganizationLinkList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OrganizationLinkList) Get(i int) *OrganizationLink {
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
func (l *OrganizationLinkList) Slice() []*OrganizationLink {
	var slice []*OrganizationLink
	if l == nil {
		slice = make([]*OrganizationLink, 0)
	} else {
		slice = make([]*OrganizationLink, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OrganizationLinkList) Each(f func(item *OrganizationLink) bool) {
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
func (l *OrganizationLinkList) Range(f func(index int, item *OrganizationLink) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
