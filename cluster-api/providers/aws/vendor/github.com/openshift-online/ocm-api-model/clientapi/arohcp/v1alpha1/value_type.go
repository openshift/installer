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

// Value represents the values of the 'value' type.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
type Value struct {
	fieldSet_ []bool
	unit      string
	value     float64
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Value) Empty() bool {
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

// Unit returns the value of the 'unit' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the unit used to measure the value.
func (o *Value) Unit() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.unit
	}
	return ""
}

// GetUnit returns the value of the 'unit' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the unit used to measure the value.
func (o *Value) GetUnit() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.unit
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Numeric value.
func (o *Value) Value() float64 {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.value
	}
	return 0.0
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Numeric value.
func (o *Value) GetValue() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.value
	}
	return
}

// ValueListKind is the name of the type used to represent list of objects of
// type 'value'.
const ValueListKind = "ValueList"

// ValueListLinkKind is the name of the type used to represent links to list
// of objects of type 'value'.
const ValueListLinkKind = "ValueListLink"

// ValueNilKind is the name of the type used to nil lists of objects of
// type 'value'.
const ValueListNilKind = "ValueListNil"

// ValueList is a list of values of the 'value' type.
type ValueList struct {
	href  string
	link  bool
	items []*Value
}

// Len returns the length of the list.
func (l *ValueList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ValueList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ValueList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ValueList) SetItems(items []*Value) {
	l.items = items
}

// Items returns the items of the list.
func (l *ValueList) Items() []*Value {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ValueList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ValueList) Get(i int) *Value {
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
func (l *ValueList) Slice() []*Value {
	var slice []*Value
	if l == nil {
		slice = make([]*Value, 0)
	} else {
		slice = make([]*Value, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ValueList) Each(f func(item *Value) bool) {
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
func (l *ValueList) Range(f func(index int, item *Value) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
