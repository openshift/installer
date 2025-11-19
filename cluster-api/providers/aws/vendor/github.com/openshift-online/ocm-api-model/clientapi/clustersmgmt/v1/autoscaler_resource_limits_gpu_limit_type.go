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

// AutoscalerResourceLimitsGPULimit represents the values of the 'autoscaler_resource_limits_GPU_limit' type.
type AutoscalerResourceLimitsGPULimit struct {
	fieldSet_ []bool
	range_    *ResourceRange
	type_     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AutoscalerResourceLimitsGPULimit) Empty() bool {
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

// Range returns the value of the 'range' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AutoscalerResourceLimitsGPULimit) Range() *ResourceRange {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.range_
	}
	return nil
}

// GetRange returns the value of the 'range' attribute and
// a flag indicating if the attribute has a value.
func (o *AutoscalerResourceLimitsGPULimit) GetRange() (value *ResourceRange, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.range_
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The type of GPU to associate with the minimum and maximum limits.
// This value is used by the Cluster Autoscaler to identify Nodes that will have GPU capacity by searching
// for it as a label value on the Node objects. For example, Nodes that carry the label key
// `cluster-api/accelerator` with the label value being the same as the Type field will be counted towards
// the resource limits by the Cluster Autoscaler.
func (o *AutoscalerResourceLimitsGPULimit) Type() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// The type of GPU to associate with the minimum and maximum limits.
// This value is used by the Cluster Autoscaler to identify Nodes that will have GPU capacity by searching
// for it as a label value on the Node objects. For example, Nodes that carry the label key
// `cluster-api/accelerator` with the label value being the same as the Type field will be counted towards
// the resource limits by the Cluster Autoscaler.
func (o *AutoscalerResourceLimitsGPULimit) GetType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.type_
	}
	return
}

// AutoscalerResourceLimitsGPULimitListKind is the name of the type used to represent list of objects of
// type 'autoscaler_resource_limits_GPU_limit'.
const AutoscalerResourceLimitsGPULimitListKind = "AutoscalerResourceLimitsGPULimitList"

// AutoscalerResourceLimitsGPULimitListLinkKind is the name of the type used to represent links to list
// of objects of type 'autoscaler_resource_limits_GPU_limit'.
const AutoscalerResourceLimitsGPULimitListLinkKind = "AutoscalerResourceLimitsGPULimitListLink"

// AutoscalerResourceLimitsGPULimitNilKind is the name of the type used to nil lists of objects of
// type 'autoscaler_resource_limits_GPU_limit'.
const AutoscalerResourceLimitsGPULimitListNilKind = "AutoscalerResourceLimitsGPULimitListNil"

// AutoscalerResourceLimitsGPULimitList is a list of values of the 'autoscaler_resource_limits_GPU_limit' type.
type AutoscalerResourceLimitsGPULimitList struct {
	href  string
	link  bool
	items []*AutoscalerResourceLimitsGPULimit
}

// Len returns the length of the list.
func (l *AutoscalerResourceLimitsGPULimitList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsGPULimitList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsGPULimitList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsGPULimitList) SetItems(items []*AutoscalerResourceLimitsGPULimit) {
	l.items = items
}

// Items returns the items of the list.
func (l *AutoscalerResourceLimitsGPULimitList) Items() []*AutoscalerResourceLimitsGPULimit {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AutoscalerResourceLimitsGPULimitList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AutoscalerResourceLimitsGPULimitList) Get(i int) *AutoscalerResourceLimitsGPULimit {
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
func (l *AutoscalerResourceLimitsGPULimitList) Slice() []*AutoscalerResourceLimitsGPULimit {
	var slice []*AutoscalerResourceLimitsGPULimit
	if l == nil {
		slice = make([]*AutoscalerResourceLimitsGPULimit, 0)
	} else {
		slice = make([]*AutoscalerResourceLimitsGPULimit, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AutoscalerResourceLimitsGPULimitList) Each(f func(item *AutoscalerResourceLimitsGPULimit) bool) {
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
func (l *AutoscalerResourceLimitsGPULimitList) Range(f func(index int, item *AutoscalerResourceLimitsGPULimit) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
