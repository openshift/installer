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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// PrivateLinkPrincipalsKind is the name of the type used to represent objects
// of type 'private_link_principals'.
const PrivateLinkPrincipalsKind = "PrivateLinkPrincipals"

// PrivateLinkPrincipalsLinkKind is the name of the type used to represent links
// to objects of type 'private_link_principals'.
const PrivateLinkPrincipalsLinkKind = "PrivateLinkPrincipalsLink"

// PrivateLinkPrincipalsNilKind is the name of the type used to nil references
// to objects of type 'private_link_principals'.
const PrivateLinkPrincipalsNilKind = "PrivateLinkPrincipalsNil"

// PrivateLinkPrincipals represents the values of the 'private_link_principals' type.
//
// Contains a list of principals for the Private Link.
type PrivateLinkPrincipals struct {
	fieldSet_  []bool
	id         string
	href       string
	principals []*PrivateLinkPrincipal
}

// Kind returns the name of the type of the object.
func (o *PrivateLinkPrincipals) Kind() string {
	if o == nil {
		return PrivateLinkPrincipalsNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return PrivateLinkPrincipalsLinkKind
	}
	return PrivateLinkPrincipalsKind
}

// Link returns true if this is a link.
func (o *PrivateLinkPrincipals) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *PrivateLinkPrincipals) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *PrivateLinkPrincipals) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *PrivateLinkPrincipals) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *PrivateLinkPrincipals) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *PrivateLinkPrincipals) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// Principals returns the value of the 'principals' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of additional principals for the Private Link
func (o *PrivateLinkPrincipals) Principals() []*PrivateLinkPrincipal {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.principals
	}
	return nil
}

// GetPrincipals returns the value of the 'principals' attribute and
// a flag indicating if the attribute has a value.
//
// List of additional principals for the Private Link
func (o *PrivateLinkPrincipals) GetPrincipals() (value []*PrivateLinkPrincipal, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.principals
	}
	return
}

// PrivateLinkPrincipalsListKind is the name of the type used to represent list of objects of
// type 'private_link_principals'.
const PrivateLinkPrincipalsListKind = "PrivateLinkPrincipalsList"

// PrivateLinkPrincipalsListLinkKind is the name of the type used to represent links to list
// of objects of type 'private_link_principals'.
const PrivateLinkPrincipalsListLinkKind = "PrivateLinkPrincipalsListLink"

// PrivateLinkPrincipalsNilKind is the name of the type used to nil lists of objects of
// type 'private_link_principals'.
const PrivateLinkPrincipalsListNilKind = "PrivateLinkPrincipalsListNil"

// PrivateLinkPrincipalsList is a list of values of the 'private_link_principals' type.
type PrivateLinkPrincipalsList struct {
	href  string
	link  bool
	items []*PrivateLinkPrincipals
}

// Kind returns the name of the type of the object.
func (l *PrivateLinkPrincipalsList) Kind() string {
	if l == nil {
		return PrivateLinkPrincipalsListNilKind
	}
	if l.link {
		return PrivateLinkPrincipalsListLinkKind
	}
	return PrivateLinkPrincipalsListKind
}

// Link returns true iif this is a link.
func (l *PrivateLinkPrincipalsList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *PrivateLinkPrincipalsList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *PrivateLinkPrincipalsList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *PrivateLinkPrincipalsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *PrivateLinkPrincipalsList) SetItems(items []*PrivateLinkPrincipals) {
	l.items = items
}

// Items returns the items of the list.
func (l *PrivateLinkPrincipalsList) Items() []*PrivateLinkPrincipals {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *PrivateLinkPrincipalsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *PrivateLinkPrincipalsList) Get(i int) *PrivateLinkPrincipals {
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
func (l *PrivateLinkPrincipalsList) Slice() []*PrivateLinkPrincipals {
	var slice []*PrivateLinkPrincipals
	if l == nil {
		slice = make([]*PrivateLinkPrincipals, 0)
	} else {
		slice = make([]*PrivateLinkPrincipals, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *PrivateLinkPrincipalsList) Each(f func(item *PrivateLinkPrincipals) bool) {
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
func (l *PrivateLinkPrincipalsList) Range(f func(index int, item *PrivateLinkPrincipals) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
