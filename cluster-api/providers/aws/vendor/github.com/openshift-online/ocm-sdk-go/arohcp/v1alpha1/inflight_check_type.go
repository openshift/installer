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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	time "time"
)

// InflightCheckKind is the name of the type used to represent objects
// of type 'inflight_check'.
const InflightCheckKind = "InflightCheck"

// InflightCheckLinkKind is the name of the type used to represent links
// to objects of type 'inflight_check'.
const InflightCheckLinkKind = "InflightCheckLink"

// InflightCheckNilKind is the name of the type used to nil references
// to objects of type 'inflight_check'.
const InflightCheckNilKind = "InflightCheckNil"

// InflightCheck represents the values of the 'inflight_check' type.
//
// Representation of check running before the cluster is provisioned.
type InflightCheck struct {
	bitmap_   uint32
	id        string
	href      string
	details   interface{}
	endedAt   time.Time
	name      string
	restarts  int
	startedAt time.Time
	state     InflightCheckState
}

// Kind returns the name of the type of the object.
func (o *InflightCheck) Kind() string {
	if o == nil {
		return InflightCheckNilKind
	}
	if o.bitmap_&1 != 0 {
		return InflightCheckLinkKind
	}
	return InflightCheckKind
}

// Link returns true if this is a link.
func (o *InflightCheck) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *InflightCheck) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *InflightCheck) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *InflightCheck) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *InflightCheck) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *InflightCheck) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Details returns the value of the 'details' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details regarding the state of the inflight check.
func (o *InflightCheck) Details() interface{} {
	if o != nil && o.bitmap_&8 != 0 {
		return o.details
	}
	return nil
}

// GetDetails returns the value of the 'details' attribute and
// a flag indicating if the attribute has a value.
//
// Details regarding the state of the inflight check.
func (o *InflightCheck) GetDetails() (value interface{}, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.details
	}
	return
}

// EndedAt returns the value of the 'ended_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The time the check finished running.
func (o *InflightCheck) EndedAt() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.endedAt
	}
	return time.Time{}
}

// GetEndedAt returns the value of the 'ended_at' attribute and
// a flag indicating if the attribute has a value.
//
// The time the check finished running.
func (o *InflightCheck) GetEndedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.endedAt
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the inflight check.
func (o *InflightCheck) Name() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the inflight check.
func (o *InflightCheck) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.name
	}
	return
}

// Restarts returns the value of the 'restarts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The number of times the inflight check restarted.
func (o *InflightCheck) Restarts() int {
	if o != nil && o.bitmap_&64 != 0 {
		return o.restarts
	}
	return 0
}

// GetRestarts returns the value of the 'restarts' attribute and
// a flag indicating if the attribute has a value.
//
// The number of times the inflight check restarted.
func (o *InflightCheck) GetRestarts() (value int, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.restarts
	}
	return
}

// StartedAt returns the value of the 'started_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The time the check started running.
func (o *InflightCheck) StartedAt() time.Time {
	if o != nil && o.bitmap_&128 != 0 {
		return o.startedAt
	}
	return time.Time{}
}

// GetStartedAt returns the value of the 'started_at' attribute and
// a flag indicating if the attribute has a value.
//
// The time the check started running.
func (o *InflightCheck) GetStartedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.startedAt
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the inflight check.
func (o *InflightCheck) State() InflightCheckState {
	if o != nil && o.bitmap_&256 != 0 {
		return o.state
	}
	return InflightCheckState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State of the inflight check.
func (o *InflightCheck) GetState() (value InflightCheckState, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.state
	}
	return
}

// InflightCheckListKind is the name of the type used to represent list of objects of
// type 'inflight_check'.
const InflightCheckListKind = "InflightCheckList"

// InflightCheckListLinkKind is the name of the type used to represent links to list
// of objects of type 'inflight_check'.
const InflightCheckListLinkKind = "InflightCheckListLink"

// InflightCheckNilKind is the name of the type used to nil lists of objects of
// type 'inflight_check'.
const InflightCheckListNilKind = "InflightCheckListNil"

// InflightCheckList is a list of values of the 'inflight_check' type.
type InflightCheckList struct {
	href  string
	link  bool
	items []*InflightCheck
}

// Kind returns the name of the type of the object.
func (l *InflightCheckList) Kind() string {
	if l == nil {
		return InflightCheckListNilKind
	}
	if l.link {
		return InflightCheckListLinkKind
	}
	return InflightCheckListKind
}

// Link returns true iif this is a link.
func (l *InflightCheckList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *InflightCheckList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *InflightCheckList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *InflightCheckList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *InflightCheckList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *InflightCheckList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *InflightCheckList) SetItems(items []*InflightCheck) {
	l.items = items
}

// Items returns the items of the list.
func (l *InflightCheckList) Items() []*InflightCheck {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *InflightCheckList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *InflightCheckList) Get(i int) *InflightCheck {
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
func (l *InflightCheckList) Slice() []*InflightCheck {
	var slice []*InflightCheck
	if l == nil {
		slice = make([]*InflightCheck, 0)
	} else {
		slice = make([]*InflightCheck, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *InflightCheckList) Each(f func(item *InflightCheck) bool) {
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
func (l *InflightCheckList) Range(f func(index int, item *InflightCheck) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
