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

// HandoffKind is the name of the type used to represent objects
// of type 'handoff'.
const HandoffKind = "Handoff"

// HandoffLinkKind is the name of the type used to represent links
// to objects of type 'handoff'.
const HandoffLinkKind = "HandoffLink"

// HandoffNilKind is the name of the type used to nil references
// to objects of type 'handoff'.
const HandoffNilKind = "HandoffNil"

// Handoff represents the values of the 'handoff' type.
//
// Definition of a Web RCA handoff.
type Handoff struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	deletedAt   time.Time
	handoffFrom *User
	handoffTo   *User
	handoffType string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *Handoff) Kind() string {
	if o == nil {
		return HandoffNilKind
	}
	if o.bitmap_&1 != 0 {
		return HandoffLinkKind
	}
	return HandoffKind
}

// Link returns true if this is a link.
func (o *Handoff) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Handoff) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Handoff) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Handoff) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Handoff) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Handoff) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Handoff) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Handoff) GetCreatedAt() (value time.Time, ok bool) {
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
func (o *Handoff) DeletedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *Handoff) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// HandoffFrom returns the value of the 'handoff_from' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Handoff) HandoffFrom() *User {
	if o != nil && o.bitmap_&32 != 0 {
		return o.handoffFrom
	}
	return nil
}

// GetHandoffFrom returns the value of the 'handoff_from' attribute and
// a flag indicating if the attribute has a value.
func (o *Handoff) GetHandoffFrom() (value *User, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.handoffFrom
	}
	return
}

// HandoffTo returns the value of the 'handoff_to' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Handoff) HandoffTo() *User {
	if o != nil && o.bitmap_&64 != 0 {
		return o.handoffTo
	}
	return nil
}

// GetHandoffTo returns the value of the 'handoff_to' attribute and
// a flag indicating if the attribute has a value.
func (o *Handoff) GetHandoffTo() (value *User, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.handoffTo
	}
	return
}

// HandoffType returns the value of the 'handoff_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Handoff) HandoffType() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.handoffType
	}
	return ""
}

// GetHandoffType returns the value of the 'handoff_type' attribute and
// a flag indicating if the attribute has a value.
func (o *Handoff) GetHandoffType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.handoffType
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Handoff) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Handoff) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// HandoffListKind is the name of the type used to represent list of objects of
// type 'handoff'.
const HandoffListKind = "HandoffList"

// HandoffListLinkKind is the name of the type used to represent links to list
// of objects of type 'handoff'.
const HandoffListLinkKind = "HandoffListLink"

// HandoffNilKind is the name of the type used to nil lists of objects of
// type 'handoff'.
const HandoffListNilKind = "HandoffListNil"

// HandoffList is a list of values of the 'handoff' type.
type HandoffList struct {
	href  string
	link  bool
	items []*Handoff
}

// Kind returns the name of the type of the object.
func (l *HandoffList) Kind() string {
	if l == nil {
		return HandoffListNilKind
	}
	if l.link {
		return HandoffListLinkKind
	}
	return HandoffListKind
}

// Link returns true iif this is a link.
func (l *HandoffList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *HandoffList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *HandoffList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *HandoffList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *HandoffList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *HandoffList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *HandoffList) SetItems(items []*Handoff) {
	l.items = items
}

// Items returns the items of the list.
func (l *HandoffList) Items() []*Handoff {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *HandoffList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *HandoffList) Get(i int) *Handoff {
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
func (l *HandoffList) Slice() []*Handoff {
	var slice []*Handoff
	if l == nil {
		slice = make([]*Handoff, 0)
	} else {
		slice = make([]*Handoff, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *HandoffList) Each(f func(item *Handoff) bool) {
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
func (l *HandoffList) Range(f func(index int, item *Handoff) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
