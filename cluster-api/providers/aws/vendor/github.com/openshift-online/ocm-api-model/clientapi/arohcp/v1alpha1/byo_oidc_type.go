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

// ByoOidc represents the values of the 'byo_oidc' type.
//
// ByoOidc configuration.
type ByoOidc struct {
	fieldSet_ []bool
	enabled   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ByoOidc) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Boolean flag indicating if the cluster should be creating using _ByoOidc_.
//
// By default this is `false`.
//
// To enable it the cluster needs to be ROSA cluster and the organization of the user needs
// to have the `byo-oidc` feature toggle enabled.
func (o *ByoOidc) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Boolean flag indicating if the cluster should be creating using _ByoOidc_.
//
// By default this is `false`.
//
// To enable it the cluster needs to be ROSA cluster and the organization of the user needs
// to have the `byo-oidc` feature toggle enabled.
func (o *ByoOidc) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.enabled
	}
	return
}

// ByoOidcListKind is the name of the type used to represent list of objects of
// type 'byo_oidc'.
const ByoOidcListKind = "ByoOidcList"

// ByoOidcListLinkKind is the name of the type used to represent links to list
// of objects of type 'byo_oidc'.
const ByoOidcListLinkKind = "ByoOidcListLink"

// ByoOidcNilKind is the name of the type used to nil lists of objects of
// type 'byo_oidc'.
const ByoOidcListNilKind = "ByoOidcListNil"

// ByoOidcList is a list of values of the 'byo_oidc' type.
type ByoOidcList struct {
	href  string
	link  bool
	items []*ByoOidc
}

// Len returns the length of the list.
func (l *ByoOidcList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ByoOidcList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ByoOidcList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ByoOidcList) SetItems(items []*ByoOidc) {
	l.items = items
}

// Items returns the items of the list.
func (l *ByoOidcList) Items() []*ByoOidc {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ByoOidcList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ByoOidcList) Get(i int) *ByoOidc {
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
func (l *ByoOidcList) Slice() []*ByoOidc {
	var slice []*ByoOidc
	if l == nil {
		slice = make([]*ByoOidc, 0)
	} else {
		slice = make([]*ByoOidc, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ByoOidcList) Each(f func(item *ByoOidc) bool) {
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
func (l *ByoOidcList) Range(f func(index int, item *ByoOidc) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
