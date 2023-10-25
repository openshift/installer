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

// AddonEnvironmentVariable represents the values of the 'addon_environment_variable' type.
//
// Representation of an addon env object.
type AddonEnvironmentVariable struct {
	bitmap_ uint32
	id      string
	name    string
	value   string
	enabled bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonEnvironmentVariable) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID for the environment variable
func (o *AddonEnvironmentVariable) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID for the environment variable
func (o *AddonEnvironmentVariable) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates is this environment variable is enabled for the addon
func (o *AddonEnvironmentVariable) Enabled() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates is this environment variable is enabled for the addon
func (o *AddonEnvironmentVariable) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.enabled
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the environment variable
func (o *AddonEnvironmentVariable) Name() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the environment variable
func (o *AddonEnvironmentVariable) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.name
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Value of the environment variable
func (o *AddonEnvironmentVariable) Value() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Value of the environment variable
func (o *AddonEnvironmentVariable) GetValue() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.value
	}
	return
}

// AddonEnvironmentVariableListKind is the name of the type used to represent list of objects of
// type 'addon_environment_variable'.
const AddonEnvironmentVariableListKind = "AddonEnvironmentVariableList"

// AddonEnvironmentVariableListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_environment_variable'.
const AddonEnvironmentVariableListLinkKind = "AddonEnvironmentVariableListLink"

// AddonEnvironmentVariableNilKind is the name of the type used to nil lists of objects of
// type 'addon_environment_variable'.
const AddonEnvironmentVariableListNilKind = "AddonEnvironmentVariableListNil"

// AddonEnvironmentVariableList is a list of values of the 'addon_environment_variable' type.
type AddonEnvironmentVariableList struct {
	href  string
	link  bool
	items []*AddonEnvironmentVariable
}

// Len returns the length of the list.
func (l *AddonEnvironmentVariableList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonEnvironmentVariableList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonEnvironmentVariableList) Get(i int) *AddonEnvironmentVariable {
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
func (l *AddonEnvironmentVariableList) Slice() []*AddonEnvironmentVariable {
	var slice []*AddonEnvironmentVariable
	if l == nil {
		slice = make([]*AddonEnvironmentVariable, 0)
	} else {
		slice = make([]*AddonEnvironmentVariable, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonEnvironmentVariableList) Each(f func(item *AddonEnvironmentVariable) bool) {
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
func (l *AddonEnvironmentVariableList) Range(f func(index int, item *AddonEnvironmentVariable) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
