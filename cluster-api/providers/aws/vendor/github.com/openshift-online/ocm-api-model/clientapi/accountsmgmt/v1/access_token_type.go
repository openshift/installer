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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// AccessToken represents the values of the 'access_token' type.
type AccessToken struct {
	fieldSet_ []bool
	auths     map[string]*AccessTokenAuth
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccessToken) Empty() bool {
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

// Auths returns the value of the 'auths' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccessToken) Auths() map[string]*AccessTokenAuth {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.auths
	}
	return nil
}

// GetAuths returns the value of the 'auths' attribute and
// a flag indicating if the attribute has a value.
func (o *AccessToken) GetAuths() (value map[string]*AccessTokenAuth, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.auths
	}
	return
}

// AccessTokenListKind is the name of the type used to represent list of objects of
// type 'access_token'.
const AccessTokenListKind = "AccessTokenList"

// AccessTokenListLinkKind is the name of the type used to represent links to list
// of objects of type 'access_token'.
const AccessTokenListLinkKind = "AccessTokenListLink"

// AccessTokenNilKind is the name of the type used to nil lists of objects of
// type 'access_token'.
const AccessTokenListNilKind = "AccessTokenListNil"

// AccessTokenList is a list of values of the 'access_token' type.
type AccessTokenList struct {
	href  string
	link  bool
	items []*AccessToken
}

// Len returns the length of the list.
func (l *AccessTokenList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AccessTokenList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AccessTokenList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AccessTokenList) SetItems(items []*AccessToken) {
	l.items = items
}

// Items returns the items of the list.
func (l *AccessTokenList) Items() []*AccessToken {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AccessTokenList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccessTokenList) Get(i int) *AccessToken {
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
func (l *AccessTokenList) Slice() []*AccessToken {
	var slice []*AccessToken
	if l == nil {
		slice = make([]*AccessToken, 0)
	} else {
		slice = make([]*AccessToken, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccessTokenList) Each(f func(item *AccessToken) bool) {
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
func (l *AccessTokenList) Range(f func(index int, item *AccessToken) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
