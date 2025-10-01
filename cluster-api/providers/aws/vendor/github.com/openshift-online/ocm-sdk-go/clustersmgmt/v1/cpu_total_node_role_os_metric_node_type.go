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

import (
	time "time"
)

// CPUTotalNodeRoleOSMetricNode represents the values of the 'CPU_total_node_role_OS_metric_node' type.
//
// Representation of information from telemetry about a the CPU capacity by node role and OS.
type CPUTotalNodeRoleOSMetricNode struct {
	bitmap_         uint32
	cpuTotal        float64
	nodeRoles       []string
	operatingSystem string
	time            time.Time
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CPUTotalNodeRoleOSMetricNode) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CPUTotal returns the value of the 'CPU_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The total CPU capacity of nodes with this set of roles and operating system.
func (o *CPUTotalNodeRoleOSMetricNode) CPUTotal() float64 {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cpuTotal
	}
	return 0.0
}

// GetCPUTotal returns the value of the 'CPU_total' attribute and
// a flag indicating if the attribute has a value.
//
// The total CPU capacity of nodes with this set of roles and operating system.
func (o *CPUTotalNodeRoleOSMetricNode) GetCPUTotal() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cpuTotal
	}
	return
}

// NodeRoles returns the value of the 'node_roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Representation of the node role for a cluster.
func (o *CPUTotalNodeRoleOSMetricNode) NodeRoles() []string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.nodeRoles
	}
	return nil
}

// GetNodeRoles returns the value of the 'node_roles' attribute and
// a flag indicating if the attribute has a value.
//
// Representation of the node role for a cluster.
func (o *CPUTotalNodeRoleOSMetricNode) GetNodeRoles() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.nodeRoles
	}
	return
}

// OperatingSystem returns the value of the 'operating_system' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The operating system.
func (o *CPUTotalNodeRoleOSMetricNode) OperatingSystem() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.operatingSystem
	}
	return ""
}

// GetOperatingSystem returns the value of the 'operating_system' attribute and
// a flag indicating if the attribute has a value.
//
// The operating system.
func (o *CPUTotalNodeRoleOSMetricNode) GetOperatingSystem() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.operatingSystem
	}
	return
}

// Time returns the value of the 'time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CPUTotalNodeRoleOSMetricNode) Time() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.time
	}
	return time.Time{}
}

// GetTime returns the value of the 'time' attribute and
// a flag indicating if the attribute has a value.
func (o *CPUTotalNodeRoleOSMetricNode) GetTime() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.time
	}
	return
}

// CPUTotalNodeRoleOSMetricNodeListKind is the name of the type used to represent list of objects of
// type 'CPU_total_node_role_OS_metric_node'.
const CPUTotalNodeRoleOSMetricNodeListKind = "CPUTotalNodeRoleOSMetricNodeList"

// CPUTotalNodeRoleOSMetricNodeListLinkKind is the name of the type used to represent links to list
// of objects of type 'CPU_total_node_role_OS_metric_node'.
const CPUTotalNodeRoleOSMetricNodeListLinkKind = "CPUTotalNodeRoleOSMetricNodeListLink"

// CPUTotalNodeRoleOSMetricNodeNilKind is the name of the type used to nil lists of objects of
// type 'CPU_total_node_role_OS_metric_node'.
const CPUTotalNodeRoleOSMetricNodeListNilKind = "CPUTotalNodeRoleOSMetricNodeListNil"

// CPUTotalNodeRoleOSMetricNodeList is a list of values of the 'CPU_total_node_role_OS_metric_node' type.
type CPUTotalNodeRoleOSMetricNodeList struct {
	href  string
	link  bool
	items []*CPUTotalNodeRoleOSMetricNode
}

// Len returns the length of the list.
func (l *CPUTotalNodeRoleOSMetricNodeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CPUTotalNodeRoleOSMetricNodeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CPUTotalNodeRoleOSMetricNodeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CPUTotalNodeRoleOSMetricNodeList) SetItems(items []*CPUTotalNodeRoleOSMetricNode) {
	l.items = items
}

// Items returns the items of the list.
func (l *CPUTotalNodeRoleOSMetricNodeList) Items() []*CPUTotalNodeRoleOSMetricNode {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CPUTotalNodeRoleOSMetricNodeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CPUTotalNodeRoleOSMetricNodeList) Get(i int) *CPUTotalNodeRoleOSMetricNode {
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
func (l *CPUTotalNodeRoleOSMetricNodeList) Slice() []*CPUTotalNodeRoleOSMetricNode {
	var slice []*CPUTotalNodeRoleOSMetricNode
	if l == nil {
		slice = make([]*CPUTotalNodeRoleOSMetricNode, 0)
	} else {
		slice = make([]*CPUTotalNodeRoleOSMetricNode, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CPUTotalNodeRoleOSMetricNodeList) Each(f func(item *CPUTotalNodeRoleOSMetricNode) bool) {
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
func (l *CPUTotalNodeRoleOSMetricNodeList) Range(f func(index int, item *CPUTotalNodeRoleOSMetricNode) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
