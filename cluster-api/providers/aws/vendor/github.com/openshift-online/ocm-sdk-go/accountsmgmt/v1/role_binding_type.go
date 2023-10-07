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

// RoleBindingKind is the name of the type used to represent objects
// of type 'role_binding'.
const RoleBindingKind = "RoleBinding"

// RoleBindingLinkKind is the name of the type used to represent links
// to objects of type 'role_binding'.
const RoleBindingLinkKind = "RoleBindingLink"

// RoleBindingNilKind is the name of the type used to nil references
// to objects of type 'role_binding'.
const RoleBindingNilKind = "RoleBindingNil"

// RoleBinding represents the values of the 'role_binding' type.
type RoleBinding struct {
	bitmap_        uint32
	id             string
	href           string
	account        *Account
	accountID      string
	createdAt      time.Time
	managedBy      string
	organization   *Organization
	organizationID string
	role           *Role
	roleID         string
	subscription   *Subscription
	subscriptionID string
	type_          string
	updatedAt      time.Time
	configManaged  bool
}

// Kind returns the name of the type of the object.
func (o *RoleBinding) Kind() string {
	if o == nil {
		return RoleBindingNilKind
	}
	if o.bitmap_&1 != 0 {
		return RoleBindingLinkKind
	}
	return RoleBindingKind
}

// Link returns true iif this is a link.
func (o *RoleBinding) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *RoleBinding) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *RoleBinding) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *RoleBinding) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *RoleBinding) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RoleBinding) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Account returns the value of the 'account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) Account() *Account {
	if o != nil && o.bitmap_&8 != 0 {
		return o.account
	}
	return nil
}

// GetAccount returns the value of the 'account' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetAccount() (value *Account, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.account
	}
	return
}

// AccountID returns the value of the 'account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) AccountID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.accountID
	}
	return ""
}

// GetAccountID returns the value of the 'account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.accountID
	}
	return
}

// ConfigManaged returns the value of the 'config_managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) ConfigManaged() bool {
	if o != nil && o.bitmap_&32 != 0 {
		return o.configManaged
	}
	return false
}

// GetConfigManaged returns the value of the 'config_managed' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetConfigManaged() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.configManaged
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) CreatedAt() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// ManagedBy returns the value of the 'managed_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) ManagedBy() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.managedBy
	}
	return ""
}

// GetManagedBy returns the value of the 'managed_by' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetManagedBy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.managedBy
	}
	return
}

// Organization returns the value of the 'organization' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) Organization() *Organization {
	if o != nil && o.bitmap_&256 != 0 {
		return o.organization
	}
	return nil
}

// GetOrganization returns the value of the 'organization' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetOrganization() (value *Organization, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.organization
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) OrganizationID() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.organizationID
	}
	return
}

// Role returns the value of the 'role' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) Role() *Role {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.role
	}
	return nil
}

// GetRole returns the value of the 'role' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetRole() (value *Role, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.role
	}
	return
}

// RoleID returns the value of the 'role_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) RoleID() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.roleID
	}
	return ""
}

// GetRoleID returns the value of the 'role_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetRoleID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.roleID
	}
	return
}

// Subscription returns the value of the 'subscription' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) Subscription() *Subscription {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.subscription
	}
	return nil
}

// GetSubscription returns the value of the 'subscription' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetSubscription() (value *Subscription, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.subscription
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) SubscriptionID() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.subscriptionID
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) Type() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.type_
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RoleBinding) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *RoleBinding) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// RoleBindingListKind is the name of the type used to represent list of objects of
// type 'role_binding'.
const RoleBindingListKind = "RoleBindingList"

// RoleBindingListLinkKind is the name of the type used to represent links to list
// of objects of type 'role_binding'.
const RoleBindingListLinkKind = "RoleBindingListLink"

// RoleBindingNilKind is the name of the type used to nil lists of objects of
// type 'role_binding'.
const RoleBindingListNilKind = "RoleBindingListNil"

// RoleBindingList is a list of values of the 'role_binding' type.
type RoleBindingList struct {
	href  string
	link  bool
	items []*RoleBinding
}

// Kind returns the name of the type of the object.
func (l *RoleBindingList) Kind() string {
	if l == nil {
		return RoleBindingListNilKind
	}
	if l.link {
		return RoleBindingListLinkKind
	}
	return RoleBindingListKind
}

// Link returns true iif this is a link.
func (l *RoleBindingList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *RoleBindingList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *RoleBindingList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *RoleBindingList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *RoleBindingList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RoleBindingList) Get(i int) *RoleBinding {
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
func (l *RoleBindingList) Slice() []*RoleBinding {
	var slice []*RoleBinding
	if l == nil {
		slice = make([]*RoleBinding, 0)
	} else {
		slice = make([]*RoleBinding, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RoleBindingList) Each(f func(item *RoleBinding) bool) {
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
func (l *RoleBindingList) Range(f func(index int, item *RoleBinding) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
