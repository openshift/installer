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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// WifSupport represents the values of the 'wif_support' type.
type WifSupport struct {
	fieldSet_ []bool
	principal string
	roles     []*WifRole
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifSupport) Empty() bool {
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

// Principal returns the value of the 'principal' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifSupport) Principal() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.principal
	}
	return ""
}

// GetPrincipal returns the value of the 'principal' attribute and
// a flag indicating if the attribute has a value.
func (o *WifSupport) GetPrincipal() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.principal
	}
	return
}

// Roles returns the value of the 'roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifSupport) Roles() []*WifRole {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.roles
	}
	return nil
}

// GetRoles returns the value of the 'roles' attribute and
// a flag indicating if the attribute has a value.
func (o *WifSupport) GetRoles() (value []*WifRole, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.roles
	}
	return
}

// WifSupportListKind is the name of the type used to represent list of objects of
// type 'wif_support'.
const WifSupportListKind = "WifSupportList"

// WifSupportListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_support'.
const WifSupportListLinkKind = "WifSupportListLink"

// WifSupportNilKind is the name of the type used to nil lists of objects of
// type 'wif_support'.
const WifSupportListNilKind = "WifSupportListNil"

// WifSupportList is a list of values of the 'wif_support' type.
type WifSupportList struct {
	href  string
	link  bool
	items []*WifSupport
}

// Len returns the length of the list.
func (l *WifSupportList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *WifSupportList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *WifSupportList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *WifSupportList) SetItems(items []*WifSupport) {
	l.items = items
}

// Items returns the items of the list.
func (l *WifSupportList) Items() []*WifSupport {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *WifSupportList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifSupportList) Get(i int) *WifSupport {
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
func (l *WifSupportList) Slice() []*WifSupport {
	var slice []*WifSupport
	if l == nil {
		slice = make([]*WifSupport, 0)
	} else {
		slice = make([]*WifSupport, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifSupportList) Each(f func(item *WifSupport) bool) {
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
func (l *WifSupportList) Range(f func(index int, item *WifSupport) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
