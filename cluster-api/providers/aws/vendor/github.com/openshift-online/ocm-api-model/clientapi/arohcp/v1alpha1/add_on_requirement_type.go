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

// AddOnRequirement represents the values of the 'add_on_requirement' type.
//
// Representation of an add-on requirement.
type AddOnRequirement struct {
	fieldSet_ []bool
	id        string
	data      map[string]interface{}
	resource  string
	status    *AddOnRequirementStatus
	enabled   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnRequirement) Empty() bool {
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

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the add-on requirement.
func (o *AddOnRequirement) ID() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the add-on requirement.
func (o *AddOnRequirement) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// Data returns the value of the 'data' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Data for the add-on requirement.
func (o *AddOnRequirement) Data() map[string]interface{} {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.data
	}
	return nil
}

// GetData returns the value of the 'data' attribute and
// a flag indicating if the attribute has a value.
//
// Data for the add-on requirement.
func (o *AddOnRequirement) GetData() (value map[string]interface{}, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.data
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this requirement is enabled for the add-on.
func (o *AddOnRequirement) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this requirement is enabled for the add-on.
func (o *AddOnRequirement) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.enabled
	}
	return
}

// Resource returns the value of the 'resource' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of resource of the add-on requirement.
func (o *AddOnRequirement) Resource() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.resource
	}
	return ""
}

// GetResource returns the value of the 'resource' attribute and
// a flag indicating if the attribute has a value.
//
// Type of resource of the add-on requirement.
func (o *AddOnRequirement) GetResource() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.resource
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional cluster specific status for the add-on.
func (o *AddOnRequirement) Status() *AddOnRequirementStatus {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Optional cluster specific status for the add-on.
func (o *AddOnRequirement) GetStatus() (value *AddOnRequirementStatus, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.status
	}
	return
}

// AddOnRequirementListKind is the name of the type used to represent list of objects of
// type 'add_on_requirement'.
const AddOnRequirementListKind = "AddOnRequirementList"

// AddOnRequirementListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_requirement'.
const AddOnRequirementListLinkKind = "AddOnRequirementListLink"

// AddOnRequirementNilKind is the name of the type used to nil lists of objects of
// type 'add_on_requirement'.
const AddOnRequirementListNilKind = "AddOnRequirementListNil"

// AddOnRequirementList is a list of values of the 'add_on_requirement' type.
type AddOnRequirementList struct {
	href  string
	link  bool
	items []*AddOnRequirement
}

// Len returns the length of the list.
func (l *AddOnRequirementList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnRequirementList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnRequirementList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnRequirementList) SetItems(items []*AddOnRequirement) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnRequirementList) Items() []*AddOnRequirement {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnRequirementList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnRequirementList) Get(i int) *AddOnRequirement {
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
func (l *AddOnRequirementList) Slice() []*AddOnRequirement {
	var slice []*AddOnRequirement
	if l == nil {
		slice = make([]*AddOnRequirement, 0)
	} else {
		slice = make([]*AddOnRequirement, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnRequirementList) Each(f func(item *AddOnRequirement) bool) {
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
func (l *AddOnRequirementList) Range(f func(index int, item *AddOnRequirement) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
