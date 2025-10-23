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

// MonitoringStackResources represents the values of the 'monitoring_stack_resources' type.
//
// Representation of Monitoring Stack Resources
type MonitoringStackResources struct {
	bitmap_  uint32
	limits   *MonitoringStackResource
	requests *MonitoringStackResource
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MonitoringStackResources) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Limits returns the value of the 'limits' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the limit of resource for monitoring stack.
func (o *MonitoringStackResources) Limits() *MonitoringStackResource {
	if o != nil && o.bitmap_&1 != 0 {
		return o.limits
	}
	return nil
}

// GetLimits returns the value of the 'limits' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the limit of resource for monitoring stack.
func (o *MonitoringStackResources) GetLimits() (value *MonitoringStackResource, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.limits
	}
	return
}

// Requests returns the value of the 'requests' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the requested amount of resource for monitoring stack.
func (o *MonitoringStackResources) Requests() *MonitoringStackResource {
	if o != nil && o.bitmap_&2 != 0 {
		return o.requests
	}
	return nil
}

// GetRequests returns the value of the 'requests' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the requested amount of resource for monitoring stack.
func (o *MonitoringStackResources) GetRequests() (value *MonitoringStackResource, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.requests
	}
	return
}

// MonitoringStackResourcesListKind is the name of the type used to represent list of objects of
// type 'monitoring_stack_resources'.
const MonitoringStackResourcesListKind = "MonitoringStackResourcesList"

// MonitoringStackResourcesListLinkKind is the name of the type used to represent links to list
// of objects of type 'monitoring_stack_resources'.
const MonitoringStackResourcesListLinkKind = "MonitoringStackResourcesListLink"

// MonitoringStackResourcesNilKind is the name of the type used to nil lists of objects of
// type 'monitoring_stack_resources'.
const MonitoringStackResourcesListNilKind = "MonitoringStackResourcesListNil"

// MonitoringStackResourcesList is a list of values of the 'monitoring_stack_resources' type.
type MonitoringStackResourcesList struct {
	href  string
	link  bool
	items []*MonitoringStackResources
}

// Len returns the length of the list.
func (l *MonitoringStackResourcesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *MonitoringStackResourcesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *MonitoringStackResourcesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *MonitoringStackResourcesList) SetItems(items []*MonitoringStackResources) {
	l.items = items
}

// Items returns the items of the list.
func (l *MonitoringStackResourcesList) Items() []*MonitoringStackResources {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *MonitoringStackResourcesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MonitoringStackResourcesList) Get(i int) *MonitoringStackResources {
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
func (l *MonitoringStackResourcesList) Slice() []*MonitoringStackResources {
	var slice []*MonitoringStackResources
	if l == nil {
		slice = make([]*MonitoringStackResources, 0)
	} else {
		slice = make([]*MonitoringStackResources, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MonitoringStackResourcesList) Each(f func(item *MonitoringStackResources) bool) {
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
func (l *MonitoringStackResourcesList) Range(f func(index int, item *MonitoringStackResources) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
