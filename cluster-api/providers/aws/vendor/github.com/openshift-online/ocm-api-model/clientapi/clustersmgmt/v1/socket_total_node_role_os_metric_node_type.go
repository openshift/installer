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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	time "time"
)

// SocketTotalNodeRoleOSMetricNode represents the values of the 'socket_total_node_role_OS_metric_node' type.
//
// Representation of information from telemetry about a the socket capacity
// by node role and OS.
type SocketTotalNodeRoleOSMetricNode struct {
	fieldSet_       []bool
	nodeRoles       []string
	operatingSystem string
	socketTotal     float64
	time            time.Time
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SocketTotalNodeRoleOSMetricNode) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// NodeRoles returns the value of the 'node_roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Representation of the node role for a cluster.
func (o *SocketTotalNodeRoleOSMetricNode) NodeRoles() []string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.nodeRoles
	}
	return nil
}

// GetNodeRoles returns the value of the 'node_roles' attribute and
// a flag indicating if the attribute has a value.
//
// Representation of the node role for a cluster.
func (o *SocketTotalNodeRoleOSMetricNode) GetNodeRoles() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.nodeRoles
	}
	return
}

// OperatingSystem returns the value of the 'operating_system' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The operating system.
func (o *SocketTotalNodeRoleOSMetricNode) OperatingSystem() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.operatingSystem
	}
	return ""
}

// GetOperatingSystem returns the value of the 'operating_system' attribute and
// a flag indicating if the attribute has a value.
//
// The operating system.
func (o *SocketTotalNodeRoleOSMetricNode) GetOperatingSystem() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.operatingSystem
	}
	return
}

// SocketTotal returns the value of the 'socket_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The total socket capacity of nodes with this set of roles and operating system.
func (o *SocketTotalNodeRoleOSMetricNode) SocketTotal() float64 {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.socketTotal
	}
	return 0.0
}

// GetSocketTotal returns the value of the 'socket_total' attribute and
// a flag indicating if the attribute has a value.
//
// The total socket capacity of nodes with this set of roles and operating system.
func (o *SocketTotalNodeRoleOSMetricNode) GetSocketTotal() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.socketTotal
	}
	return
}

// Time returns the value of the 'time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *SocketTotalNodeRoleOSMetricNode) Time() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.time
	}
	return time.Time{}
}

// GetTime returns the value of the 'time' attribute and
// a flag indicating if the attribute has a value.
func (o *SocketTotalNodeRoleOSMetricNode) GetTime() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.time
	}
	return
}

// SocketTotalNodeRoleOSMetricNodeListKind is the name of the type used to represent list of objects of
// type 'socket_total_node_role_OS_metric_node'.
const SocketTotalNodeRoleOSMetricNodeListKind = "SocketTotalNodeRoleOSMetricNodeList"

// SocketTotalNodeRoleOSMetricNodeListLinkKind is the name of the type used to represent links to list
// of objects of type 'socket_total_node_role_OS_metric_node'.
const SocketTotalNodeRoleOSMetricNodeListLinkKind = "SocketTotalNodeRoleOSMetricNodeListLink"

// SocketTotalNodeRoleOSMetricNodeNilKind is the name of the type used to nil lists of objects of
// type 'socket_total_node_role_OS_metric_node'.
const SocketTotalNodeRoleOSMetricNodeListNilKind = "SocketTotalNodeRoleOSMetricNodeListNil"

// SocketTotalNodeRoleOSMetricNodeList is a list of values of the 'socket_total_node_role_OS_metric_node' type.
type SocketTotalNodeRoleOSMetricNodeList struct {
	href  string
	link  bool
	items []*SocketTotalNodeRoleOSMetricNode
}

// Len returns the length of the list.
func (l *SocketTotalNodeRoleOSMetricNodeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SocketTotalNodeRoleOSMetricNodeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SocketTotalNodeRoleOSMetricNodeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SocketTotalNodeRoleOSMetricNodeList) SetItems(items []*SocketTotalNodeRoleOSMetricNode) {
	l.items = items
}

// Items returns the items of the list.
func (l *SocketTotalNodeRoleOSMetricNodeList) Items() []*SocketTotalNodeRoleOSMetricNode {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SocketTotalNodeRoleOSMetricNodeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SocketTotalNodeRoleOSMetricNodeList) Get(i int) *SocketTotalNodeRoleOSMetricNode {
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
func (l *SocketTotalNodeRoleOSMetricNodeList) Slice() []*SocketTotalNodeRoleOSMetricNode {
	var slice []*SocketTotalNodeRoleOSMetricNode
	if l == nil {
		slice = make([]*SocketTotalNodeRoleOSMetricNode, 0)
	} else {
		slice = make([]*SocketTotalNodeRoleOSMetricNode, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SocketTotalNodeRoleOSMetricNodeList) Each(f func(item *SocketTotalNodeRoleOSMetricNode) bool) {
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
func (l *SocketTotalNodeRoleOSMetricNodeList) Range(f func(index int, item *SocketTotalNodeRoleOSMetricNode) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
