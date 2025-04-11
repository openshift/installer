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

// RolePolicyBinding represents the values of the 'role_policy_binding' type.
type RolePolicyBinding struct {
	bitmap_             uint32
	arn                 string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	name                string
	policies            []*RolePolicy
	status              *RolePolicyBindingStatus
	type_               string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RolePolicyBinding) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Arn returns the value of the 'arn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) Arn() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.arn
	}
	return ""
}

// GetArn returns the value of the 'arn' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetArn() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.arn
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&2 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// LastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) LastUpdateTimestamp() time.Time {
	if o != nil && o.bitmap_&4 != 0 {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetLastUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.lastUpdateTimestamp
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) Name() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.name
	}
	return
}

// Policies returns the value of the 'policies' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) Policies() []*RolePolicy {
	if o != nil && o.bitmap_&16 != 0 {
		return o.policies
	}
	return nil
}

// GetPolicies returns the value of the 'policies' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetPolicies() (value []*RolePolicy, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.policies
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) Status() *RolePolicyBindingStatus {
	if o != nil && o.bitmap_&32 != 0 {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetStatus() (value *RolePolicyBindingStatus, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.status
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicyBinding) Type() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicyBinding) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.type_
	}
	return
}

// RolePolicyBindingListKind is the name of the type used to represent list of objects of
// type 'role_policy_binding'.
const RolePolicyBindingListKind = "RolePolicyBindingList"

// RolePolicyBindingListLinkKind is the name of the type used to represent links to list
// of objects of type 'role_policy_binding'.
const RolePolicyBindingListLinkKind = "RolePolicyBindingListLink"

// RolePolicyBindingNilKind is the name of the type used to nil lists of objects of
// type 'role_policy_binding'.
const RolePolicyBindingListNilKind = "RolePolicyBindingListNil"

// RolePolicyBindingList is a list of values of the 'role_policy_binding' type.
type RolePolicyBindingList struct {
	href  string
	link  bool
	items []*RolePolicyBinding
}

// Len returns the length of the list.
func (l *RolePolicyBindingList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *RolePolicyBindingList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RolePolicyBindingList) Get(i int) *RolePolicyBinding {
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
func (l *RolePolicyBindingList) Slice() []*RolePolicyBinding {
	var slice []*RolePolicyBinding
	if l == nil {
		slice = make([]*RolePolicyBinding, 0)
	} else {
		slice = make([]*RolePolicyBinding, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RolePolicyBindingList) Each(f func(item *RolePolicyBinding) bool) {
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
func (l *RolePolicyBindingList) Range(f func(index int, item *RolePolicyBinding) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
