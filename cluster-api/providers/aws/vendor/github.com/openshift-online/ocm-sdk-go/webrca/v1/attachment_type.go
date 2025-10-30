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

// AttachmentKind is the name of the type used to represent objects
// of type 'attachment'.
const AttachmentKind = "Attachment"

// AttachmentLinkKind is the name of the type used to represent links
// to objects of type 'attachment'.
const AttachmentLinkKind = "AttachmentLink"

// AttachmentNilKind is the name of the type used to nil references
// to objects of type 'attachment'.
const AttachmentNilKind = "AttachmentNil"

// Attachment represents the values of the 'attachment' type.
//
// Definition of a Web RCA attachment.
type Attachment struct {
	bitmap_     uint32
	id          string
	href        string
	contentType string
	createdAt   time.Time
	creator     *User
	deletedAt   time.Time
	event       *Event
	fileSize    int
	name        string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *Attachment) Kind() string {
	if o == nil {
		return AttachmentNilKind
	}
	if o.bitmap_&1 != 0 {
		return AttachmentLinkKind
	}
	return AttachmentKind
}

// Link returns true if this is a link.
func (o *Attachment) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Attachment) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Attachment) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Attachment) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Attachment) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Attachment) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ContentType returns the value of the 'content_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Attachment) ContentType() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.contentType
	}
	return ""
}

// GetContentType returns the value of the 'content_type' attribute and
// a flag indicating if the attribute has a value.
func (o *Attachment) GetContentType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.contentType
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Attachment) CreatedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Attachment) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Creator returns the value of the 'creator' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Attachment) Creator() *User {
	if o != nil && o.bitmap_&32 != 0 {
		return o.creator
	}
	return nil
}

// GetCreator returns the value of the 'creator' attribute and
// a flag indicating if the attribute has a value.
func (o *Attachment) GetCreator() (value *User, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.creator
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *Attachment) DeletedAt() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *Attachment) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// Event returns the value of the 'event' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Attachment) Event() *Event {
	if o != nil && o.bitmap_&128 != 0 {
		return o.event
	}
	return nil
}

// GetEvent returns the value of the 'event' attribute and
// a flag indicating if the attribute has a value.
func (o *Attachment) GetEvent() (value *Event, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.event
	}
	return
}

// FileSize returns the value of the 'file_size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Attachment) FileSize() int {
	if o != nil && o.bitmap_&256 != 0 {
		return o.fileSize
	}
	return 0
}

// GetFileSize returns the value of the 'file_size' attribute and
// a flag indicating if the attribute has a value.
func (o *Attachment) GetFileSize() (value int, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.fileSize
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Attachment) Name() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *Attachment) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.name
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Attachment) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Attachment) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// AttachmentListKind is the name of the type used to represent list of objects of
// type 'attachment'.
const AttachmentListKind = "AttachmentList"

// AttachmentListLinkKind is the name of the type used to represent links to list
// of objects of type 'attachment'.
const AttachmentListLinkKind = "AttachmentListLink"

// AttachmentNilKind is the name of the type used to nil lists of objects of
// type 'attachment'.
const AttachmentListNilKind = "AttachmentListNil"

// AttachmentList is a list of values of the 'attachment' type.
type AttachmentList struct {
	href  string
	link  bool
	items []*Attachment
}

// Kind returns the name of the type of the object.
func (l *AttachmentList) Kind() string {
	if l == nil {
		return AttachmentListNilKind
	}
	if l.link {
		return AttachmentListLinkKind
	}
	return AttachmentListKind
}

// Link returns true iif this is a link.
func (l *AttachmentList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AttachmentList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AttachmentList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AttachmentList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AttachmentList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AttachmentList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AttachmentList) SetItems(items []*Attachment) {
	l.items = items
}

// Items returns the items of the list.
func (l *AttachmentList) Items() []*Attachment {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AttachmentList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AttachmentList) Get(i int) *Attachment {
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
func (l *AttachmentList) Slice() []*Attachment {
	var slice []*Attachment
	if l == nil {
		slice = make([]*Attachment, 0)
	} else {
		slice = make([]*Attachment, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AttachmentList) Each(f func(item *Attachment) bool) {
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
func (l *AttachmentList) Range(f func(index int, item *Attachment) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
