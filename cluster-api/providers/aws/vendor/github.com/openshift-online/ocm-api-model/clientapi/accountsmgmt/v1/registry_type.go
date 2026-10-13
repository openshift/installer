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

// RegistryKind is the name of the type used to represent objects
// of type 'registry'.
const RegistryKind = "Registry"

// RegistryLinkKind is the name of the type used to represent links
// to objects of type 'registry'.
const RegistryLinkKind = "RegistryLink"

// RegistryNilKind is the name of the type used to nil references
// to objects of type 'registry'.
const RegistryNilKind = "RegistryNil"

// Registry represents the values of the 'registry' type.
type Registry struct {
	fieldSet_  []bool
	id         string
	href       string
	url        string
	createdAt  time.Time
	name       string
	orgName    string
	teamName   string
	type_      string
	updatedAt  time.Time
	cloudAlias bool
}

// Kind returns the name of the type of the object.
func (o *Registry) Kind() string {
	if o == nil {
		return RegistryNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return RegistryLinkKind
	}
	return RegistryKind
}

// Link returns true if this is a link.
func (o *Registry) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Registry) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Registry) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Registry) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Registry) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Registry) Empty() bool {
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

// URL returns the value of the 'URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) URL() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.url
	}
	return ""
}

// GetURL returns the value of the 'URL' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetURL() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.url
	}
	return
}

// CloudAlias returns the value of the 'cloud_alias' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) CloudAlias() bool {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.cloudAlias
	}
	return false
}

// GetCloudAlias returns the value of the 'cloud_alias' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetCloudAlias() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.cloudAlias
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.createdAt
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) Name() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.name
	}
	return
}

// OrgName returns the value of the 'org_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) OrgName() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.orgName
	}
	return ""
}

// GetOrgName returns the value of the 'org_name' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetOrgName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.orgName
	}
	return
}

// TeamName returns the value of the 'team_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) TeamName() string {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.teamName
	}
	return ""
}

// GetTeamName returns the value of the 'team_name' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetTeamName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.teamName
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) Type() string {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.type_
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Registry) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Registry) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.updatedAt
	}
	return
}

// RegistryListKind is the name of the type used to represent list of objects of
// type 'registry'.
const RegistryListKind = "RegistryList"

// RegistryListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry'.
const RegistryListLinkKind = "RegistryListLink"

// RegistryNilKind is the name of the type used to nil lists of objects of
// type 'registry'.
const RegistryListNilKind = "RegistryListNil"

// RegistryList is a list of values of the 'registry' type.
type RegistryList struct {
	href  string
	link  bool
	items []*Registry
}

// Kind returns the name of the type of the object.
func (l *RegistryList) Kind() string {
	if l == nil {
		return RegistryListNilKind
	}
	if l.link {
		return RegistryListLinkKind
	}
	return RegistryListKind
}

// Link returns true iif this is a link.
func (l *RegistryList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *RegistryList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *RegistryList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *RegistryList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *RegistryList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *RegistryList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *RegistryList) SetItems(items []*Registry) {
	l.items = items
}

// Items returns the items of the list.
func (l *RegistryList) Items() []*Registry {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *RegistryList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RegistryList) Get(i int) *Registry {
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
func (l *RegistryList) Slice() []*Registry {
	var slice []*Registry
	if l == nil {
		slice = make([]*Registry, 0)
	} else {
		slice = make([]*Registry, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RegistryList) Each(f func(item *Registry) bool) {
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
func (l *RegistryList) Range(f func(index int, item *Registry) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
