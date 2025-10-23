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

import (
	time "time"
)

// RegistryCredentialKind is the name of the type used to represent objects
// of type 'registry_credential'.
const RegistryCredentialKind = "RegistryCredential"

// RegistryCredentialLinkKind is the name of the type used to represent links
// to objects of type 'registry_credential'.
const RegistryCredentialLinkKind = "RegistryCredentialLink"

// RegistryCredentialNilKind is the name of the type used to nil references
// to objects of type 'registry_credential'.
const RegistryCredentialNilKind = "RegistryCredentialNil"

// RegistryCredential represents the values of the 'registry_credential' type.
type RegistryCredential struct {
	fieldSet_          []bool
	id                 string
	href               string
	account            *Account
	createdAt          time.Time
	externalResourceID string
	registry           *Registry
	token              string
	updatedAt          time.Time
	username           string
}

// Kind returns the name of the type of the object.
func (o *RegistryCredential) Kind() string {
	if o == nil {
		return RegistryCredentialNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return RegistryCredentialLinkKind
	}
	return RegistryCredentialKind
}

// Link returns true if this is a link.
func (o *RegistryCredential) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *RegistryCredential) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *RegistryCredential) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *RegistryCredential) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *RegistryCredential) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RegistryCredential) Empty() bool {
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

// Account returns the value of the 'account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) Account() *Account {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.account
	}
	return nil
}

// GetAccount returns the value of the 'account' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetAccount() (value *Account, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.account
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.createdAt
	}
	return
}

// ExternalResourceID returns the value of the 'external_resource_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) ExternalResourceID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.externalResourceID
	}
	return ""
}

// GetExternalResourceID returns the value of the 'external_resource_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetExternalResourceID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.externalResourceID
	}
	return
}

// Registry returns the value of the 'registry' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) Registry() *Registry {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.registry
	}
	return nil
}

// GetRegistry returns the value of the 'registry' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetRegistry() (value *Registry, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.registry
	}
	return
}

// Token returns the value of the 'token' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) Token() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.token
	}
	return ""
}

// GetToken returns the value of the 'token' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetToken() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.token
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.updatedAt
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RegistryCredential) Username() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
func (o *RegistryCredential) GetUsername() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.username
	}
	return
}

// RegistryCredentialListKind is the name of the type used to represent list of objects of
// type 'registry_credential'.
const RegistryCredentialListKind = "RegistryCredentialList"

// RegistryCredentialListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry_credential'.
const RegistryCredentialListLinkKind = "RegistryCredentialListLink"

// RegistryCredentialNilKind is the name of the type used to nil lists of objects of
// type 'registry_credential'.
const RegistryCredentialListNilKind = "RegistryCredentialListNil"

// RegistryCredentialList is a list of values of the 'registry_credential' type.
type RegistryCredentialList struct {
	href  string
	link  bool
	items []*RegistryCredential
}

// Kind returns the name of the type of the object.
func (l *RegistryCredentialList) Kind() string {
	if l == nil {
		return RegistryCredentialListNilKind
	}
	if l.link {
		return RegistryCredentialListLinkKind
	}
	return RegistryCredentialListKind
}

// Link returns true iif this is a link.
func (l *RegistryCredentialList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *RegistryCredentialList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *RegistryCredentialList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *RegistryCredentialList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *RegistryCredentialList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *RegistryCredentialList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *RegistryCredentialList) SetItems(items []*RegistryCredential) {
	l.items = items
}

// Items returns the items of the list.
func (l *RegistryCredentialList) Items() []*RegistryCredential {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *RegistryCredentialList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RegistryCredentialList) Get(i int) *RegistryCredential {
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
func (l *RegistryCredentialList) Slice() []*RegistryCredential {
	var slice []*RegistryCredential
	if l == nil {
		slice = make([]*RegistryCredential, 0)
	} else {
		slice = make([]*RegistryCredential, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RegistryCredentialList) Each(f func(item *RegistryCredential) bool) {
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
func (l *RegistryCredentialList) Range(f func(index int, item *RegistryCredential) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
