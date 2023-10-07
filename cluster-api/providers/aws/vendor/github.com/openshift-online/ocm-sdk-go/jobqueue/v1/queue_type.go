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

package v1 // github.com/openshift-online/ocm-sdk-go/jobqueue/v1

import (
	time "time"
)

// QueueKind is the name of the type used to represent objects
// of type 'queue'.
const QueueKind = "Queue"

// QueueLinkKind is the name of the type used to represent links
// to objects of type 'queue'.
const QueueLinkKind = "QueueLink"

// QueueNilKind is the name of the type used to nil references
// to objects of type 'queue'.
const QueueNilKind = "QueueNil"

// Queue represents the values of the 'queue' type.
type Queue struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	maxAttempts int
	maxRunTime  int
	name        string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *Queue) Kind() string {
	if o == nil {
		return QueueNilKind
	}
	if o.bitmap_&1 != 0 {
		return QueueLinkKind
	}
	return QueueKind
}

// Link returns true iif this is a link.
func (o *Queue) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Queue) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Queue) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Queue) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Queue) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Queue) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Queue) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Queue) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// MaxAttempts returns the value of the 'max_attempts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// SQS Visibility Timeout
func (o *Queue) MaxAttempts() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.maxAttempts
	}
	return 0
}

// GetMaxAttempts returns the value of the 'max_attempts' attribute and
// a flag indicating if the attribute has a value.
//
// SQS Visibility Timeout
func (o *Queue) GetMaxAttempts() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.maxAttempts
	}
	return
}

// MaxRunTime returns the value of the 'max_run_time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Queue) MaxRunTime() int {
	if o != nil && o.bitmap_&32 != 0 {
		return o.maxRunTime
	}
	return 0
}

// GetMaxRunTime returns the value of the 'max_run_time' attribute and
// a flag indicating if the attribute has a value.
func (o *Queue) GetMaxRunTime() (value int, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.maxRunTime
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Queue) Name() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Queue) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.name
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Queue) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&128 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Queue) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// QueueListKind is the name of the type used to represent list of objects of
// type 'queue'.
const QueueListKind = "QueueList"

// QueueListLinkKind is the name of the type used to represent links to list
// of objects of type 'queue'.
const QueueListLinkKind = "QueueListLink"

// QueueNilKind is the name of the type used to nil lists of objects of
// type 'queue'.
const QueueListNilKind = "QueueListNil"

// QueueList is a list of values of the 'queue' type.
type QueueList struct {
	href  string
	link  bool
	items []*Queue
}

// Kind returns the name of the type of the object.
func (l *QueueList) Kind() string {
	if l == nil {
		return QueueListNilKind
	}
	if l.link {
		return QueueListLinkKind
	}
	return QueueListKind
}

// Link returns true iif this is a link.
func (l *QueueList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *QueueList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *QueueList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *QueueList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *QueueList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *QueueList) Get(i int) *Queue {
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
func (l *QueueList) Slice() []*Queue {
	var slice []*Queue
	if l == nil {
		slice = make([]*Queue, 0)
	} else {
		slice = make([]*Queue, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *QueueList) Each(f func(item *Queue) bool) {
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
func (l *QueueList) Range(f func(index int, item *Queue) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
