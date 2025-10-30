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

// AdminCredentials represents the values of the 'admin_credentials' type.
//
// Temporary administrator credentials generated during the installation of the
// cluster.
type AdminCredentials struct {
	bitmap_  uint32
	password string
	user     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AdminCredentials) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Password returns the value of the 'password' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster administrator password.
func (o *AdminCredentials) Password() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.password
	}
	return ""
}

// GetPassword returns the value of the 'password' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster administrator password.
func (o *AdminCredentials) GetPassword() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.password
	}
	return
}

// User returns the value of the 'user' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster administrator user name.
func (o *AdminCredentials) User() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.user
	}
	return ""
}

// GetUser returns the value of the 'user' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster administrator user name.
func (o *AdminCredentials) GetUser() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.user
	}
	return
}

// AdminCredentialsListKind is the name of the type used to represent list of objects of
// type 'admin_credentials'.
const AdminCredentialsListKind = "AdminCredentialsList"

// AdminCredentialsListLinkKind is the name of the type used to represent links to list
// of objects of type 'admin_credentials'.
const AdminCredentialsListLinkKind = "AdminCredentialsListLink"

// AdminCredentialsNilKind is the name of the type used to nil lists of objects of
// type 'admin_credentials'.
const AdminCredentialsListNilKind = "AdminCredentialsListNil"

// AdminCredentialsList is a list of values of the 'admin_credentials' type.
type AdminCredentialsList struct {
	href  string
	link  bool
	items []*AdminCredentials
}

// Len returns the length of the list.
func (l *AdminCredentialsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AdminCredentialsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AdminCredentialsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AdminCredentialsList) SetItems(items []*AdminCredentials) {
	l.items = items
}

// Items returns the items of the list.
func (l *AdminCredentialsList) Items() []*AdminCredentials {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AdminCredentialsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AdminCredentialsList) Get(i int) *AdminCredentials {
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
func (l *AdminCredentialsList) Slice() []*AdminCredentials {
	var slice []*AdminCredentials
	if l == nil {
		slice = make([]*AdminCredentials, 0)
	} else {
		slice = make([]*AdminCredentials, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AdminCredentialsList) Each(f func(item *AdminCredentials) bool) {
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
func (l *AdminCredentialsList) Range(f func(index int, item *AdminCredentials) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
