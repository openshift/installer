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

// AddonParameter represents the values of the 'addon_parameter' type.
//
// Representation of an addon parameter.
type AddonParameter struct {
	bitmap_           uint32
	id                string
	addon             *Addon
	conditions        []*AddonRequirement
	defaultValue      string
	description       string
	editableDirection string
	name              string
	options           []*AddonParameterOption
	order             int
	validation        string
	validationErrMsg  string
	valueType         AddonParameterValueType
	editable          bool
	enabled           bool
	required          bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddonParameter) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID for this addon parameter
func (o *AddonParameter) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID for this addon parameter
func (o *AddonParameter) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// Addon returns the value of the 'addon' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *AddonParameter) Addon() *Addon {
	if o != nil && o.bitmap_&2 != 0 {
		return o.addon
	}
	return nil
}

// GetAddon returns the value of the 'addon' attribute and
// a flag indicating if the attribute has a value.
func (o *AddonParameter) GetAddon() (value *Addon, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.addon
	}
	return
}

// Conditions returns the value of the 'conditions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Conditions in which this parameter is valid for
func (o *AddonParameter) Conditions() []*AddonRequirement {
	if o != nil && o.bitmap_&4 != 0 {
		return o.conditions
	}
	return nil
}

// GetConditions returns the value of the 'conditions' attribute and
// a flag indicating if the attribute has a value.
//
// Conditions in which this parameter is valid for
func (o *AddonParameter) GetConditions() (value []*AddonRequirement, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.conditions
	}
	return
}

// DefaultValue returns the value of the 'default_value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the value default for the addon parameter.
func (o *AddonParameter) DefaultValue() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.defaultValue
	}
	return ""
}

// GetDefaultValue returns the value of the 'default_value' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the value default for the addon parameter.
func (o *AddonParameter) GetDefaultValue() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.defaultValue
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the addon parameter.
func (o *AddonParameter) Description() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the addon parameter.
func (o *AddonParameter) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.description
	}
	return
}

// Editable returns the value of the 'editable' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this parameter can be edited after creation.
func (o *AddonParameter) Editable() bool {
	if o != nil && o.bitmap_&32 != 0 {
		return o.editable
	}
	return false
}

// GetEditable returns the value of the 'editable' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this parameter can be edited after creation.
func (o *AddonParameter) GetEditable() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.editable
	}
	return
}

// EditableDirection returns the value of the 'editable_direction' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Restricts if the parameter can be upscaled/downscaled
// Expected values are "up", "down", or "" (no restriction).
func (o *AddonParameter) EditableDirection() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.editableDirection
	}
	return ""
}

// GetEditableDirection returns the value of the 'editable_direction' attribute and
// a flag indicating if the attribute has a value.
//
// Restricts if the parameter can be upscaled/downscaled
// Expected values are "up", "down", or "" (no restriction).
func (o *AddonParameter) GetEditableDirection() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.editableDirection
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this parameter is enabled for the addon.
func (o *AddonParameter) Enabled() bool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this parameter is enabled for the addon.
func (o *AddonParameter) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.enabled
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the addon parameter.
func (o *AddonParameter) Name() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the addon parameter.
func (o *AddonParameter) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.name
	}
	return
}

// Options returns the value of the 'options' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of options for the addon parameter value.
func (o *AddonParameter) Options() []*AddonParameterOption {
	if o != nil && o.bitmap_&512 != 0 {
		return o.options
	}
	return nil
}

// GetOptions returns the value of the 'options' attribute and
// a flag indicating if the attribute has a value.
//
// List of options for the addon parameter value.
func (o *AddonParameter) GetOptions() (value []*AddonParameterOption, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.options
	}
	return
}

// Order returns the value of the 'order' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the weight of the AddonParameter which would be used by sort order
func (o *AddonParameter) Order() int {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.order
	}
	return 0
}

// GetOrder returns the value of the 'order' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the weight of the AddonParameter which would be used by sort order
func (o *AddonParameter) GetOrder() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.order
	}
	return
}

// Required returns the value of the 'required' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this parameter is required by the addon.
func (o *AddonParameter) Required() bool {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.required
	}
	return false
}

// GetRequired returns the value of the 'required' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this parameter is required by the addon.
func (o *AddonParameter) GetRequired() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.required
	}
	return
}

// Validation returns the value of the 'validation' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Validation rule for the addon parameter.
func (o *AddonParameter) Validation() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.validation
	}
	return ""
}

// GetValidation returns the value of the 'validation' attribute and
// a flag indicating if the attribute has a value.
//
// Validation rule for the addon parameter.
func (o *AddonParameter) GetValidation() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.validation
	}
	return
}

// ValidationErrMsg returns the value of the 'validation_err_msg' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Error message to return should the parameter be invalid.
func (o *AddonParameter) ValidationErrMsg() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.validationErrMsg
	}
	return ""
}

// GetValidationErrMsg returns the value of the 'validation_err_msg' attribute and
// a flag indicating if the attribute has a value.
//
// Error message to return should the parameter be invalid.
func (o *AddonParameter) GetValidationErrMsg() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.validationErrMsg
	}
	return
}

// ValueType returns the value of the 'value_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of value of the addon parameter.
func (o *AddonParameter) ValueType() AddonParameterValueType {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.valueType
	}
	return AddonParameterValueType("")
}

// GetValueType returns the value of the 'value_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of value of the addon parameter.
func (o *AddonParameter) GetValueType() (value AddonParameterValueType, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.valueType
	}
	return
}

// AddonParameterListKind is the name of the type used to represent list of objects of
// type 'addon_parameter'.
const AddonParameterListKind = "AddonParameterList"

// AddonParameterListLinkKind is the name of the type used to represent links to list
// of objects of type 'addon_parameter'.
const AddonParameterListLinkKind = "AddonParameterListLink"

// AddonParameterNilKind is the name of the type used to nil lists of objects of
// type 'addon_parameter'.
const AddonParameterListNilKind = "AddonParameterListNil"

// AddonParameterList is a list of values of the 'addon_parameter' type.
type AddonParameterList struct {
	href  string
	link  bool
	items []*AddonParameter
}

// Len returns the length of the list.
func (l *AddonParameterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddonParameterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddonParameterList) Get(i int) *AddonParameter {
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
func (l *AddonParameterList) Slice() []*AddonParameter {
	var slice []*AddonParameter
	if l == nil {
		slice = make([]*AddonParameter, 0)
	} else {
		slice = make([]*AddonParameter, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddonParameterList) Each(f func(item *AddonParameter) bool) {
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
func (l *AddonParameterList) Range(f func(index int, item *AddonParameter) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
