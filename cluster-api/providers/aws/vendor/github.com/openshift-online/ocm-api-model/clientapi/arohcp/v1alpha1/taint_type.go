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

// Taint represents the values of the 'taint' type.
//
// Representation of a Taint set on a MachinePool in a cluster.
type Taint struct {
	fieldSet_ []bool
	effect    string
	key       string
	value     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Taint) Empty() bool {
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

// Effect returns the value of the 'effect' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The effect on the node for the pods matching the taint, i.e: NoSchedule, NoExecute, PreferNoSchedule.
func (o *Taint) Effect() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.effect
	}
	return ""
}

// GetEffect returns the value of the 'effect' attribute and
// a flag indicating if the attribute has a value.
//
// The effect on the node for the pods matching the taint, i.e: NoSchedule, NoExecute, PreferNoSchedule.
func (o *Taint) GetEffect() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.effect
	}
	return
}

// Key returns the value of the 'key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The key for the taint
func (o *Taint) Key() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.key
	}
	return ""
}

// GetKey returns the value of the 'key' attribute and
// a flag indicating if the attribute has a value.
//
// The key for the taint
func (o *Taint) GetKey() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.key
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The value for the taint.
func (o *Taint) Value() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// The value for the taint.
func (o *Taint) GetValue() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.value
	}
	return
}

// TaintListKind is the name of the type used to represent list of objects of
// type 'taint'.
const TaintListKind = "TaintList"

// TaintListLinkKind is the name of the type used to represent links to list
// of objects of type 'taint'.
const TaintListLinkKind = "TaintListLink"

// TaintNilKind is the name of the type used to nil lists of objects of
// type 'taint'.
const TaintListNilKind = "TaintListNil"

// TaintList is a list of values of the 'taint' type.
type TaintList struct {
	href  string
	link  bool
	items []*Taint
}

// Len returns the length of the list.
func (l *TaintList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TaintList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TaintList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TaintList) SetItems(items []*Taint) {
	l.items = items
}

// Items returns the items of the list.
func (l *TaintList) Items() []*Taint {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TaintList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TaintList) Get(i int) *Taint {
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
func (l *TaintList) Slice() []*Taint {
	var slice []*Taint
	if l == nil {
		slice = make([]*Taint, 0)
	} else {
		slice = make([]*Taint, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TaintList) Each(f func(item *Taint) bool) {
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
func (l *TaintList) Range(f func(index int, item *Taint) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
