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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// MonitoringStackResource represents the values of the 'monitoring_stack_resource' type.
//
// Representation of Monitoring Stack Resource
type MonitoringStackResource struct {
	bitmap_ uint32
	cpu     string
	memory  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MonitoringStackResource) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Cpu returns the value of the 'cpu' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the CPU resource for monitoring stack.
func (o *MonitoringStackResource) Cpu() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.cpu
	}
	return ""
}

// GetCpu returns the value of the 'cpu' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the CPU resource for monitoring stack.
func (o *MonitoringStackResource) GetCpu() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.cpu
	}
	return
}

// Memory returns the value of the 'memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the memory resource for monitoring stack.
func (o *MonitoringStackResource) Memory() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.memory
	}
	return ""
}

// GetMemory returns the value of the 'memory' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the memory resource for monitoring stack.
func (o *MonitoringStackResource) GetMemory() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.memory
	}
	return
}

// MonitoringStackResourceListKind is the name of the type used to represent list of objects of
// type 'monitoring_stack_resource'.
const MonitoringStackResourceListKind = "MonitoringStackResourceList"

// MonitoringStackResourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'monitoring_stack_resource'.
const MonitoringStackResourceListLinkKind = "MonitoringStackResourceListLink"

// MonitoringStackResourceNilKind is the name of the type used to nil lists of objects of
// type 'monitoring_stack_resource'.
const MonitoringStackResourceListNilKind = "MonitoringStackResourceListNil"

// MonitoringStackResourceList is a list of values of the 'monitoring_stack_resource' type.
type MonitoringStackResourceList struct {
	href  string
	link  bool
	items []*MonitoringStackResource
}

// Len returns the length of the list.
func (l *MonitoringStackResourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *MonitoringStackResourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MonitoringStackResourceList) Get(i int) *MonitoringStackResource {
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
func (l *MonitoringStackResourceList) Slice() []*MonitoringStackResource {
	var slice []*MonitoringStackResource
	if l == nil {
		slice = make([]*MonitoringStackResource, 0)
	} else {
		slice = make([]*MonitoringStackResource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MonitoringStackResourceList) Each(f func(item *MonitoringStackResource) bool) {
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
func (l *MonitoringStackResourceList) Range(f func(index int, item *MonitoringStackResource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
