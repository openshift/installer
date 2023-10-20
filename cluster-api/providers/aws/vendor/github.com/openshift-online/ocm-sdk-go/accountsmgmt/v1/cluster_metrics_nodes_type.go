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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// ClusterMetricsNodes represents the values of the 'cluster_metrics_nodes' type.
type ClusterMetricsNodes struct {
	bitmap_ uint32
	compute float64
	infra   float64
	master  float64
	total   float64
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterMetricsNodes) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Compute returns the value of the 'compute' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterMetricsNodes) Compute() float64 {
	if o != nil && o.bitmap_&1 != 0 {
		return o.compute
	}
	return 0.0
}

// GetCompute returns the value of the 'compute' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterMetricsNodes) GetCompute() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.compute
	}
	return
}

// Infra returns the value of the 'infra' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterMetricsNodes) Infra() float64 {
	if o != nil && o.bitmap_&2 != 0 {
		return o.infra
	}
	return 0.0
}

// GetInfra returns the value of the 'infra' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterMetricsNodes) GetInfra() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.infra
	}
	return
}

// Master returns the value of the 'master' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterMetricsNodes) Master() float64 {
	if o != nil && o.bitmap_&4 != 0 {
		return o.master
	}
	return 0.0
}

// GetMaster returns the value of the 'master' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterMetricsNodes) GetMaster() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.master
	}
	return
}

// Total returns the value of the 'total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ClusterMetricsNodes) Total() float64 {
	if o != nil && o.bitmap_&8 != 0 {
		return o.total
	}
	return 0.0
}

// GetTotal returns the value of the 'total' attribute and
// a flag indicating if the attribute has a value.
func (o *ClusterMetricsNodes) GetTotal() (value float64, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.total
	}
	return
}

// ClusterMetricsNodesListKind is the name of the type used to represent list of objects of
// type 'cluster_metrics_nodes'.
const ClusterMetricsNodesListKind = "ClusterMetricsNodesList"

// ClusterMetricsNodesListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_metrics_nodes'.
const ClusterMetricsNodesListLinkKind = "ClusterMetricsNodesListLink"

// ClusterMetricsNodesNilKind is the name of the type used to nil lists of objects of
// type 'cluster_metrics_nodes'.
const ClusterMetricsNodesListNilKind = "ClusterMetricsNodesListNil"

// ClusterMetricsNodesList is a list of values of the 'cluster_metrics_nodes' type.
type ClusterMetricsNodesList struct {
	href  string
	link  bool
	items []*ClusterMetricsNodes
}

// Len returns the length of the list.
func (l *ClusterMetricsNodesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterMetricsNodesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterMetricsNodesList) Get(i int) *ClusterMetricsNodes {
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
func (l *ClusterMetricsNodesList) Slice() []*ClusterMetricsNodes {
	var slice []*ClusterMetricsNodes
	if l == nil {
		slice = make([]*ClusterMetricsNodes, 0)
	} else {
		slice = make([]*ClusterMetricsNodes, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterMetricsNodesList) Each(f func(item *ClusterMetricsNodes) bool) {
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
func (l *ClusterMetricsNodesList) Range(f func(index int, item *ClusterMetricsNodes) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
