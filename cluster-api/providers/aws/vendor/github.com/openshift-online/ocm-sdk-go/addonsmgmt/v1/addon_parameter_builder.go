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

// AddonParameterBuilder contains the data and logic needed to build 'addon_parameter' objects.
//
// Representation of an addon parameter.
type AddonParameterBuilder struct {
	bitmap_           uint32
	id                string
	addon             *AddonBuilder
	conditions        []*AddonRequirementBuilder
	defaultValue      string
	description       string
	editableDirection string
	name              string
	options           []*AddonParameterOptionBuilder
	order             int
	validation        string
	validationErrMsg  string
	valueType         AddonParameterValueType
	editable          bool
	enabled           bool
	required          bool
}

// NewAddonParameter creates a new builder of 'addon_parameter' objects.
func NewAddonParameter() *AddonParameterBuilder {
	return &AddonParameterBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonParameterBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AddonParameterBuilder) ID(value string) *AddonParameterBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an addon that can be installed in a cluster.
func (b *AddonParameterBuilder) Addon(value *AddonBuilder) *AddonParameterBuilder {
	b.addon = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Conditions sets the value of the 'conditions' attribute to the given values.
func (b *AddonParameterBuilder) Conditions(values ...*AddonRequirementBuilder) *AddonParameterBuilder {
	b.conditions = make([]*AddonRequirementBuilder, len(values))
	copy(b.conditions, values)
	b.bitmap_ |= 4
	return b
}

// DefaultValue sets the value of the 'default_value' attribute to the given value.
func (b *AddonParameterBuilder) DefaultValue(value string) *AddonParameterBuilder {
	b.defaultValue = value
	b.bitmap_ |= 8
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddonParameterBuilder) Description(value string) *AddonParameterBuilder {
	b.description = value
	b.bitmap_ |= 16
	return b
}

// Editable sets the value of the 'editable' attribute to the given value.
func (b *AddonParameterBuilder) Editable(value bool) *AddonParameterBuilder {
	b.editable = value
	b.bitmap_ |= 32
	return b
}

// EditableDirection sets the value of the 'editable_direction' attribute to the given value.
func (b *AddonParameterBuilder) EditableDirection(value string) *AddonParameterBuilder {
	b.editableDirection = value
	b.bitmap_ |= 64
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonParameterBuilder) Enabled(value bool) *AddonParameterBuilder {
	b.enabled = value
	b.bitmap_ |= 128
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonParameterBuilder) Name(value string) *AddonParameterBuilder {
	b.name = value
	b.bitmap_ |= 256
	return b
}

// Options sets the value of the 'options' attribute to the given values.
func (b *AddonParameterBuilder) Options(values ...*AddonParameterOptionBuilder) *AddonParameterBuilder {
	b.options = make([]*AddonParameterOptionBuilder, len(values))
	copy(b.options, values)
	b.bitmap_ |= 512
	return b
}

// Order sets the value of the 'order' attribute to the given value.
func (b *AddonParameterBuilder) Order(value int) *AddonParameterBuilder {
	b.order = value
	b.bitmap_ |= 1024
	return b
}

// Required sets the value of the 'required' attribute to the given value.
func (b *AddonParameterBuilder) Required(value bool) *AddonParameterBuilder {
	b.required = value
	b.bitmap_ |= 2048
	return b
}

// Validation sets the value of the 'validation' attribute to the given value.
func (b *AddonParameterBuilder) Validation(value string) *AddonParameterBuilder {
	b.validation = value
	b.bitmap_ |= 4096
	return b
}

// ValidationErrMsg sets the value of the 'validation_err_msg' attribute to the given value.
func (b *AddonParameterBuilder) ValidationErrMsg(value string) *AddonParameterBuilder {
	b.validationErrMsg = value
	b.bitmap_ |= 8192
	return b
}

// ValueType sets the value of the 'value_type' attribute to the given value.
//
// Representation of the value type for this specific addon parameter
func (b *AddonParameterBuilder) ValueType(value AddonParameterValueType) *AddonParameterBuilder {
	b.valueType = value
	b.bitmap_ |= 16384
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonParameterBuilder) Copy(object *AddonParameter) *AddonParameterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	if object.addon != nil {
		b.addon = NewAddon().Copy(object.addon)
	} else {
		b.addon = nil
	}
	if object.conditions != nil {
		b.conditions = make([]*AddonRequirementBuilder, len(object.conditions))
		for i, v := range object.conditions {
			b.conditions[i] = NewAddonRequirement().Copy(v)
		}
	} else {
		b.conditions = nil
	}
	b.defaultValue = object.defaultValue
	b.description = object.description
	b.editable = object.editable
	b.editableDirection = object.editableDirection
	b.enabled = object.enabled
	b.name = object.name
	if object.options != nil {
		b.options = make([]*AddonParameterOptionBuilder, len(object.options))
		for i, v := range object.options {
			b.options[i] = NewAddonParameterOption().Copy(v)
		}
	} else {
		b.options = nil
	}
	b.order = object.order
	b.required = object.required
	b.validation = object.validation
	b.validationErrMsg = object.validationErrMsg
	b.valueType = object.valueType
	return b
}

// Build creates a 'addon_parameter' object using the configuration stored in the builder.
func (b *AddonParameterBuilder) Build() (object *AddonParameter, err error) {
	object = new(AddonParameter)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	if b.addon != nil {
		object.addon, err = b.addon.Build()
		if err != nil {
			return
		}
	}
	if b.conditions != nil {
		object.conditions = make([]*AddonRequirement, len(b.conditions))
		for i, v := range b.conditions {
			object.conditions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.defaultValue = b.defaultValue
	object.description = b.description
	object.editable = b.editable
	object.editableDirection = b.editableDirection
	object.enabled = b.enabled
	object.name = b.name
	if b.options != nil {
		object.options = make([]*AddonParameterOption, len(b.options))
		for i, v := range b.options {
			object.options[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.order = b.order
	object.required = b.required
	object.validation = b.validation
	object.validationErrMsg = b.validationErrMsg
	object.valueType = b.valueType
	return
}
