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

// KeyRingKind is the name of the type used to represent objects
// of type 'key_ring'.
const KeyRingKind = "KeyRing"

// KeyRingLinkKind is the name of the type used to represent links
// to objects of type 'key_ring'.
const KeyRingLinkKind = "KeyRingLink"

// KeyRingNilKind is the name of the type used to nil references
// to objects of type 'key_ring'.
const KeyRingNilKind = "KeyRingNil"

// KeyRing represents the values of the 'key_ring' type.
//
// Description of a cloud provider key ring.
type KeyRing struct {
	bitmap_ uint32
	id      string
	href    string
	name    string
}

// Kind returns the name of the type of the object.
func (o *KeyRing) Kind() string {
	if o == nil {
		return KeyRingNilKind
	}
	if o.bitmap_&1 != 0 {
		return KeyRingLinkKind
	}
	return KeyRingKind
}

// Link returns true if this is a link.
func (o *KeyRing) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *KeyRing) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *KeyRing) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *KeyRing) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *KeyRing) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *KeyRing) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the key ring.
func (o *KeyRing) Name() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the key ring.
func (o *KeyRing) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.name
	}
	return
}

// KeyRingListKind is the name of the type used to represent list of objects of
// type 'key_ring'.
const KeyRingListKind = "KeyRingList"

// KeyRingListLinkKind is the name of the type used to represent links to list
// of objects of type 'key_ring'.
const KeyRingListLinkKind = "KeyRingListLink"

// KeyRingNilKind is the name of the type used to nil lists of objects of
// type 'key_ring'.
const KeyRingListNilKind = "KeyRingListNil"

// KeyRingList is a list of values of the 'key_ring' type.
type KeyRingList struct {
	href  string
	link  bool
	items []*KeyRing
}

// Kind returns the name of the type of the object.
func (l *KeyRingList) Kind() string {
	if l == nil {
		return KeyRingListNilKind
	}
	if l.link {
		return KeyRingListLinkKind
	}
	return KeyRingListKind
}

// Link returns true iif this is a link.
func (l *KeyRingList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *KeyRingList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *KeyRingList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *KeyRingList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *KeyRingList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *KeyRingList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *KeyRingList) SetItems(items []*KeyRing) {
	l.items = items
}

// Items returns the items of the list.
func (l *KeyRingList) Items() []*KeyRing {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *KeyRingList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *KeyRingList) Get(i int) *KeyRing {
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
func (l *KeyRingList) Slice() []*KeyRing {
	var slice []*KeyRing
	if l == nil {
		slice = make([]*KeyRing, 0)
	} else {
		slice = make([]*KeyRing, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *KeyRingList) Each(f func(item *KeyRing) bool) {
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
func (l *KeyRingList) Range(f func(index int, item *KeyRing) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
