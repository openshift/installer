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

// FollowUpKind is the name of the type used to represent objects
// of type 'follow_up'.
const FollowUpKind = "FollowUp"

// FollowUpLinkKind is the name of the type used to represent links
// to objects of type 'follow_up'.
const FollowUpLinkKind = "FollowUpLink"

// FollowUpNilKind is the name of the type used to nil references
// to objects of type 'follow_up'.
const FollowUpNilKind = "FollowUpNil"

// FollowUp represents the values of the 'follow_up' type.
//
// Definition of a Web RCA event.
type FollowUp struct {
	bitmap_      uint32
	id           string
	href         string
	createdAt    time.Time
	deletedAt    time.Time
	followUpType string
	incident     *Incident
	owner        string
	priority     string
	status       string
	title        string
	updatedAt    time.Time
	url          string
	workedAt     time.Time
	archived     bool
	done         bool
}

// Kind returns the name of the type of the object.
func (o *FollowUp) Kind() string {
	if o == nil {
		return FollowUpNilKind
	}
	if o.bitmap_&1 != 0 {
		return FollowUpLinkKind
	}
	return FollowUpKind
}

// Link returns true iif this is a link.
func (o *FollowUp) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *FollowUp) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *FollowUp) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *FollowUp) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *FollowUp) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *FollowUp) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Archived returns the value of the 'archived' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Archived() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.archived
	}
	return false
}

// GetArchived returns the value of the 'archived' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetArchived() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.archived
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *FollowUp) CreatedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *FollowUp) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *FollowUp) DeletedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *FollowUp) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// Done returns the value of the 'done' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Done() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.done
	}
	return false
}

// GetDone returns the value of the 'done' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetDone() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.done
	}
	return
}

// FollowUpType returns the value of the 'follow_up_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) FollowUpType() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.followUpType
	}
	return ""
}

// GetFollowUpType returns the value of the 'follow_up_type' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetFollowUpType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.followUpType
	}
	return
}

// Incident returns the value of the 'incident' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Incident() *Incident {
	if o != nil && o.bitmap_&256 != 0 {
		return o.incident
	}
	return nil
}

// GetIncident returns the value of the 'incident' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetIncident() (value *Incident, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.incident
	}
	return
}

// Owner returns the value of the 'owner' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Owner() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.owner
	}
	return ""
}

// GetOwner returns the value of the 'owner' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetOwner() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.owner
	}
	return
}

// Priority returns the value of the 'priority' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Priority() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.priority
	}
	return ""
}

// GetPriority returns the value of the 'priority' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetPriority() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.priority
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Status() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.status
	}
	return
}

// Title returns the value of the 'title' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Title() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.title
	}
	return ""
}

// GetTitle returns the value of the 'title' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetTitle() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.title
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *FollowUp) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *FollowUp) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// Url returns the value of the 'url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) Url() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.url
	}
	return ""
}

// GetUrl returns the value of the 'url' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.url
	}
	return
}

// WorkedAt returns the value of the 'worked_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *FollowUp) WorkedAt() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.workedAt
	}
	return time.Time{}
}

// GetWorkedAt returns the value of the 'worked_at' attribute and
// a flag indicating if the attribute has a value.
func (o *FollowUp) GetWorkedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.workedAt
	}
	return
}

// FollowUpListKind is the name of the type used to represent list of objects of
// type 'follow_up'.
const FollowUpListKind = "FollowUpList"

// FollowUpListLinkKind is the name of the type used to represent links to list
// of objects of type 'follow_up'.
const FollowUpListLinkKind = "FollowUpListLink"

// FollowUpNilKind is the name of the type used to nil lists of objects of
// type 'follow_up'.
const FollowUpListNilKind = "FollowUpListNil"

// FollowUpList is a list of values of the 'follow_up' type.
type FollowUpList struct {
	href  string
	link  bool
	items []*FollowUp
}

// Kind returns the name of the type of the object.
func (l *FollowUpList) Kind() string {
	if l == nil {
		return FollowUpListNilKind
	}
	if l.link {
		return FollowUpListLinkKind
	}
	return FollowUpListKind
}

// Link returns true iif this is a link.
func (l *FollowUpList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *FollowUpList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *FollowUpList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *FollowUpList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *FollowUpList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FollowUpList) Get(i int) *FollowUp {
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
func (l *FollowUpList) Slice() []*FollowUp {
	var slice []*FollowUp
	if l == nil {
		slice = make([]*FollowUp, 0)
	} else {
		slice = make([]*FollowUp, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FollowUpList) Each(f func(item *FollowUp) bool) {
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
func (l *FollowUpList) Range(f func(index int, item *FollowUp) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
