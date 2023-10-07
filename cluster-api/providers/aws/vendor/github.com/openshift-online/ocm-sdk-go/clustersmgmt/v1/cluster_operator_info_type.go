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

// ClusterOperatorInfo represents the values of the 'cluster_operator_info' type.
type ClusterOperatorInfo struct {
	bitmap_   uint32
	condition ClusterOperatorState
	name      string
	reason    string
	time      time.Time
	version   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterOperatorInfo) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Condition returns the value of the 'condition' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Operator status.  Empty string if unknown.
func (o *ClusterOperatorInfo) Condition() ClusterOperatorState {
	if o != nil && o.bitmap_&1 != 0 {
		return o.condition
	}
	return ClusterOperatorState("")
}

// GetCondition returns the value of the 'condition' attribute and
// a flag indicating if the attribute has a value.
//
// Operator status.  Empty string if unknown.
func (o *ClusterOperatorInfo) GetCondition() (value ClusterOperatorState, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.condition
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the operator.
func (o *ClusterOperatorInfo) Name() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the operator.
func (o *ClusterOperatorInfo) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.name
	}
	return
}

// Reason returns the value of the 'reason' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Extra detail on condition, if available.  Empty string if unknown.
func (o *ClusterOperatorInfo) Reason() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.reason
	}
	return ""
}

// GetReason returns the value of the 'reason' attribute and
// a flag indicating if the attribute has a value.
//
// Extra detail on condition, if available.  Empty string if unknown.
func (o *ClusterOperatorInfo) GetReason() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.reason
	}
	return
}

// Time returns the value of the 'time' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Time when the sample was obtained, in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
func (o *ClusterOperatorInfo) Time() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.time
	}
	return time.Time{}
}

// GetTime returns the value of the 'time' attribute and
// a flag indicating if the attribute has a value.
//
// Time when the sample was obtained, in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
func (o *ClusterOperatorInfo) GetTime() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.time
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Current version of the operator.  Empty string if unknown.
func (o *ClusterOperatorInfo) Version() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.version
	}
	return ""
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Current version of the operator.  Empty string if unknown.
func (o *ClusterOperatorInfo) GetVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.version
	}
	return
}

// ClusterOperatorInfoListKind is the name of the type used to represent list of objects of
// type 'cluster_operator_info'.
const ClusterOperatorInfoListKind = "ClusterOperatorInfoList"

// ClusterOperatorInfoListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_operator_info'.
const ClusterOperatorInfoListLinkKind = "ClusterOperatorInfoListLink"

// ClusterOperatorInfoNilKind is the name of the type used to nil lists of objects of
// type 'cluster_operator_info'.
const ClusterOperatorInfoListNilKind = "ClusterOperatorInfoListNil"

// ClusterOperatorInfoList is a list of values of the 'cluster_operator_info' type.
type ClusterOperatorInfoList struct {
	href  string
	link  bool
	items []*ClusterOperatorInfo
}

// Len returns the length of the list.
func (l *ClusterOperatorInfoList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ClusterOperatorInfoList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterOperatorInfoList) Get(i int) *ClusterOperatorInfo {
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
func (l *ClusterOperatorInfoList) Slice() []*ClusterOperatorInfo {
	var slice []*ClusterOperatorInfo
	if l == nil {
		slice = make([]*ClusterOperatorInfo, 0)
	} else {
		slice = make([]*ClusterOperatorInfo, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterOperatorInfoList) Each(f func(item *ClusterOperatorInfo) bool) {
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
func (l *ClusterOperatorInfoList) Range(f func(index int, item *ClusterOperatorInfo) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
