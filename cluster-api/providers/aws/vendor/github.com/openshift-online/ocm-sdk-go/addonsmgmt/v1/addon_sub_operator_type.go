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

// AddonSubOperator represents the values of the 'addon_sub_operator' type.
//
// Representation of an addon sub operator. A sub operator is an operator
// who's life cycle is controlled by the addon umbrella operator.
type AddonSubOperator struct {
	bitmap_           uint32
	operatorName      string
	operatorNamespace string
	enabled           bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonSubOperator) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if the sub operator is enabled for the addon
func (o *AddonSubOperator) Enabled() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if the sub operator is enabled for the addon
func (o *AddonSubOperator) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.enabled
	}
	return
}

// OperatorName returns the value of the 'operator_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the addon sub operator
func (o *AddonSubOperator) OperatorName() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.operatorName
	}
	return ""
}

// GetOperatorName returns the value of the 'operator_name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the addon sub operator
func (o *AddonSubOperator) GetOperatorName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.operatorName
	}
	return
}

// OperatorNamespace returns the value of the 'operator_namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespace of the addon sub operator
func (o *AddonSubOperator) OperatorNamespace() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.operatorNamespace
	}
	return ""
}

// GetOperatorNamespace returns the value of the 'operator_namespace' attribute and
// a flag indicating if the attribute has a value.
//
// Namespace of the addon sub operator
func (o *AddonSubOperator) GetOperatorNamespace() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.operatorNamespace
	}
	return
}

// AddonSubOperatorListKind is the name of the type used to represent list of objects of
// type 'addon_sub_operator'.
const AddonSubOperatorListKind = "AddonSubOperatorList"

// AddonSubOperatorListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_sub_operator'.
const AddonSubOperatorListLinkKind = "AddonSubOperatorListLink"

// AddonSubOperatorNilKind is the name of the type used to nil lists of objects of
// type 'addon_sub_operator'.
const AddonSubOperatorListNilKind = "AddonSubOperatorListNil"

// AddonSubOperatorList is a list of values of the 'addon_sub_operator' type.
type AddonSubOperatorList struct {
	href  string
	link  bool
	items []*AddonSubOperator
}

// Len returns the length of the list.
func (l *AddonSubOperatorList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonSubOperatorList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonSubOperatorList) Get(i int) *AddonSubOperator {
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
func (l *AddonSubOperatorList) Slice() []*AddonSubOperator {
	var slice []*AddonSubOperator
	if l == nil {
		slice = make([]*AddonSubOperator, 0)
	} else {
		slice = make([]*AddonSubOperator, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonSubOperatorList) Each(f func(item *AddonSubOperator) bool) {
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
func (l *AddonSubOperatorList) Range(f func(index int, item *AddonSubOperator) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
