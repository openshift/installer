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

// NodeInfo represents the values of the 'node_info' type.
//
// Provides information about a node from specific type in the cluster.
type NodeInfo struct {
	bitmap_ uint32
	amount  int
	type_   NodeType
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *NodeInfo) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Amount returns the value of the 'amount' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The amount of the nodes from this type.
func (o *NodeInfo) Amount() int {
	if o != nil && o.bitmap_&1 != 0 {
		return o.amount
	}
	return 0
}

// GetAmount returns the value of the 'amount' attribute and
// a flag indicating if the attribute has a value.
//
// The amount of the nodes from this type.
func (o *NodeInfo) GetAmount() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.amount
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Node type.
func (o *NodeInfo) Type() NodeType {
	if o != nil && o.bitmap_&2 != 0 {
		return o.type_
	}
	return NodeType("")
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// The Node type.
func (o *NodeInfo) GetType() (value NodeType, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.type_
	}
	return
}

// NodeInfoListKind is the name of the type used to represent list of objects of
// type 'node_info'.
const NodeInfoListKind = "NodeInfoList"

// NodeInfoListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_info'.
const NodeInfoListLinkKind = "NodeInfoListLink"

// NodeInfoNilKind is the name of the type used to nil lists of objects of
// type 'node_info'.
const NodeInfoListNilKind = "NodeInfoListNil"

// NodeInfoList is a list of values of the 'node_info' type.
type NodeInfoList struct {
	href  string
	link  bool
	items []*NodeInfo
}

// Len returns the length of the list.
func (l *NodeInfoList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *NodeInfoList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *NodeInfoList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *NodeInfoList) SetItems(items []*NodeInfo) {
	l.items = items
}

// Items returns the items of the list.
func (l *NodeInfoList) Items() []*NodeInfo {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *NodeInfoList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *NodeInfoList) Get(i int) *NodeInfo {
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
func (l *NodeInfoList) Slice() []*NodeInfo {
	var slice []*NodeInfo
	if l == nil {
		slice = make([]*NodeInfo, 0)
	} else {
		slice = make([]*NodeInfo, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *NodeInfoList) Each(f func(item *NodeInfo) bool) {
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
func (l *NodeInfoList) Range(f func(index int, item *NodeInfo) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
