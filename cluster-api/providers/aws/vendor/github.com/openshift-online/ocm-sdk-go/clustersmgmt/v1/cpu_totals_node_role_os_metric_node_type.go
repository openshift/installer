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

// CPUTotalsNodeRoleOSMetricNode represents the values of the 'CPU_totals_node_role_OS_metric_node' type.
//
// Representation of information from telemetry about the CPU capacity by node
// role and OS of a cluster.
type CPUTotalsNodeRoleOSMetricNode struct {
	bitmap_   uint32
	cpuTotals []*CPUTotalNodeRoleOSMetricNode
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CPUTotalsNodeRoleOSMetricNode) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CPUTotals returns the value of the 'CPU_totals' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *CPUTotalsNodeRoleOSMetricNode) CPUTotals() []*CPUTotalNodeRoleOSMetricNode {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cpuTotals
	}
	return nil
}

// GetCPUTotals returns the value of the 'CPU_totals' attribute and
// a flag indicating if the attribute has a value.
func (o *CPUTotalsNodeRoleOSMetricNode) GetCPUTotals() (value []*CPUTotalNodeRoleOSMetricNode, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cpuTotals
	}
	return
}

// CPUTotalsNodeRoleOSMetricNodeListKind is the name of the type used to represent list of objects of
// type 'CPU_totals_node_role_OS_metric_node'.
const CPUTotalsNodeRoleOSMetricNodeListKind = "CPUTotalsNodeRoleOSMetricNodeList"

// CPUTotalsNodeRoleOSMetricNodeListLinkKind is the name of the type used to represent links to list
// of objects of type 'CPU_totals_node_role_OS_metric_node'.
const CPUTotalsNodeRoleOSMetricNodeListLinkKind = "CPUTotalsNodeRoleOSMetricNodeListLink"

// CPUTotalsNodeRoleOSMetricNodeNilKind is the name of the type used to nil lists of objects of
// type 'CPU_totals_node_role_OS_metric_node'.
const CPUTotalsNodeRoleOSMetricNodeListNilKind = "CPUTotalsNodeRoleOSMetricNodeListNil"

// CPUTotalsNodeRoleOSMetricNodeList is a list of values of the 'CPU_totals_node_role_OS_metric_node' type.
type CPUTotalsNodeRoleOSMetricNodeList struct {
	href  string
	link  bool
	items []*CPUTotalsNodeRoleOSMetricNode
}

// Len returns the length of the list.
func (l *CPUTotalsNodeRoleOSMetricNodeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CPUTotalsNodeRoleOSMetricNodeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CPUTotalsNodeRoleOSMetricNodeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CPUTotalsNodeRoleOSMetricNodeList) SetItems(items []*CPUTotalsNodeRoleOSMetricNode) {
	l.items = items
}

// Items returns the items of the list.
func (l *CPUTotalsNodeRoleOSMetricNodeList) Items() []*CPUTotalsNodeRoleOSMetricNode {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CPUTotalsNodeRoleOSMetricNodeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CPUTotalsNodeRoleOSMetricNodeList) Get(i int) *CPUTotalsNodeRoleOSMetricNode {
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
func (l *CPUTotalsNodeRoleOSMetricNodeList) Slice() []*CPUTotalsNodeRoleOSMetricNode {
	var slice []*CPUTotalsNodeRoleOSMetricNode
	if l == nil {
		slice = make([]*CPUTotalsNodeRoleOSMetricNode, 0)
	} else {
		slice = make([]*CPUTotalsNodeRoleOSMetricNode, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CPUTotalsNodeRoleOSMetricNodeList) Each(f func(item *CPUTotalsNodeRoleOSMetricNode) bool) {
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
func (l *CPUTotalsNodeRoleOSMetricNodeList) Range(f func(index int, item *CPUTotalsNodeRoleOSMetricNode) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
