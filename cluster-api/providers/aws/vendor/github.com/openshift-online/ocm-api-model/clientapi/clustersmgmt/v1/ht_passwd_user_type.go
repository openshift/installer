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

// HTPasswdUser represents the values of the 'HT_passwd_user' type.
type HTPasswdUser struct {
	fieldSet_      []bool
	id             string
	hashedPassword string
	password       string
	username       string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *HTPasswdUser) Empty() bool {
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

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID for a secondary user in the _HTPasswd_ data file.
func (o *HTPasswdUser) ID() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID for a secondary user in the _HTPasswd_ data file.
func (o *HTPasswdUser) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// HashedPassword returns the value of the 'hashed_password' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HTPasswd Hashed Password for a user in the _HTPasswd_ data file.
// The value of this field is set as-is in the _HTPasswd_ data file for the HTPasswd IDP
func (o *HTPasswdUser) HashedPassword() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.hashedPassword
	}
	return ""
}

// GetHashedPassword returns the value of the 'hashed_password' attribute and
// a flag indicating if the attribute has a value.
//
// HTPasswd Hashed Password for a user in the _HTPasswd_ data file.
// The value of this field is set as-is in the _HTPasswd_ data file for the HTPasswd IDP
func (o *HTPasswdUser) GetHashedPassword() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.hashedPassword
	}
	return
}

// Password returns the value of the 'password' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Password in plain-text for a  user in the _HTPasswd_ data file.
// The value of this field is hashed before setting it in the  _HTPasswd_ data file for the HTPasswd IDP
func (o *HTPasswdUser) Password() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.password
	}
	return ""
}

// GetPassword returns the value of the 'password' attribute and
// a flag indicating if the attribute has a value.
//
// Password in plain-text for a  user in the _HTPasswd_ data file.
// The value of this field is hashed before setting it in the  _HTPasswd_ data file for the HTPasswd IDP
func (o *HTPasswdUser) GetPassword() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.password
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Username for a secondary user in the _HTPasswd_ data file.
func (o *HTPasswdUser) Username() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
//
// Username for a secondary user in the _HTPasswd_ data file.
func (o *HTPasswdUser) GetUsername() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.username
	}
	return
}

// HTPasswdUserListKind is the name of the type used to represent list of objects of
// type 'HT_passwd_user'.
const HTPasswdUserListKind = "HTPasswdUserList"

// HTPasswdUserListLinkKind is the name of the type used to represent links to list
// of objects of type 'HT_passwd_user'.
const HTPasswdUserListLinkKind = "HTPasswdUserListLink"

// HTPasswdUserNilKind is the name of the type used to nil lists of objects of
// type 'HT_passwd_user'.
const HTPasswdUserListNilKind = "HTPasswdUserListNil"

// HTPasswdUserList is a list of values of the 'HT_passwd_user' type.
type HTPasswdUserList struct {
	href  string
	link  bool
	items []*HTPasswdUser
}

// Len returns the length of the list.
func (l *HTPasswdUserList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *HTPasswdUserList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *HTPasswdUserList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *HTPasswdUserList) SetItems(items []*HTPasswdUser) {
	l.items = items
}

// Items returns the items of the list.
func (l *HTPasswdUserList) Items() []*HTPasswdUser {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *HTPasswdUserList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *HTPasswdUserList) Get(i int) *HTPasswdUser {
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
func (l *HTPasswdUserList) Slice() []*HTPasswdUser {
	var slice []*HTPasswdUser
	if l == nil {
		slice = make([]*HTPasswdUser, 0)
	} else {
		slice = make([]*HTPasswdUser, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *HTPasswdUserList) Each(f func(item *HTPasswdUser) bool) {
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
func (l *HTPasswdUserList) Range(f func(index int, item *HTPasswdUser) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
