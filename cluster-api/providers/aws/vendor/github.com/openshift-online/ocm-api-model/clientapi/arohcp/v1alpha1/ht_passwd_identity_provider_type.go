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

// HTPasswdIdentityProvider represents the values of the 'HT_passwd_identity_provider' type.
//
// Details for `htpasswd` identity providers.
type HTPasswdIdentityProvider struct {
	fieldSet_ []bool
	password  string
	username  string
	users     *HTPasswdUserList
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *HTPasswdIdentityProvider) Empty() bool {
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

// Password returns the value of the 'password' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Password to be used in the _HTPasswd_ data file.
func (o *HTPasswdIdentityProvider) Password() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.password
	}
	return ""
}

// GetPassword returns the value of the 'password' attribute and
// a flag indicating if the attribute has a value.
//
// Password to be used in the _HTPasswd_ data file.
func (o *HTPasswdIdentityProvider) GetPassword() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.password
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Username to be used in the _HTPasswd_ data file.
func (o *HTPasswdIdentityProvider) Username() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
//
// Username to be used in the _HTPasswd_ data file.
func (o *HTPasswdIdentityProvider) GetUsername() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.username
	}
	return
}

// Users returns the value of the 'users' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the collection of _HTPasswd_ users.
func (o *HTPasswdIdentityProvider) Users() *HTPasswdUserList {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.users
	}
	return nil
}

// GetUsers returns the value of the 'users' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the collection of _HTPasswd_ users.
func (o *HTPasswdIdentityProvider) GetUsers() (value *HTPasswdUserList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.users
	}
	return
}

// HTPasswdIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'HT_passwd_identity_provider'.
const HTPasswdIdentityProviderListKind = "HTPasswdIdentityProviderList"

// HTPasswdIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'HT_passwd_identity_provider'.
const HTPasswdIdentityProviderListLinkKind = "HTPasswdIdentityProviderListLink"

// HTPasswdIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'HT_passwd_identity_provider'.
const HTPasswdIdentityProviderListNilKind = "HTPasswdIdentityProviderListNil"

// HTPasswdIdentityProviderList is a list of values of the 'HT_passwd_identity_provider' type.
type HTPasswdIdentityProviderList struct {
	href  string
	link  bool
	items []*HTPasswdIdentityProvider
}

// Len returns the length of the list.
func (l *HTPasswdIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *HTPasswdIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *HTPasswdIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *HTPasswdIdentityProviderList) SetItems(items []*HTPasswdIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *HTPasswdIdentityProviderList) Items() []*HTPasswdIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *HTPasswdIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *HTPasswdIdentityProviderList) Get(i int) *HTPasswdIdentityProvider {
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
func (l *HTPasswdIdentityProviderList) Slice() []*HTPasswdIdentityProvider {
	var slice []*HTPasswdIdentityProvider
	if l == nil {
		slice = make([]*HTPasswdIdentityProvider, 0)
	} else {
		slice = make([]*HTPasswdIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *HTPasswdIdentityProviderList) Each(f func(item *HTPasswdIdentityProvider) bool) {
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
func (l *HTPasswdIdentityProviderList) Range(f func(index int, item *HTPasswdIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
