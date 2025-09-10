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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// Event represents the values of the 'event' type.
//
// Representation of a trackable event.
type Event struct {
	bitmap_ uint32
	body    map[string]string
	key     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Event) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Body returns the value of the 'body' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Body of the event to track the details of the tracking event as Key value pair
func (o *Event) Body() map[string]string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.body
	}
	return nil
}

// GetBody returns the value of the 'body' attribute and
// a flag indicating if the attribute has a value.
//
// Body of the event to track the details of the tracking event as Key value pair
func (o *Event) GetBody() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.body
	}
	return
}

// Key returns the value of the 'key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Key of the event to be tracked. This key should start with an
// uppercase letter followed by alphanumeric characters or
// underscores. The entire key needs to be smaller than 64 characters.
func (o *Event) Key() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.key
	}
	return ""
}

// GetKey returns the value of the 'key' attribute and
// a flag indicating if the attribute has a value.
//
// Key of the event to be tracked. This key should start with an
// uppercase letter followed by alphanumeric characters or
// underscores. The entire key needs to be smaller than 64 characters.
func (o *Event) GetKey() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.key
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

// Len returns the length of the list.
func (l *EventList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *EventList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *EventList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *EventList) SetItems(items []*Event) {
	l.items = items
}

// Items returns the items of the list.
func (l *EventList) Items() []*Event {
	if l == nil {
		return nil
	}
	return l.items
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
