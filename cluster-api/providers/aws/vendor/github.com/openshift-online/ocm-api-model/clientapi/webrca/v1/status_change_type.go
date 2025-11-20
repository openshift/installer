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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/webrca/v1

import (
	time "time"
)

// StatusChangeKind is the name of the type used to represent objects
// of type 'status_change'.
const StatusChangeKind = "StatusChange"

// StatusChangeLinkKind is the name of the type used to represent links
// to objects of type 'status_change'.
const StatusChangeLinkKind = "StatusChangeLink"

// StatusChangeNilKind is the name of the type used to nil references
// to objects of type 'status_change'.
const StatusChangeNilKind = "StatusChangeNil"

// StatusChange represents the values of the 'status_change' type.
//
// Definition of a Web RCA event.
type StatusChange struct {
	fieldSet_ []bool
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	status    interface{}
	statusId  string
	updatedAt time.Time
}

// Kind returns the name of the type of the object.
func (o *StatusChange) Kind() string {
	if o == nil {
		return StatusChangeNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return StatusChangeLinkKind
	}
	return StatusChangeKind
}

// Link returns true if this is a link.
func (o *StatusChange) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *StatusChange) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *StatusChange) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *StatusChange) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *StatusChange) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *StatusChange) Empty() bool {
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
//
// Object creation timestamp.
func (o *StatusChange) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *StatusChange) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.createdAt
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *StatusChange) DeletedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *StatusChange) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.deletedAt
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *StatusChange) Status() interface{} {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *StatusChange) GetStatus() (value interface{}, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.status
	}
	return
}

// StatusId returns the value of the 'status_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *StatusChange) StatusId() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.statusId
	}
	return ""
}

// GetStatusId returns the value of the 'status_id' attribute and
// a flag indicating if the attribute has a value.
func (o *StatusChange) GetStatusId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.statusId
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *StatusChange) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *StatusChange) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.updatedAt
	}
	return
}

// StatusChangeListKind is the name of the type used to represent list of objects of
// type 'status_change'.
const StatusChangeListKind = "StatusChangeList"

// StatusChangeListLinkKind is the name of the type used to represent links to list
// of objects of type 'status_change'.
const StatusChangeListLinkKind = "StatusChangeListLink"

// StatusChangeNilKind is the name of the type used to nil lists of objects of
// type 'status_change'.
const StatusChangeListNilKind = "StatusChangeListNil"

// StatusChangeList is a list of values of the 'status_change' type.
type StatusChangeList struct {
	href  string
	link  bool
	items []*StatusChange
}

// Kind returns the name of the type of the object.
func (l *StatusChangeList) Kind() string {
	if l == nil {
		return StatusChangeListNilKind
	}
	if l.link {
		return StatusChangeListLinkKind
	}
	return StatusChangeListKind
}

// Link returns true iif this is a link.
func (l *StatusChangeList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *StatusChangeList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *StatusChangeList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *StatusChangeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *StatusChangeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *StatusChangeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *StatusChangeList) SetItems(items []*StatusChange) {
	l.items = items
}

// Items returns the items of the list.
func (l *StatusChangeList) Items() []*StatusChange {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *StatusChangeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *StatusChangeList) Get(i int) *StatusChange {
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
func (l *StatusChangeList) Slice() []*StatusChange {
	var slice []*StatusChange
	if l == nil {
		slice = make([]*StatusChange, 0)
	} else {
		slice = make([]*StatusChange, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *StatusChangeList) Each(f func(item *StatusChange) bool) {
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
func (l *StatusChangeList) Range(f func(index int, item *StatusChange) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
