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

// OpenIDClaims represents the values of the 'open_ID_claims' type.
//
// _OpenID_ identity provider claims.
type OpenIDClaims struct {
	fieldSet_         []bool
	email             []string
	groups            []string
	name              []string
	preferredUsername []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OpenIDClaims) Empty() bool {
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

// Email returns the value of the 'email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of claims to use as the mail address.
func (o *OpenIDClaims) Email() []string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.email
	}
	return nil
}

// GetEmail returns the value of the 'email' attribute and
// a flag indicating if the attribute has a value.
//
// List of claims to use as the mail address.
func (o *OpenIDClaims) GetEmail() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.email
	}
	return
}

// Groups returns the value of the 'groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of claims to use as the group name.
func (o *OpenIDClaims) Groups() []string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.groups
	}
	return nil
}

// GetGroups returns the value of the 'groups' attribute and
// a flag indicating if the attribute has a value.
//
// List of claims to use as the group name.
func (o *OpenIDClaims) GetGroups() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.groups
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of claims to use as the display name.
func (o *OpenIDClaims) Name() []string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.name
	}
	return nil
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// List of claims to use as the display name.
func (o *OpenIDClaims) GetName() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.name
	}
	return
}

// PreferredUsername returns the value of the 'preferred_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of claims to use as the preferred user name when provisioning a user.
func (o *OpenIDClaims) PreferredUsername() []string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.preferredUsername
	}
	return nil
}

// GetPreferredUsername returns the value of the 'preferred_username' attribute and
// a flag indicating if the attribute has a value.
//
// List of claims to use as the preferred user name when provisioning a user.
func (o *OpenIDClaims) GetPreferredUsername() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.preferredUsername
	}
	return
}

// OpenIDClaimsListKind is the name of the type used to represent list of objects of
// type 'open_ID_claims'.
const OpenIDClaimsListKind = "OpenIDClaimsList"

// OpenIDClaimsListLinkKind is the name of the type used to represent links to list
// of objects of type 'open_ID_claims'.
const OpenIDClaimsListLinkKind = "OpenIDClaimsListLink"

// OpenIDClaimsNilKind is the name of the type used to nil lists of objects of
// type 'open_ID_claims'.
const OpenIDClaimsListNilKind = "OpenIDClaimsListNil"

// OpenIDClaimsList is a list of values of the 'open_ID_claims' type.
type OpenIDClaimsList struct {
	href  string
	link  bool
	items []*OpenIDClaims
}

// Len returns the length of the list.
func (l *OpenIDClaimsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *OpenIDClaimsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *OpenIDClaimsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *OpenIDClaimsList) SetItems(items []*OpenIDClaims) {
	l.items = items
}

// Items returns the items of the list.
func (l *OpenIDClaimsList) Items() []*OpenIDClaims {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *OpenIDClaimsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OpenIDClaimsList) Get(i int) *OpenIDClaims {
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
func (l *OpenIDClaimsList) Slice() []*OpenIDClaims {
	var slice []*OpenIDClaims
	if l == nil {
		slice = make([]*OpenIDClaims, 0)
	} else {
		slice = make([]*OpenIDClaims, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OpenIDClaimsList) Each(f func(item *OpenIDClaims) bool) {
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
func (l *OpenIDClaimsList) Range(f func(index int, item *OpenIDClaims) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
