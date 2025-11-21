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

// AutoscalerResourceLimits represents the values of the 'autoscaler_resource_limits' type.
type AutoscalerResourceLimits struct {
	fieldSet_     []bool
	gpus          []*AutoscalerResourceLimitsGPULimit
	cores         *ResourceRange
	maxNodesTotal int
	memory        *ResourceRange
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AutoscalerResourceLimits) Empty() bool {
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

// GPUS returns the value of the 'GPUS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Minimum and maximum number of different GPUs in cluster, in the format <gpu_type>:<min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers. Can be passed multiple times.
func (o *AutoscalerResourceLimits) GPUS() []*AutoscalerResourceLimitsGPULimit {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.gpus
	}
	return nil
}

// GetGPUS returns the value of the 'GPUS' attribute and
// a flag indicating if the attribute has a value.
//
// Minimum and maximum number of different GPUs in cluster, in the format <gpu_type>:<min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers. Can be passed multiple times.
func (o *AutoscalerResourceLimits) GetGPUS() (value []*AutoscalerResourceLimitsGPULimit, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.gpus
	}
	return
}

// Cores returns the value of the 'cores' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Minimum and maximum number of cores in cluster, in the format <min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers.
func (o *AutoscalerResourceLimits) Cores() *ResourceRange {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.cores
	}
	return nil
}

// GetCores returns the value of the 'cores' attribute and
// a flag indicating if the attribute has a value.
//
// Minimum and maximum number of cores in cluster, in the format <min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers.
func (o *AutoscalerResourceLimits) GetCores() (value *ResourceRange, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.cores
	}
	return
}

// MaxNodesTotal returns the value of the 'max_nodes_total' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Maximum number of nodes in all node groups.
// Cluster autoscaler will not grow the cluster beyond this number.
func (o *AutoscalerResourceLimits) MaxNodesTotal() int {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.maxNodesTotal
	}
	return 0
}

// GetMaxNodesTotal returns the value of the 'max_nodes_total' attribute and
// a flag indicating if the attribute has a value.
//
// Maximum number of nodes in all node groups.
// Cluster autoscaler will not grow the cluster beyond this number.
func (o *AutoscalerResourceLimits) GetMaxNodesTotal() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.maxNodesTotal
	}
	return
}

// Memory returns the value of the 'memory' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Minimum and maximum number of gigabytes of memory in cluster, in the format <min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers.
func (o *AutoscalerResourceLimits) Memory() *ResourceRange {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.memory
	}
	return nil
}

// GetMemory returns the value of the 'memory' attribute and
// a flag indicating if the attribute has a value.
//
// Minimum and maximum number of gigabytes of memory in cluster, in the format <min>:<max>.
// Cluster autoscaler will not scale the cluster beyond these numbers.
func (o *AutoscalerResourceLimits) GetMemory() (value *ResourceRange, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.memory
	}
	return
}

// AutoscalerResourceLimitsListKind is the name of the type used to represent list of objects of
// type 'autoscaler_resource_limits'.
const AutoscalerResourceLimitsListKind = "AutoscalerResourceLimitsList"

// AutoscalerResourceLimitsListLinkKind is the name of the type used to represent links to list
// of objects of type 'autoscaler_resource_limits'.
const AutoscalerResourceLimitsListLinkKind = "AutoscalerResourceLimitsListLink"

// AutoscalerResourceLimitsNilKind is the name of the type used to nil lists of objects of
// type 'autoscaler_resource_limits'.
const AutoscalerResourceLimitsListNilKind = "AutoscalerResourceLimitsListNil"

// AutoscalerResourceLimitsList is a list of values of the 'autoscaler_resource_limits' type.
type AutoscalerResourceLimitsList struct {
	href  string
	link  bool
	items []*AutoscalerResourceLimits
}

// Len returns the length of the list.
func (l *AutoscalerResourceLimitsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AutoscalerResourceLimitsList) SetItems(items []*AutoscalerResourceLimits) {
	l.items = items
}

// Items returns the items of the list.
func (l *AutoscalerResourceLimitsList) Items() []*AutoscalerResourceLimits {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AutoscalerResourceLimitsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AutoscalerResourceLimitsList) Get(i int) *AutoscalerResourceLimits {
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
func (l *AutoscalerResourceLimitsList) Slice() []*AutoscalerResourceLimits {
	var slice []*AutoscalerResourceLimits
	if l == nil {
		slice = make([]*AutoscalerResourceLimits, 0)
	} else {
		slice = make([]*AutoscalerResourceLimits, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AutoscalerResourceLimitsList) Each(f func(item *AutoscalerResourceLimits) bool) {
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
func (l *AutoscalerResourceLimitsList) Range(f func(index int, item *AutoscalerResourceLimits) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
