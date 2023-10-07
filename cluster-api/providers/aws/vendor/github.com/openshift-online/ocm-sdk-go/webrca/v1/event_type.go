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

// EventKind is the name of the type used to represent objects
// of type 'event'.
const EventKind = "Event"

// EventLinkKind is the name of the type used to represent links
// to objects of type 'event'.
const EventLinkKind = "EventLink"

// EventNilKind is the name of the type used to nil references
// to objects of type 'event'.
const EventNilKind = "EventNil"

// Event represents the values of the 'event' type.
//
// Definition of a Web RCA event.
type Event struct {
	bitmap_              uint32
	id                   string
	href                 string
	createdAt            time.Time
	creator              *User
	deletedAt            time.Time
	escalation           *Escalation
	eventType            string
	externalReferenceUrl string
	followUp             *FollowUp
	followUpChange       *FollowUpChange
	handoff              *Handoff
	incident             *Incident
	note                 string
	statusChange         *StatusChange
	updatedAt            time.Time
}

// Kind returns the name of the type of the object.
func (o *Event) Kind() string {
	if o == nil {
		return EventNilKind
	}
	if o.bitmap_&1 != 0 {
		return EventLinkKind
	}
	return EventKind
}

// Link returns true iif this is a link.
func (o *Event) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Event) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Event) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Event) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Event) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Event) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Event) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Event) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// Creator returns the value of the 'creator' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) Creator() *User {
	if o != nil && o.bitmap_&16 != 0 {
		return o.creator
	}
	return nil
}

// GetCreator returns the value of the 'creator' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetCreator() (value *User, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.creator
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *Event) DeletedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *Event) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// Escalation returns the value of the 'escalation' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) Escalation() *Escalation {
	if o != nil && o.bitmap_&64 != 0 {
		return o.escalation
	}
	return nil
}

// GetEscalation returns the value of the 'escalation' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetEscalation() (value *Escalation, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.escalation
	}
	return
}

// EventType returns the value of the 'event_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) EventType() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.eventType
	}
	return ""
}

// GetEventType returns the value of the 'event_type' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetEventType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.eventType
	}
	return
}

// ExternalReferenceUrl returns the value of the 'external_reference_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) ExternalReferenceUrl() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.externalReferenceUrl
	}
	return ""
}

// GetExternalReferenceUrl returns the value of the 'external_reference_url' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetExternalReferenceUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.externalReferenceUrl
	}
	return
}

// FollowUp returns the value of the 'follow_up' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) FollowUp() *FollowUp {
	if o != nil && o.bitmap_&512 != 0 {
		return o.followUp
	}
	return nil
}

// GetFollowUp returns the value of the 'follow_up' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetFollowUp() (value *FollowUp, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.followUp
	}
	return
}

// FollowUpChange returns the value of the 'follow_up_change' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) FollowUpChange() *FollowUpChange {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.followUpChange
	}
	return nil
}

// GetFollowUpChange returns the value of the 'follow_up_change' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetFollowUpChange() (value *FollowUpChange, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.followUpChange
	}
	return
}

// Handoff returns the value of the 'handoff' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) Handoff() *Handoff {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.handoff
	}
	return nil
}

// GetHandoff returns the value of the 'handoff' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetHandoff() (value *Handoff, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.handoff
	}
	return
}

// Incident returns the value of the 'incident' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) Incident() *Incident {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.incident
	}
	return nil
}

// GetIncident returns the value of the 'incident' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetIncident() (value *Incident, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.incident
	}
	return
}

// Note returns the value of the 'note' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) Note() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.note
	}
	return ""
}

// GetNote returns the value of the 'note' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetNote() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.note
	}
	return
}

// StatusChange returns the value of the 'status_change' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Event) StatusChange() *StatusChange {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.statusChange
	}
	return nil
}

// GetStatusChange returns the value of the 'status_change' attribute and
// a flag indicating if the attribute has a value.
func (o *Event) GetStatusChange() (value *StatusChange, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.statusChange
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Event) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Event) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// EventListKind is the name of the type used to represent list of objects of
// type 'event'.
const EventListKind = "EventList"

// EventListLinkKind is the name of the type used to represent links to list
// of objects of type 'event'.
const EventListLinkKind = "EventListLink"

// EventNilKind is the name of the type used to nil lists of objects of
// type 'event'.
const EventListNilKind = "EventListNil"

// EventList is a list of values of the 'event' type.
type EventList struct {
	href  string
	link  bool
	items []*Event
}

// Kind returns the name of the type of the object.
func (l *EventList) Kind() string {
	if l == nil {
		return EventListNilKind
	}
	if l.link {
		return EventListLinkKind
	}
	return EventListKind
}

// Link returns true iif this is a link.
func (l *EventList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *EventList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *EventList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *EventList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *EventList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *EventList) Get(i int) *Event {
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
func (l *EventList) Slice() []*Event {
	var slice []*Event
	if l == nil {
		slice = make([]*Event, 0)
	} else {
		slice = make([]*Event, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *EventList) Each(f func(item *Event) bool) {
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
func (l *EventList) Range(f func(index int, item *Event) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
