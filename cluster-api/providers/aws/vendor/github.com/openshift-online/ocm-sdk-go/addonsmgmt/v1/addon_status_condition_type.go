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

// AddonStatusCondition represents the values of the 'addon_status_condition' type.
//
// Representation of an addon status condition type.
type AddonStatusCondition struct {
	bitmap_     uint32
	message     string
	reason      string
	statusType  AddonStatusConditionType
	statusValue AddonStatusConditionValue
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonStatusCondition) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Message returns the value of the 'message' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Message for the condition
func (o *AddonStatusCondition) Message() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.message
	}
	return ""
}

// GetMessage returns the value of the 'message' attribute and
// a flag indicating if the attribute has a value.
//
// Message for the condition
func (o *AddonStatusCondition) GetMessage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.message
	}
	return
}

// Reason returns the value of the 'reason' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Reason for the condition
func (o *AddonStatusCondition) Reason() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.reason
	}
	return ""
}

// GetReason returns the value of the 'reason' attribute and
// a flag indicating if the attribute has a value.
//
// Reason for the condition
func (o *AddonStatusCondition) GetReason() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.reason
	}
	return
}

// StatusType returns the value of the 'status_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of the reported addon status condition
func (o *AddonStatusCondition) StatusType() AddonStatusConditionType {
	if o != nil && o.bitmap_&4 != 0 {
		return o.statusType
	}
	return AddonStatusConditionType("")
}

// GetStatusType returns the value of the 'status_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of the reported addon status condition
func (o *AddonStatusCondition) GetStatusType() (value AddonStatusConditionType, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.statusType
	}
	return
}

// StatusValue returns the value of the 'status_value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Value of the reported addon status condition
func (o *AddonStatusCondition) StatusValue() AddonStatusConditionValue {
	if o != nil && o.bitmap_&8 != 0 {
		return o.statusValue
	}
	return AddonStatusConditionValue("")
}

// GetStatusValue returns the value of the 'status_value' attribute and
// a flag indicating if the attribute has a value.
//
// Value of the reported addon status condition
func (o *AddonStatusCondition) GetStatusValue() (value AddonStatusConditionValue, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.statusValue
	}
	return
}

// AddonStatusConditionListKind is the name of the type used to represent list of objects of
// type 'addon_status_condition'.
const AddonStatusConditionListKind = "AddonStatusConditionList"

// AddonStatusConditionListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_status_condition'.
const AddonStatusConditionListLinkKind = "AddonStatusConditionListLink"

// AddonStatusConditionNilKind is the name of the type used to nil lists of objects of
// type 'addon_status_condition'.
const AddonStatusConditionListNilKind = "AddonStatusConditionListNil"

// AddonStatusConditionList is a list of values of the 'addon_status_condition' type.
type AddonStatusConditionList struct {
	href  string
	link  bool
	items []*AddonStatusCondition
}

// Len returns the length of the list.
func (l *AddonStatusConditionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonStatusConditionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonStatusConditionList) Get(i int) *AddonStatusCondition {
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
func (l *AddonStatusConditionList) Slice() []*AddonStatusCondition {
	var slice []*AddonStatusCondition
	if l == nil {
		slice = make([]*AddonStatusCondition, 0)
	} else {
		slice = make([]*AddonStatusCondition, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonStatusConditionList) Each(f func(item *AddonStatusCondition) bool) {
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
func (l *AddonStatusConditionList) Range(f func(index int, item *AddonStatusCondition) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
