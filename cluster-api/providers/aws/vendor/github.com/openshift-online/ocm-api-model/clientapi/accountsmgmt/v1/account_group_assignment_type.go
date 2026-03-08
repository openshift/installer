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

// AccountGroupAssignmentKind is the name of the type used to represent objects
// of type 'account_group_assignment'.
const AccountGroupAssignmentKind = "AccountGroupAssignment"

// AccountGroupAssignmentLinkKind is the name of the type used to represent links
// to objects of type 'account_group_assignment'.
const AccountGroupAssignmentLinkKind = "AccountGroupAssignmentLink"

// AccountGroupAssignmentNilKind is the name of the type used to nil references
// to objects of type 'account_group_assignment'.
const AccountGroupAssignmentNilKind = "AccountGroupAssignmentNil"

// AccountGroupAssignment represents the values of the 'account_group_assignment' type.
type AccountGroupAssignment struct {
	fieldSet_      []bool
	id             string
	href           string
	accountID      string
	accountGroup   *AccountGroup
	accountGroupID string
	createdAt      time.Time
	managedBy      AccountGroupAssignmentManagedBy
	updatedAt      time.Time
}

// Kind returns the name of the type of the object.
func (o *AccountGroupAssignment) Kind() string {
	if o == nil {
		return AccountGroupAssignmentNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AccountGroupAssignmentLinkKind
	}
	return AccountGroupAssignmentKind
}

// Link returns true if this is a link.
func (o *AccountGroupAssignment) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AccountGroupAssignment) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AccountGroupAssignment) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AccountGroupAssignment) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AccountGroupAssignment) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccountGroupAssignment) Empty() bool {
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

// AccountID returns the value of the 'account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroupAssignment) AccountID() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.accountID
	}
	return ""
}

// GetAccountID returns the value of the 'account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroupAssignment) GetAccountID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.accountID
	}
	return
}

// AccountGroup returns the value of the 'account_group' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroupAssignment) AccountGroup() *AccountGroup {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.accountGroup
	}
	return nil
}

// GetAccountGroup returns the value of the 'account_group' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroupAssignment) GetAccountGroup() (value *AccountGroup, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.accountGroup
	}
	return
}

// AccountGroupID returns the value of the 'account_group_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroupAssignment) AccountGroupID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.accountGroupID
	}
	return ""
}

// GetAccountGroupID returns the value of the 'account_group_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroupAssignment) GetAccountGroupID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.accountGroupID
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroupAssignment) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroupAssignment) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.createdAt
	}
	return
}

// ManagedBy returns the value of the 'managed_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which system manages this account group assignment
// Can be either:
// RBAC - Assignment managed by remote RBAC service, synchronized by job
// OCM - Assignment managed by OCM's APIs directly
func (o *AccountGroupAssignment) ManagedBy() AccountGroupAssignmentManagedBy {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.managedBy
	}
	return AccountGroupAssignmentManagedBy("")
}

// GetManagedBy returns the value of the 'managed_by' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which system manages this account group assignment
// Can be either:
// RBAC - Assignment managed by remote RBAC service, synchronized by job
// OCM - Assignment managed by OCM's APIs directly
func (o *AccountGroupAssignment) GetManagedBy() (value AccountGroupAssignmentManagedBy, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.managedBy
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AccountGroupAssignment) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *AccountGroupAssignment) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.updatedAt
	}
	return
}

// AccountGroupAssignmentListKind is the name of the type used to represent list of objects of
// type 'account_group_assignment'.
const AccountGroupAssignmentListKind = "AccountGroupAssignmentList"

// AccountGroupAssignmentListLinkKind is the name of the type used to represent links to list
// of objects of type 'account_group_assignment'.
const AccountGroupAssignmentListLinkKind = "AccountGroupAssignmentListLink"

// AccountGroupAssignmentNilKind is the name of the type used to nil lists of objects of
// type 'account_group_assignment'.
const AccountGroupAssignmentListNilKind = "AccountGroupAssignmentListNil"

// AccountGroupAssignmentList is a list of values of the 'account_group_assignment' type.
type AccountGroupAssignmentList struct {
	href  string
	link  bool
	items []*AccountGroupAssignment
}

// Kind returns the name of the type of the object.
func (l *AccountGroupAssignmentList) Kind() string {
	if l == nil {
		return AccountGroupAssignmentListNilKind
	}
	if l.link {
		return AccountGroupAssignmentListLinkKind
	}
	return AccountGroupAssignmentListKind
}

// Link returns true iif this is a link.
func (l *AccountGroupAssignmentList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AccountGroupAssignmentList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AccountGroupAssignmentList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AccountGroupAssignmentList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AccountGroupAssignmentList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AccountGroupAssignmentList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AccountGroupAssignmentList) SetItems(items []*AccountGroupAssignment) {
	l.items = items
}

// Items returns the items of the list.
func (l *AccountGroupAssignmentList) Items() []*AccountGroupAssignment {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AccountGroupAssignmentList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccountGroupAssignmentList) Get(i int) *AccountGroupAssignment {
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
func (l *AccountGroupAssignmentList) Slice() []*AccountGroupAssignment {
	var slice []*AccountGroupAssignment
	if l == nil {
		slice = make([]*AccountGroupAssignment, 0)
	} else {
		slice = make([]*AccountGroupAssignment, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccountGroupAssignmentList) Each(f func(item *AccountGroupAssignment) bool) {
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
func (l *AccountGroupAssignmentList) Range(f func(index int, item *AccountGroupAssignment) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
