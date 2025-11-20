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

// ResourceRange represents the values of the 'resource_range' type.
type ResourceRange struct {
	fieldSet_ []bool
	max       int
	min       int
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ResourceRange) Empty() bool {
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

// Max returns the value of the 'max' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ResourceRange) Max() int {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.max
	}
	return 0
}

// GetMax returns the value of the 'max' attribute and
// a flag indicating if the attribute has a value.
func (o *ResourceRange) GetMax() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.max
	}
	return
}

// Min returns the value of the 'min' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ResourceRange) Min() int {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.min
	}
	return 0
}

// GetMin returns the value of the 'min' attribute and
// a flag indicating if the attribute has a value.
func (o *ResourceRange) GetMin() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.min
	}
	return
}

// ResourceRangeListKind is the name of the type used to represent list of objects of
// type 'resource_range'.
const ResourceRangeListKind = "ResourceRangeList"

// ResourceRangeListLinkKind is the name of the type used to represent links to list
// of objects of type 'resource_range'.
const ResourceRangeListLinkKind = "ResourceRangeListLink"

// ResourceRangeNilKind is the name of the type used to nil lists of objects of
// type 'resource_range'.
const ResourceRangeListNilKind = "ResourceRangeListNil"

// ResourceRangeList is a list of values of the 'resource_range' type.
type ResourceRangeList struct {
	href  string
	link  bool
	items []*ResourceRange
}

// Len returns the length of the list.
func (l *ResourceRangeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ResourceRangeList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ResourceRangeList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ResourceRangeList) SetItems(items []*ResourceRange) {
	l.items = items
}

// Items returns the items of the list.
func (l *ResourceRangeList) Items() []*ResourceRange {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ResourceRangeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ResourceRangeList) Get(i int) *ResourceRange {
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
func (l *ResourceRangeList) Slice() []*ResourceRange {
	var slice []*ResourceRange
	if l == nil {
		slice = make([]*ResourceRange, 0)
	} else {
		slice = make([]*ResourceRange, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ResourceRangeList) Each(f func(item *ResourceRange) bool) {
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
func (l *ResourceRangeList) Range(f func(index int, item *ResourceRange) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
