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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	time "time"
)

// UserKind is the name of the type used to represent objects
// of type 'user'.
const UserKind = "User"

// UserLinkKind is the name of the type used to represent links
// to objects of type 'user'.
const UserLinkKind = "UserLink"

// UserNilKind is the name of the type used to nil references
// to objects of type 'user'.
const UserNilKind = "UserNil"

// User represents the values of the 'user' type.
//
// Definition of a Web RCA user.
type User struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	email     string
	name      string
	updatedAt time.Time
	username  string
	fromAuth  bool
}

// Kind returns the name of the type of the object.
func (o *User) Kind() string {
	if o == nil {
		return UserNilKind
	}
	if o.bitmap_&1 != 0 {
		return UserLinkKind
	}
	return UserKind
}

// Link returns true iif this is a link.
func (o *User) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *User) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *User) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *User) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *User) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *User) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *User) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *User) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *User) DeletedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *User) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// Email returns the value of the 'email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *User) Email() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.email
	}
	return ""
}

// GetEmail returns the value of the 'email' attribute and
// a flag indicating if the attribute has a value.
func (o *User) GetEmail() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.email
	}
	return
}

// FromAuth returns the value of the 'from_auth' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *User) FromAuth() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.fromAuth
	}
	return false
}

// GetFromAuth returns the value of the 'from_auth' attribute and
// a flag indicating if the attribute has a value.
func (o *User) GetFromAuth() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.fromAuth
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *User) Name() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *User) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.name
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *User) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *User) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *User) Username() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
func (o *User) GetUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.username
	}
	return
}

// UserListKind is the name of the type used to represent list of objects of
// type 'user'.
const UserListKind = "UserList"

// UserListLinkKind is the name of the type used to represent links to list
// of objects of type 'user'.
const UserListLinkKind = "UserListLink"

// UserNilKind is the name of the type used to nil lists of objects of
// type 'user'.
const UserListNilKind = "UserListNil"

// UserList is a list of values of the 'user' type.
type UserList struct {
	href  string
	link  bool
	items []*User
}

// Kind returns the name of the type of the object.
func (l *UserList) Kind() string {
	if l == nil {
		return UserListNilKind
	}
	if l.link {
		return UserListLinkKind
	}
	return UserListKind
}

// Link returns true iif this is a link.
func (l *UserList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *UserList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *UserList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *UserList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *UserList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *UserList) Get(i int) *User {
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
func (l *UserList) Slice() []*User {
	var slice []*User
	if l == nil {
		slice = make([]*User, 0)
	} else {
		slice = make([]*User, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *UserList) Each(f func(item *User) bool) {
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
func (l *UserList) Range(f func(index int, item *User) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
