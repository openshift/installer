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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// ValueUnit represents the values of the 'value_unit' type.
type ValueUnit struct {
	fieldSet_ []bool
	unit      string
	value     float64
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ValueUnit) Empty() bool {
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
func (o *ValueUnit) Unit() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.unit
	}
	return ""
}

// GetUnit returns the value of the 'unit' attribute and
// a flag indicating if the attribute has a value.
func (o *ValueUnit) GetUnit() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.unit
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *ValueUnit) Value() float64 {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.value
	}
	return 0.0
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
func (o *ValueUnit) GetValue() (value float64, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.value
	}
	return
}

// ValueUnitListKind is the name of the type used to represent list of objects of
// type 'value_unit'.
const ValueUnitListKind = "ValueUnitList"

// ValueUnitListLinkKind is the name of the type used to represent links to list
// of objects of type 'value_unit'.
const ValueUnitListLinkKind = "ValueUnitListLink"

// ValueUnitNilKind is the name of the type used to nil lists of objects of
// type 'value_unit'.
const ValueUnitListNilKind = "ValueUnitListNil"

// ValueUnitList is a list of values of the 'value_unit' type.
type ValueUnitList struct {
	href  string
	link  bool
	items []*ValueUnit
}

// Len returns the length of the list.
func (l *ValueUnitList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ValueUnitList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ValueUnitList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ValueUnitList) SetItems(items []*ValueUnit) {
	l.items = items
}

// Items returns the items of the list.
func (l *ValueUnitList) Items() []*ValueUnit {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ValueUnitList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ValueUnitList) Get(i int) *ValueUnit {
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
func (l *ValueUnitList) Slice() []*ValueUnit {
	var slice []*ValueUnit
	if l == nil {
		slice = make([]*ValueUnit, 0)
	} else {
		slice = make([]*ValueUnit, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ValueUnitList) Each(f func(item *ValueUnit) bool) {
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
func (l *ValueUnitList) Range(f func(index int, item *ValueUnit) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
