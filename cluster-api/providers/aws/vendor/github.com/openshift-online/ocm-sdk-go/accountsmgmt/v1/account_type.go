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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	time "time"
)

// AccountKind is the name of the type used to represent objects
// of type 'account'.
const AccountKind = "Account"

// AccountLinkKind is the name of the type used to represent links
// to objects of type 'account'.
const AccountLinkKind = "AccountLink"

// AccountNilKind is the name of the type used to nil references
// to objects of type 'account'.
const AccountNilKind = "AccountNil"

// Account represents the values of the 'account' type.
type Account struct {
	bitmap_        uint32
	id             string
	href           string
	banCode        string
	banDescription string
	capabilities   []*Capability
	createdAt      time.Time
	email          string
	firstName      string
	labels         []*Label
	lastName       string
	organization   *Organization
	rhitAccountID  string
	rhitWebUserId  string
	updatedAt      time.Time
	username       string
	banned         bool
	serviceAccount bool
}

// Kind returns the name of the type of the object.
func (o *Account) Kind() string {
	if o == nil {
		return AccountNilKind
	}
	if o.bitmap_&1 != 0 {
		return AccountLinkKind
	}
	return AccountKind
}

// Link returns true iif this is a link.
func (o *Account) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Account) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Account) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Account) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Account) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Account) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// BanCode returns the value of the 'ban_code' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) BanCode() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.banCode
	}
	return ""
}

// GetBanCode returns the value of the 'ban_code' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetBanCode() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.banCode
	}
	return
}

// BanDescription returns the value of the 'ban_description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) BanDescription() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.banDescription
	}
	return ""
}

// GetBanDescription returns the value of the 'ban_description' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetBanDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.banDescription
	}
	return
}

// Banned returns the value of the 'banned' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Banned() bool {
	if o != nil && o.bitmap_&32 != 0 {
		return o.banned
	}
	return false
}

// GetBanned returns the value of the 'banned' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetBanned() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.banned
	}
	return
}

// Capabilities returns the value of the 'capabilities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Capabilities() []*Capability {
	if o != nil && o.bitmap_&64 != 0 {
		return o.capabilities
	}
	return nil
}

// GetCapabilities returns the value of the 'capabilities' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetCapabilities() (value []*Capability, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.capabilities
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) CreatedAt() time.Time {
	if o != nil && o.bitmap_&128 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Email returns the value of the 'email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Email() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.email
	}
	return ""
}

// GetEmail returns the value of the 'email' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetEmail() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.email
	}
	return
}

// FirstName returns the value of the 'first_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) FirstName() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.firstName
	}
	return ""
}

// GetFirstName returns the value of the 'first_name' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetFirstName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.firstName
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Labels() []*Label {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetLabels() (value []*Label, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.labels
	}
	return
}

// LastName returns the value of the 'last_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) LastName() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.lastName
	}
	return ""
}

// GetLastName returns the value of the 'last_name' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetLastName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.lastName
	}
	return
}

// Organization returns the value of the 'organization' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Organization() *Organization {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.organization
	}
	return nil
}

// GetOrganization returns the value of the 'organization' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetOrganization() (value *Organization, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.organization
	}
	return
}

// RhitAccountID returns the value of the 'rhit_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// RhitAccountID will be deprecated in favor of RhitWebUserId
func (o *Account) RhitAccountID() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.rhitAccountID
	}
	return ""
}

// GetRhitAccountID returns the value of the 'rhit_account_ID' attribute and
// a flag indicating if the attribute has a value.
//
// RhitAccountID will be deprecated in favor of RhitWebUserId
func (o *Account) GetRhitAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.rhitAccountID
	}
	return
}

// RhitWebUserId returns the value of the 'rhit_web_user_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) RhitWebUserId() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.rhitWebUserId
	}
	return ""
}

// GetRhitWebUserId returns the value of the 'rhit_web_user_id' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetRhitWebUserId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.rhitWebUserId
	}
	return
}

// ServiceAccount returns the value of the 'service_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) ServiceAccount() bool {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.serviceAccount
	}
	return false
}

// GetServiceAccount returns the value of the 'service_account' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetServiceAccount() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.serviceAccount
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Account) Username() string {
	if o != nil && o.bitmap_&131072 != 0 {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
func (o *Account) GetUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&131072 != 0
	if ok {
		value = o.username
	}
	return
}

// AccountListKind is the name of the type used to represent list of objects of
// type 'account'.
const AccountListKind = "AccountList"

// AccountListLinkKind is the name of the type used to represent links to list
// of objects of type 'account'.
const AccountListLinkKind = "AccountListLink"

// AccountNilKind is the name of the type used to nil lists of objects of
// type 'account'.
const AccountListNilKind = "AccountListNil"

// AccountList is a list of values of the 'account' type.
type AccountList struct {
	href  string
	link  bool
	items []*Account
}

// Kind returns the name of the type of the object.
func (l *AccountList) Kind() string {
	if l == nil {
		return AccountListNilKind
	}
	if l.link {
		return AccountListLinkKind
	}
	return AccountListKind
}

// Link returns true iif this is a link.
func (l *AccountList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AccountList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AccountList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AccountList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AccountList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccountList) Get(i int) *Account {
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
func (l *AccountList) Slice() []*Account {
	var slice []*Account
	if l == nil {
		slice = make([]*Account, 0)
	} else {
		slice = make([]*Account, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccountList) Each(f func(item *Account) bool) {
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
func (l *AccountList) Range(f func(index int, item *Account) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
