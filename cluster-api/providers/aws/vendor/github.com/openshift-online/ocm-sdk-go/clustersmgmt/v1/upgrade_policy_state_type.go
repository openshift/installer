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

// UpgradePolicyStateKind is the name of the type used to represent objects
// of type 'upgrade_policy_state'.
const UpgradePolicyStateKind = "UpgradePolicyState"

// UpgradePolicyStateLinkKind is the name of the type used to represent links
// to objects of type 'upgrade_policy_state'.
const UpgradePolicyStateLinkKind = "UpgradePolicyStateLink"

// UpgradePolicyStateNilKind is the name of the type used to nil references
// to objects of type 'upgrade_policy_state'.
const UpgradePolicyStateNilKind = "UpgradePolicyStateNil"

// UpgradePolicyState represents the values of the 'upgrade_policy_state' type.
//
// Representation of an upgrade policy state that that is set for a cluster.
type UpgradePolicyState struct {
	bitmap_     uint32
	id          string
	href        string
	description string
	value       UpgradePolicyStateValue
}

// Kind returns the name of the type of the object.
func (o *UpgradePolicyState) Kind() string {
	if o == nil {
		return UpgradePolicyStateNilKind
	}
	if o.bitmap_&1 != 0 {
		return UpgradePolicyStateLinkKind
	}
	return UpgradePolicyStateKind
}

// Link returns true if this is a link.
func (o *UpgradePolicyState) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *UpgradePolicyState) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *UpgradePolicyState) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *UpgradePolicyState) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *UpgradePolicyState) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *UpgradePolicyState) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the state.
func (o *UpgradePolicyState) Description() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the state.
func (o *UpgradePolicyState) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.description
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State value can be 'pending', 'scheduled', 'cancelled', 'started', 'delayed',
// 'failed' or 'completed'.
func (o *UpgradePolicyState) Value() UpgradePolicyStateValue {
	if o != nil && o.bitmap_&16 != 0 {
		return o.value
	}
	return UpgradePolicyStateValue("")
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// State value can be 'pending', 'scheduled', 'cancelled', 'started', 'delayed',
// 'failed' or 'completed'.
func (o *UpgradePolicyState) GetValue() (value UpgradePolicyStateValue, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.value
	}
	return
}

// UpgradePolicyStateListKind is the name of the type used to represent list of objects of
// type 'upgrade_policy_state'.
const UpgradePolicyStateListKind = "UpgradePolicyStateList"

// UpgradePolicyStateListLinkKind is the name of the type used to represent links to list
// of objects of type 'upgrade_policy_state'.
const UpgradePolicyStateListLinkKind = "UpgradePolicyStateListLink"

// UpgradePolicyStateNilKind is the name of the type used to nil lists of objects of
// type 'upgrade_policy_state'.
const UpgradePolicyStateListNilKind = "UpgradePolicyStateListNil"

// UpgradePolicyStateList is a list of values of the 'upgrade_policy_state' type.
type UpgradePolicyStateList struct {
	href  string
	link  bool
	items []*UpgradePolicyState
}

// Kind returns the name of the type of the object.
func (l *UpgradePolicyStateList) Kind() string {
	if l == nil {
		return UpgradePolicyStateListNilKind
	}
	if l.link {
		return UpgradePolicyStateListLinkKind
	}
	return UpgradePolicyStateListKind
}

// Link returns true iif this is a link.
func (l *UpgradePolicyStateList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *UpgradePolicyStateList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *UpgradePolicyStateList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *UpgradePolicyStateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *UpgradePolicyStateList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *UpgradePolicyStateList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *UpgradePolicyStateList) SetItems(items []*UpgradePolicyState) {
	l.items = items
}

// Items returns the items of the list.
func (l *UpgradePolicyStateList) Items() []*UpgradePolicyState {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *UpgradePolicyStateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *UpgradePolicyStateList) Get(i int) *UpgradePolicyState {
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
func (l *UpgradePolicyStateList) Slice() []*UpgradePolicyState {
	var slice []*UpgradePolicyState
	if l == nil {
		slice = make([]*UpgradePolicyState, 0)
	} else {
		slice = make([]*UpgradePolicyState, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *UpgradePolicyStateList) Each(f func(item *UpgradePolicyState) bool) {
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
func (l *UpgradePolicyStateList) Range(f func(index int, item *UpgradePolicyState) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
