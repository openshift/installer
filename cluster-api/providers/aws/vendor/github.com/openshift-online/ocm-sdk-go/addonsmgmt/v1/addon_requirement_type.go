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

// AddonRequirement represents the values of the 'addon_requirement' type.
//
// Representation of an addon requirement.
type AddonRequirement struct {
	bitmap_  uint32
	id       string
	data     map[string]interface{}
	resource AddonRequirementResource
	status   *AddonRequirementStatus
	enabled  bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonRequirement) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the addon requirement.
func (o *AddonRequirement) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the addon requirement.
func (o *AddonRequirement) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// Data returns the value of the 'data' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Data for the addon requirement.
func (o *AddonRequirement) Data() map[string]interface{} {
	if o != nil && o.bitmap_&2 != 0 {
		return o.data
	}
	return nil
}

// GetData returns the value of the 'data' attribute and
// a flag indicating if the attribute has a value.
//
// Data for the addon requirement.
func (o *AddonRequirement) GetData() (value map[string]interface{}, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.data
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this requirement is enabled for the addon.
func (o *AddonRequirement) Enabled() bool {
	if o != nil && o.bitmap_&4 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this requirement is enabled for the addon.
func (o *AddonRequirement) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.enabled
	}
	return
}

// Resource returns the value of the 'resource' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of resource of the addon requirement.
func (o *AddonRequirement) Resource() AddonRequirementResource {
	if o != nil && o.bitmap_&8 != 0 {
		return o.resource
	}
	return AddonRequirementResource("")
}

// GetResource returns the value of the 'resource' attribute and
// a flag indicating if the attribute has a value.
//
// Type of resource of the addon requirement.
func (o *AddonRequirement) GetResource() (value AddonRequirementResource, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.resource
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional cluster specific status for the addon.
func (o *AddonRequirement) Status() *AddonRequirementStatus {
	if o != nil && o.bitmap_&16 != 0 {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Optional cluster specific status for the addon.
func (o *AddonRequirement) GetStatus() (value *AddonRequirementStatus, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.status
	}
	return
}

// AddonRequirementListKind is the name of the type used to represent list of objects of
// type 'addon_requirement'.
const AddonRequirementListKind = "AddonRequirementList"

// AddonRequirementListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_requirement'.
const AddonRequirementListLinkKind = "AddonRequirementListLink"

// AddonRequirementNilKind is the name of the type used to nil lists of objects of
// type 'addon_requirement'.
const AddonRequirementListNilKind = "AddonRequirementListNil"

// AddonRequirementList is a list of values of the 'addon_requirement' type.
type AddonRequirementList struct {
	href  string
	link  bool
	items []*AddonRequirement
}

// Len returns the length of the list.
func (l *AddonRequirementList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonRequirementList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonRequirementList) Get(i int) *AddonRequirement {
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
func (l *AddonRequirementList) Slice() []*AddonRequirement {
	var slice []*AddonRequirement
	if l == nil {
		slice = make([]*AddonRequirement, 0)
	} else {
		slice = make([]*AddonRequirement, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonRequirementList) Each(f func(item *AddonRequirement) bool) {
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
func (l *AddonRequirementList) Range(f func(index int, item *AddonRequirement) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
