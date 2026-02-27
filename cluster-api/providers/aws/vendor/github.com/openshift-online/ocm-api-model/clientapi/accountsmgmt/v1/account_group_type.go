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

// AccountGroupKind is the name of the type used to represent objects
// of type 'account_group'.
const AccountGroupKind = "AccountGroup"

// AccountGroupLinkKind is the name of the type used to represent links
// to objects of type 'account_group'.
const AccountGroupLinkKind = "AccountGroupLink"

// AccountGroupNilKind is the name of the type used to nil references
// to objects of type 'account_group'.
const AccountGroupNilKind = "AccountGroupNil"

// AccountGroup represents the values of the 'account_group' type.
type AccountGroup struct {
	fieldSet_      []bool
	id             string
	href           string
	createdAt      time.Time
	description    string
	externalID     string
	managedBy      AccountGroupManagedBy
	name           string
	organizationID string
	updatedAt      time.Time
}

// Kind returns the name of the type of the object.
func (o *AccountGroup) Kind() string {
	if o == nil {
		return AccountGroupNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AccountGroupLinkKind
	}
	return AccountGroupKind
}

// Link returns true if this is a link.
func (o *AccountGroup) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AccountGroup) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AccountGroup) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AccountGroup) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AccountGroup) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccountGroup) Empty() bool {
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

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroup) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroup) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.createdAt
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroup) Description() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroup) GetDescription() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.description
	}
	return
}

// ExternalID returns the value of the 'external_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// An optional field for the Platform RBAC's identifier for this account group.
// Used when the group is managed by an external RBAC service to track the
// corresponding ID in that system.
func (o *AccountGroup) ExternalID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.externalID
	}
	return ""
}

// GetExternalID returns the value of the 'external_ID' attribute and
// a flag indicating if the attribute has a value.
//
// An optional field for the Platform RBAC's identifier for this account group.
// Used when the group is managed by an external RBAC service to track the
// corresponding ID in that system.
func (o *AccountGroup) GetExternalID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.externalID
	}
	return
}

// ManagedBy returns the value of the 'managed_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which system manages this account group
// Can be either:
// RBAC - Group managed by remote RBAC service, synchronized by job
// OCM - Group managed by OCM's APIs directly
func (o *AccountGroup) ManagedBy() AccountGroupManagedBy {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.managedBy
	}
	return AccountGroupManagedBy("")
}

// GetManagedBy returns the value of the 'managed_by' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which system manages this account group
// Can be either:
// RBAC - Group managed by remote RBAC service, synchronized by job
// OCM - Group managed by OCM's APIs directly
func (o *AccountGroup) GetManagedBy() (value AccountGroupManagedBy, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.managedBy
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroup) Name() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroup) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.name
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroup) OrganizationID() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroup) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.organizationID
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroup) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroup) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.updatedAt
	}
	return
}

// AccountGroupListKind is the name of the type used to represent list of objects of
// type 'account_group'.
const AccountGroupListKind = "AccountGroupList"

// AccountGroupListLinkKind is the name of the type used to represent links to list
// of objects of type 'account_group'.
const AccountGroupListLinkKind = "AccountGroupListLink"

// AccountGroupNilKind is the name of the type used to nil lists of objects of
// type 'account_group'.
const AccountGroupListNilKind = "AccountGroupListNil"

// AccountGroupList is a list of values of the 'account_group' type.
type AccountGroupList struct {
	href  string
	link  bool
	items []*AccountGroup
}

// Kind returns the name of the type of the object.
func (l *AccountGroupList) Kind() string {
	if l == nil {
		return AccountGroupListNilKind
	}
	if l.link {
		return AccountGroupListLinkKind
	}
	return AccountGroupListKind
}

// Link returns true iif this is a link.
func (l *AccountGroupList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AccountGroupList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AccountGroupList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AccountGroupList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AccountGroupList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AccountGroupList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AccountGroupList) SetItems(items []*AccountGroup) {
	l.items = items
}

// Items returns the items of the list.
func (l *AccountGroupList) Items() []*AccountGroup {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AccountGroupList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccountGroupList) Get(i int) *AccountGroup {
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
func (l *AccountGroupList) Slice() []*AccountGroup {
	var slice []*AccountGroup
	if l == nil {
		slice = make([]*AccountGroup, 0)
	} else {
		slice = make([]*AccountGroup, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccountGroupList) Each(f func(item *AccountGroup) bool) {
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
func (l *AccountGroupList) Range(f func(index int, item *AccountGroup) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
