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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	time "time"
)

// NodePoolStateKind is the name of the type used to represent objects
// of type 'node_pool_state'.
const NodePoolStateKind = "NodePoolState"

// NodePoolStateLinkKind is the name of the type used to represent links
// to objects of type 'node_pool_state'.
const NodePoolStateLinkKind = "NodePoolStateLink"

// NodePoolStateNilKind is the name of the type used to nil references
// to objects of type 'node_pool_state'.
const NodePoolStateNilKind = "NodePoolStateNil"

// NodePoolState represents the values of the 'node_pool_state' type.
//
// Representation of the status of a node pool.
type NodePoolState struct {
	fieldSet_            []bool
	id                   string
	href                 string
	lastUpdatedTimestamp time.Time
	nodePoolStateValue   string
}

// Kind returns the name of the type of the object.
func (o *NodePoolState) Kind() string {
	if o == nil {
		return NodePoolStateNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return NodePoolStateLinkKind
	}
	return NodePoolStateKind
}

// Link returns true if this is a link.
func (o *NodePoolState) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *NodePoolState) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NodePoolState) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NodePoolState) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NodePoolState) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodePoolState) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// LastUpdatedTimestamp returns the value of the 'last_updated_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The current number of replicas for the node pool.
func (o *NodePoolState) LastUpdatedTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.lastUpdatedTimestamp
	}
	return time.Time{}
}

// GetLastUpdatedTimestamp returns the value of the 'last_updated_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// The current number of replicas for the node pool.
func (o *NodePoolState) GetLastUpdatedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.lastUpdatedTimestamp
	}
	return
}

// NodePoolStateValue returns the value of the 'node_pool_state_value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The current state of the node pool
func (o *NodePoolState) NodePoolStateValue() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.nodePoolStateValue
	}
	return ""
}

// GetNodePoolStateValue returns the value of the 'node_pool_state_value' attribute and
// a flag indicating if the attribute has a value.
//
// The current state of the node pool
func (o *NodePoolState) GetNodePoolStateValue() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.nodePoolStateValue
	}
	return
}

// NodePoolStateListKind is the name of the type used to represent list of objects of
// type 'node_pool_state'.
const NodePoolStateListKind = "NodePoolStateList"

// NodePoolStateListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_pool_state'.
const NodePoolStateListLinkKind = "NodePoolStateListLink"

// NodePoolStateNilKind is the name of the type used to nil lists of objects of
// type 'node_pool_state'.
const NodePoolStateListNilKind = "NodePoolStateListNil"

// NodePoolStateList is a list of values of the 'node_pool_state' type.
type NodePoolStateList struct {
	href  string
	link  bool
	items []*NodePoolState
}

// Kind returns the name of the type of the object.
func (l *NodePoolStateList) Kind() string {
	if l == nil {
		return NodePoolStateListNilKind
	}
	if l.link {
		return NodePoolStateListLinkKind
	}
	return NodePoolStateListKind
}

// Link returns true iif this is a link.
func (l *NodePoolStateList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NodePoolStateList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NodePoolStateList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NodePoolStateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NodePoolStateList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NodePoolStateList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NodePoolStateList) SetItems(items []*NodePoolState) {
	l.items = items
}

// Items returns the items of the list.
func (l *NodePoolStateList) Items() []*NodePoolState {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NodePoolStateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodePoolStateList) Get(i int) *NodePoolState {
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
func (l *NodePoolStateList) Slice() []*NodePoolState {
	var slice []*NodePoolState
	if l == nil {
		slice = make([]*NodePoolState, 0)
	} else {
		slice = make([]*NodePoolState, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodePoolStateList) Each(f func(item *NodePoolState) bool) {
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
func (l *NodePoolStateList) Range(f func(index int, item *NodePoolState) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
