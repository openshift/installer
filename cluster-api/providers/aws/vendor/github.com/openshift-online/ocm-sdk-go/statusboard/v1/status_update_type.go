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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	time "time"
)

// StatusUpdateKind is the name of the type used to represent objects
// of type 'status_update'.
const StatusUpdateKind = "StatusUpdate"

// StatusUpdateLinkKind is the name of the type used to represent links
// to objects of type 'status_update'.
const StatusUpdateLinkKind = "StatusUpdateLink"

// StatusUpdateNilKind is the name of the type used to nil references
// to objects of type 'status_update'.
const StatusUpdateNilKind = "StatusUpdateNil"

// StatusUpdate represents the values of the 'status_update' type.
type StatusUpdate struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	metadata    interface{}
	service     *Service
	serviceInfo *ServiceInfo
	status      string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *StatusUpdate) Kind() string {
	if o == nil {
		return StatusUpdateNilKind
	}
	if o.bitmap_&1 != 0 {
		return StatusUpdateLinkKind
	}
	return StatusUpdateKind
}

// Link returns true if this is a link.
func (o *StatusUpdate) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *StatusUpdate) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *StatusUpdate) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *StatusUpdate) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *StatusUpdate) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *StatusUpdate) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *StatusUpdate) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *StatusUpdate) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Metadata returns the value of the 'metadata' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Miscellaneous metadata about the application.
func (o *StatusUpdate) Metadata() interface{} {
	if o != nil && o.bitmap_&16 != 0 {
		return o.metadata
	}
	return nil
}

// GetMetadata returns the value of the 'metadata' attribute and
// a flag indicating if the attribute has a value.
//
// Miscellaneous metadata about the application.
func (o *StatusUpdate) GetMetadata() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.metadata
	}
	return
}

// Service returns the value of the 'service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The associated service ID.
func (o *StatusUpdate) Service() *Service {
	if o != nil && o.bitmap_&32 != 0 {
		return o.service
	}
	return nil
}

// GetService returns the value of the 'service' attribute and
// a flag indicating if the attribute has a value.
//
// The associated service ID.
func (o *StatusUpdate) GetService() (value *Service, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.service
	}
	return
}

// ServiceInfo returns the value of the 'service_info' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional Service related data.
func (o *StatusUpdate) ServiceInfo() *ServiceInfo {
	if o != nil && o.bitmap_&64 != 0 {
		return o.serviceInfo
	}
	return nil
}

// GetServiceInfo returns the value of the 'service_info' attribute and
// a flag indicating if the attribute has a value.
//
// Additional Service related data.
func (o *StatusUpdate) GetServiceInfo() (value *ServiceInfo, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.serviceInfo
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A status message for the given service.
func (o *StatusUpdate) Status() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// A status message for the given service.
func (o *StatusUpdate) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.status
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *StatusUpdate) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&256 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *StatusUpdate) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// StatusUpdateListKind is the name of the type used to represent list of objects of
// type 'status_update'.
const StatusUpdateListKind = "StatusUpdateList"

// StatusUpdateListLinkKind is the name of the type used to represent links to list
// of objects of type 'status_update'.
const StatusUpdateListLinkKind = "StatusUpdateListLink"

// StatusUpdateNilKind is the name of the type used to nil lists of objects of
// type 'status_update'.
const StatusUpdateListNilKind = "StatusUpdateListNil"

// StatusUpdateList is a list of values of the 'status_update' type.
type StatusUpdateList struct {
	href  string
	link  bool
	items []*StatusUpdate
}

// Kind returns the name of the type of the object.
func (l *StatusUpdateList) Kind() string {
	if l == nil {
		return StatusUpdateListNilKind
	}
	if l.link {
		return StatusUpdateListLinkKind
	}
	return StatusUpdateListKind
}

// Link returns true iif this is a link.
func (l *StatusUpdateList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *StatusUpdateList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *StatusUpdateList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *StatusUpdateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *StatusUpdateList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *StatusUpdateList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *StatusUpdateList) SetItems(items []*StatusUpdate) {
	l.items = items
}

// Items returns the items of the list.
func (l *StatusUpdateList) Items() []*StatusUpdate {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *StatusUpdateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *StatusUpdateList) Get(i int) *StatusUpdate {
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
func (l *StatusUpdateList) Slice() []*StatusUpdate {
	var slice []*StatusUpdate
	if l == nil {
		slice = make([]*StatusUpdate, 0)
	} else {
		slice = make([]*StatusUpdate, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *StatusUpdateList) Each(f func(item *StatusUpdate) bool) {
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
func (l *StatusUpdateList) Range(f func(index int, item *StatusUpdate) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
