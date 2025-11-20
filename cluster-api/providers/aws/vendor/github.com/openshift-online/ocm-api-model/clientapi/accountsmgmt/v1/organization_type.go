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

// OrganizationKind is the name of the type used to represent objects
// of type 'organization'.
const OrganizationKind = "Organization"

// OrganizationLinkKind is the name of the type used to represent links
// to objects of type 'organization'.
const OrganizationLinkKind = "OrganizationLink"

// OrganizationNilKind is the name of the type used to nil references
// to objects of type 'organization'.
const OrganizationNilKind = "OrganizationNil"

// Organization represents the values of the 'organization' type.
type Organization struct {
	fieldSet_    []bool
	id           string
	href         string
	capabilities []*Capability
	createdAt    time.Time
	ebsAccountID string
	externalID   string
	labels       []*Label
	name         string
	updatedAt    time.Time
}

// Kind returns the name of the type of the object.
func (o *Organization) Kind() string {
	if o == nil {
		return OrganizationNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return OrganizationLinkKind
	}
	return OrganizationKind
}

// Link returns true if this is a link.
func (o *Organization) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Organization) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Organization) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Organization) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Organization) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Organization) Empty() bool {
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

// Capabilities returns the value of the 'capabilities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) Capabilities() []*Capability {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.capabilities
	}
	return nil
}

// GetCapabilities returns the value of the 'capabilities' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetCapabilities() (value []*Capability, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.capabilities
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.createdAt
	}
	return
}

// EbsAccountID returns the value of the 'ebs_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) EbsAccountID() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.ebsAccountID
	}
	return ""
}

// GetEbsAccountID returns the value of the 'ebs_account_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetEbsAccountID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.ebsAccountID
	}
	return
}

// ExternalID returns the value of the 'external_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) ExternalID() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.externalID
	}
	return ""
}

// GetExternalID returns the value of the 'external_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetExternalID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.externalID
	}
	return
}

// Labels returns the value of the 'labels' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) Labels() []*Label {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.labels
	}
	return nil
}

// GetLabels returns the value of the 'labels' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetLabels() (value []*Label, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.labels
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) Name() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.name
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Organization) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Organization) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.updatedAt
	}
	return
}

// OrganizationListKind is the name of the type used to represent list of objects of
// type 'organization'.
const OrganizationListKind = "OrganizationList"

// OrganizationListLinkKind is the name of the type used to represent links to list
// of objects of type 'organization'.
const OrganizationListLinkKind = "OrganizationListLink"

// OrganizationNilKind is the name of the type used to nil lists of objects of
// type 'organization'.
const OrganizationListNilKind = "OrganizationListNil"

// OrganizationList is a list of values of the 'organization' type.
type OrganizationList struct {
	href  string
	link  bool
	items []*Organization
}

// Kind returns the name of the type of the object.
func (l *OrganizationList) Kind() string {
	if l == nil {
		return OrganizationListNilKind
	}
	if l.link {
		return OrganizationListLinkKind
	}
	return OrganizationListKind
}

// Link returns true iif this is a link.
func (l *OrganizationList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *OrganizationList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *OrganizationList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *OrganizationList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *OrganizationList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *OrganizationList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *OrganizationList) SetItems(items []*Organization) {
	l.items = items
}

// Items returns the items of the list.
func (l *OrganizationList) Items() []*Organization {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *OrganizationList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OrganizationList) Get(i int) *Organization {
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
func (l *OrganizationList) Slice() []*Organization {
	var slice []*Organization
	if l == nil {
		slice = make([]*Organization, 0)
	} else {
		slice = make([]*Organization, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OrganizationList) Each(f func(item *Organization) bool) {
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
func (l *OrganizationList) Range(f func(index int, item *Organization) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
