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

// FollowUpChangeKind is the name of the type used to represent objects
// of type 'follow_up_change'.
const FollowUpChangeKind = "FollowUpChange"

// FollowUpChangeLinkKind is the name of the type used to represent links
// to objects of type 'follow_up_change'.
const FollowUpChangeLinkKind = "FollowUpChangeLink"

// FollowUpChangeNilKind is the name of the type used to nil references
// to objects of type 'follow_up_change'.
const FollowUpChangeNilKind = "FollowUpChangeNil"

// FollowUpChange represents the values of the 'follow_up_change' type.
//
// Definition of a Web RCA event.
type FollowUpChange struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	followUp  *FollowUp
	status    interface{}
	updatedAt time.Time
}

// Kind returns the name of the type of the object.
func (o *FollowUpChange) Kind() string {
	if o == nil {
		return FollowUpChangeNilKind
	}
	if o.bitmap_&1 != 0 {
		return FollowUpChangeLinkKind
	}
	return FollowUpChangeKind
}

// Link returns true iif this is a link.
func (o *FollowUpChange) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *FollowUpChange) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *FollowUpChange) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *FollowUpChange) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *FollowUpChange) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *FollowUpChange) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *FollowUpChange) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *FollowUpChange) GetCreatedAt() (value time.Time, ok bool) {
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
func (o *FollowUpChange) DeletedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *FollowUpChange) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// FollowUp returns the value of the 'follow_up' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUpChange) FollowUp() *FollowUp {
	if o != nil && o.bitmap_&32 != 0 {
		return o.followUp
	}
	return nil
}

// GetFollowUp returns the value of the 'follow_up' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUpChange) GetFollowUp() (value *FollowUp, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.followUp
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUpChange) Status() interface{} {
	if o != nil && o.bitmap_&64 != 0 {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUpChange) GetStatus() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.status
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *FollowUpChange) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&128 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *FollowUpChange) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// FollowUpChangeListKind is the name of the type used to represent list of objects of
// type 'follow_up_change'.
const FollowUpChangeListKind = "FollowUpChangeList"

// FollowUpChangeListLinkKind is the name of the type used to represent links to list
// of objects of type 'follow_up_change'.
const FollowUpChangeListLinkKind = "FollowUpChangeListLink"

// FollowUpChangeNilKind is the name of the type used to nil lists of objects of
// type 'follow_up_change'.
const FollowUpChangeListNilKind = "FollowUpChangeListNil"

// FollowUpChangeList is a list of values of the 'follow_up_change' type.
type FollowUpChangeList struct {
	href  string
	link  bool
	items []*FollowUpChange
}

// Kind returns the name of the type of the object.
func (l *FollowUpChangeList) Kind() string {
	if l == nil {
		return FollowUpChangeListNilKind
	}
	if l.link {
		return FollowUpChangeListLinkKind
	}
	return FollowUpChangeListKind
}

// Link returns true iif this is a link.
func (l *FollowUpChangeList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *FollowUpChangeList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *FollowUpChangeList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *FollowUpChangeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *FollowUpChangeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FollowUpChangeList) Get(i int) *FollowUpChange {
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
func (l *FollowUpChangeList) Slice() []*FollowUpChange {
	var slice []*FollowUpChange
	if l == nil {
		slice = make([]*FollowUpChange, 0)
	} else {
		slice = make([]*FollowUpChange, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FollowUpChangeList) Each(f func(item *FollowUpChange) bool) {
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
func (l *FollowUpChangeList) Range(f func(index int, item *FollowUpChange) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
