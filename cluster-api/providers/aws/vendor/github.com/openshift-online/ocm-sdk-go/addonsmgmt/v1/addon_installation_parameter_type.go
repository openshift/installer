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

// AddonInstallationParameter represents the values of the 'addon_installation_parameter' type.
//
// representation of addon installation parameter
type AddonInstallationParameter struct {
	bitmap_ uint32
	href    string
	id      string
	kind    string
	value   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonInstallationParameter) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Href returns the value of the 'href' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Self link
func (o *AddonInstallationParameter) Href() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.href
	}
	return ""
}

// GetHref returns the value of the 'href' attribute and
// a flag indicating if the attribute has a value.
//
// Self link
func (o *AddonInstallationParameter) GetHref() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.href
	}
	return
}

// Id returns the value of the 'id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Unique identifier of the object
func (o *AddonInstallationParameter) Id() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetId returns the value of the 'id' attribute and
// a flag indicating if the attribute has a value.
//
// Unique identifier of the object
func (o *AddonInstallationParameter) GetId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// Kind returns the value of the 'kind' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the type of this object
func (o *AddonInstallationParameter) Kind() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.kind
	}
	return ""
}

// GetKind returns the value of the 'kind' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the type of this object
func (o *AddonInstallationParameter) GetKind() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.kind
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Value of the parameter
func (o *AddonInstallationParameter) Value() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Value of the parameter
func (o *AddonInstallationParameter) GetValue() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.value
	}
	return
}

// AddonInstallationParameterListKind is the name of the type used to represent list of objects of
// type 'addon_installation_parameter'.
const AddonInstallationParameterListKind = "AddonInstallationParameterList"

// AddonInstallationParameterListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_installation_parameter'.
const AddonInstallationParameterListLinkKind = "AddonInstallationParameterListLink"

// AddonInstallationParameterNilKind is the name of the type used to nil lists of objects of
// type 'addon_installation_parameter'.
const AddonInstallationParameterListNilKind = "AddonInstallationParameterListNil"

// AddonInstallationParameterList is a list of values of the 'addon_installation_parameter' type.
type AddonInstallationParameterList struct {
	href  string
	link  bool
	items []*AddonInstallationParameter
}

// Len returns the length of the list.
func (l *AddonInstallationParameterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonInstallationParameterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonInstallationParameterList) Get(i int) *AddonInstallationParameter {
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
func (l *AddonInstallationParameterList) Slice() []*AddonInstallationParameter {
	var slice []*AddonInstallationParameter
	if l == nil {
		slice = make([]*AddonInstallationParameter, 0)
	} else {
		slice = make([]*AddonInstallationParameter, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonInstallationParameterList) Each(f func(item *AddonInstallationParameter) bool) {
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
func (l *AddonInstallationParameterList) Range(f func(index int, item *AddonInstallationParameter) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
