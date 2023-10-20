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

// NodePoolStatusKind is the name of the type used to represent objects
// of type 'node_pool_status'.
const NodePoolStatusKind = "NodePoolStatus"

// NodePoolStatusLinkKind is the name of the type used to represent links
// to objects of type 'node_pool_status'.
const NodePoolStatusLinkKind = "NodePoolStatusLink"

// NodePoolStatusNilKind is the name of the type used to nil references
// to objects of type 'node_pool_status'.
const NodePoolStatusNilKind = "NodePoolStatusNil"

// NodePoolStatus represents the values of the 'node_pool_status' type.
//
// Representation of the status of a node pool.
type NodePoolStatus struct {
	bitmap_         uint32
	id              string
	href            string
	currentReplicas int
	message         string
}

// Kind returns the name of the type of the object.
func (o *NodePoolStatus) Kind() string {
	if o == nil {
		return NodePoolStatusNilKind
	}
	if o.bitmap_&1 != 0 {
		return NodePoolStatusLinkKind
	}
	return NodePoolStatusKind
}

// Link returns true iif this is a link.
func (o *NodePoolStatus) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *NodePoolStatus) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *NodePoolStatus) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *NodePoolStatus) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *NodePoolStatus) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodePoolStatus) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CurrentReplicas returns the value of the 'current_replicas' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The current number of replicas for the node pool.
func (o *NodePoolStatus) CurrentReplicas() int {
	if o != nil && o.bitmap_&8 != 0 {
		return o.currentReplicas
	}
	return 0
}

// GetCurrentReplicas returns the value of the 'current_replicas' attribute and
// a flag indicating if the attribute has a value.
//
// The current number of replicas for the node pool.
func (o *NodePoolStatus) GetCurrentReplicas() (value int, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.currentReplicas
	}
	return
}

// Message returns the value of the 'message' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Adds additional information about the NodePool status when the node pool doesn't reach the desired replicas.
func (o *NodePoolStatus) Message() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.message
	}
	return ""
}

// GetMessage returns the value of the 'message' attribute and
// a flag indicating if the attribute has a value.
//
// Adds additional information about the NodePool status when the node pool doesn't reach the desired replicas.
func (o *NodePoolStatus) GetMessage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.message
	}
	return
}

// NodePoolStatusListKind is the name of the type used to represent list of objects of
// type 'node_pool_status'.
const NodePoolStatusListKind = "NodePoolStatusList"

// NodePoolStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_pool_status'.
const NodePoolStatusListLinkKind = "NodePoolStatusListLink"

// NodePoolStatusNilKind is the name of the type used to nil lists of objects of
// type 'node_pool_status'.
const NodePoolStatusListNilKind = "NodePoolStatusListNil"

// NodePoolStatusList is a list of values of the 'node_pool_status' type.
type NodePoolStatusList struct {
	href  string
	link  bool
	items []*NodePoolStatus
}

// Kind returns the name of the type of the object.
func (l *NodePoolStatusList) Kind() string {
	if l == nil {
		return NodePoolStatusListNilKind
	}
	if l.link {
		return NodePoolStatusListLinkKind
	}
	return NodePoolStatusListKind
}

// Link returns true iif this is a link.
func (l *NodePoolStatusList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *NodePoolStatusList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *NodePoolStatusList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *NodePoolStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *NodePoolStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodePoolStatusList) Get(i int) *NodePoolStatus {
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
func (l *NodePoolStatusList) Slice() []*NodePoolStatus {
	var slice []*NodePoolStatus
	if l == nil {
		slice = make([]*NodePoolStatus, 0)
	} else {
		slice = make([]*NodePoolStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodePoolStatusList) Each(f func(item *NodePoolStatus) bool) {
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
func (l *NodePoolStatusList) Range(f func(index int, item *NodePoolStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
