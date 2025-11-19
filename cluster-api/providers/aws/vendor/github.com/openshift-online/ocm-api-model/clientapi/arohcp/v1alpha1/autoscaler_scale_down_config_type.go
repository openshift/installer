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

// AutoscalerScaleDownConfig represents the values of the 'autoscaler_scale_down_config' type.
type AutoscalerScaleDownConfig struct {
	fieldSet_            []bool
	delayAfterAdd        string
	delayAfterDelete     string
	delayAfterFailure    string
	unneededTime         string
	utilizationThreshold string
	enabled              bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AutoscalerScaleDownConfig) Empty() bool {
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

// DelayAfterAdd returns the value of the 'delay_after_add' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long after scale up that scale down evaluation resumes.
func (o *AutoscalerScaleDownConfig) DelayAfterAdd() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.delayAfterAdd
	}
	return ""
}

// GetDelayAfterAdd returns the value of the 'delay_after_add' attribute and
// a flag indicating if the attribute has a value.
//
// How long after scale up that scale down evaluation resumes.
func (o *AutoscalerScaleDownConfig) GetDelayAfterAdd() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.delayAfterAdd
	}
	return
}

// DelayAfterDelete returns the value of the 'delay_after_delete' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long after node deletion that scale down evaluation resumes, defaults to scan-interval.
func (o *AutoscalerScaleDownConfig) DelayAfterDelete() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.delayAfterDelete
	}
	return ""
}

// GetDelayAfterDelete returns the value of the 'delay_after_delete' attribute and
// a flag indicating if the attribute has a value.
//
// How long after node deletion that scale down evaluation resumes, defaults to scan-interval.
func (o *AutoscalerScaleDownConfig) GetDelayAfterDelete() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.delayAfterDelete
	}
	return
}

// DelayAfterFailure returns the value of the 'delay_after_failure' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long after scale down failure that scale down evaluation resumes.
func (o *AutoscalerScaleDownConfig) DelayAfterFailure() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.delayAfterFailure
	}
	return ""
}

// GetDelayAfterFailure returns the value of the 'delay_after_failure' attribute and
// a flag indicating if the attribute has a value.
//
// How long after scale down failure that scale down evaluation resumes.
func (o *AutoscalerScaleDownConfig) GetDelayAfterFailure() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.delayAfterFailure
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Should cluster-autoscaler scale down the cluster.
func (o *AutoscalerScaleDownConfig) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Should cluster-autoscaler scale down the cluster.
func (o *AutoscalerScaleDownConfig) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.enabled
	}
	return
}

// UnneededTime returns the value of the 'unneeded_time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// How long a node should be unneeded before it is eligible for scale down.
func (o *AutoscalerScaleDownConfig) UnneededTime() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.unneededTime
	}
	return ""
}

// GetUnneededTime returns the value of the 'unneeded_time' attribute and
// a flag indicating if the attribute has a value.
//
// How long a node should be unneeded before it is eligible for scale down.
func (o *AutoscalerScaleDownConfig) GetUnneededTime() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.unneededTime
	}
	return
}

// UtilizationThreshold returns the value of the 'utilization_threshold' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down.
func (o *AutoscalerScaleDownConfig) UtilizationThreshold() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.utilizationThreshold
	}
	return ""
}

// GetUtilizationThreshold returns the value of the 'utilization_threshold' attribute and
// a flag indicating if the attribute has a value.
//
// Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down.
func (o *AutoscalerScaleDownConfig) GetUtilizationThreshold() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.utilizationThreshold
	}
	return
}

// AutoscalerScaleDownConfigListKind is the name of the type used to represent list of objects of
// type 'autoscaler_scale_down_config'.
const AutoscalerScaleDownConfigListKind = "AutoscalerScaleDownConfigList"

// AutoscalerScaleDownConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'autoscaler_scale_down_config'.
const AutoscalerScaleDownConfigListLinkKind = "AutoscalerScaleDownConfigListLink"

// AutoscalerScaleDownConfigNilKind is the name of the type used to nil lists of objects of
// type 'autoscaler_scale_down_config'.
const AutoscalerScaleDownConfigListNilKind = "AutoscalerScaleDownConfigListNil"

// AutoscalerScaleDownConfigList is a list of values of the 'autoscaler_scale_down_config' type.
type AutoscalerScaleDownConfigList struct {
	href  string
	link  bool
	items []*AutoscalerScaleDownConfig
}

// Len returns the length of the list.
func (l *AutoscalerScaleDownConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AutoscalerScaleDownConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AutoscalerScaleDownConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AutoscalerScaleDownConfigList) SetItems(items []*AutoscalerScaleDownConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *AutoscalerScaleDownConfigList) Items() []*AutoscalerScaleDownConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AutoscalerScaleDownConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AutoscalerScaleDownConfigList) Get(i int) *AutoscalerScaleDownConfig {
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
func (l *AutoscalerScaleDownConfigList) Slice() []*AutoscalerScaleDownConfig {
	var slice []*AutoscalerScaleDownConfig
	if l == nil {
		slice = make([]*AutoscalerScaleDownConfig, 0)
	} else {
		slice = make([]*AutoscalerScaleDownConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AutoscalerScaleDownConfigList) Each(f func(item *AutoscalerScaleDownConfig) bool) {
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
func (l *AutoscalerScaleDownConfigList) Range(f func(index int, item *AutoscalerScaleDownConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
