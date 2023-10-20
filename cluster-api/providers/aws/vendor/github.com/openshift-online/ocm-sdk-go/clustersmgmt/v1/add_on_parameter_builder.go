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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AddOnParameterBuilder contains the data and logic needed to build 'add_on_parameter' objects.
//
// Representation of an add-on parameter.
type AddOnParameterBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	addon             *AddOnBuilder
	conditions        []*AddOnRequirementBuilder
	defaultValue      string
	description       string
	editableDirection string
	name              string
	options           []*AddOnParameterOptionBuilder
	validation        string
	validationErrMsg  string
	valueType         string
	editable          bool
	enabled           bool
	required          bool
}

// NewAddOnParameter creates a new builder of 'add_on_parameter' objects.
func NewAddOnParameter() *AddOnParameterBuilder {
	return &AddOnParameterBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnParameterBuilder) Link(value bool) *AddOnParameterBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnParameterBuilder) ID(value string) *AddOnParameterBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnParameterBuilder) HREF(value string) *AddOnParameterBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnParameterBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an add-on that can be installed in a cluster.
func (b *AddOnParameterBuilder) Addon(value *AddOnBuilder) *AddOnParameterBuilder {
	b.addon = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Conditions sets the value of the 'conditions' attribute to the given values.
func (b *AddOnParameterBuilder) Conditions(values ...*AddOnRequirementBuilder) *AddOnParameterBuilder {
	b.conditions = make([]*AddOnRequirementBuilder, len(values))
	copy(b.conditions, values)
	b.bitmap_ |= 16
	return b
}

// DefaultValue sets the value of the 'default_value' attribute to the given value.
func (b *AddOnParameterBuilder) DefaultValue(value string) *AddOnParameterBuilder {
	b.defaultValue = value
	b.bitmap_ |= 32
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddOnParameterBuilder) Description(value string) *AddOnParameterBuilder {
	b.description = value
	b.bitmap_ |= 64
	return b
}

// Editable sets the value of the 'editable' attribute to the given value.
func (b *AddOnParameterBuilder) Editable(value bool) *AddOnParameterBuilder {
	b.editable = value
	b.bitmap_ |= 128
	return b
}

// EditableDirection sets the value of the 'editable_direction' attribute to the given value.
func (b *AddOnParameterBuilder) EditableDirection(value string) *AddOnParameterBuilder {
	b.editableDirection = value
	b.bitmap_ |= 256
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnParameterBuilder) Enabled(value bool) *AddOnParameterBuilder {
	b.enabled = value
	b.bitmap_ |= 512
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnParameterBuilder) Name(value string) *AddOnParameterBuilder {
	b.name = value
	b.bitmap_ |= 1024
	return b
}

// Options sets the value of the 'options' attribute to the given values.
func (b *AddOnParameterBuilder) Options(values ...*AddOnParameterOptionBuilder) *AddOnParameterBuilder {
	b.options = make([]*AddOnParameterOptionBuilder, len(values))
	copy(b.options, values)
	b.bitmap_ |= 2048
	return b
}

// Required sets the value of the 'required' attribute to the given value.
func (b *AddOnParameterBuilder) Required(value bool) *AddOnParameterBuilder {
	b.required = value
	b.bitmap_ |= 4096
	return b
}

// Validation sets the value of the 'validation' attribute to the given value.
func (b *AddOnParameterBuilder) Validation(value string) *AddOnParameterBuilder {
	b.validation = value
	b.bitmap_ |= 8192
	return b
}

// ValidationErrMsg sets the value of the 'validation_err_msg' attribute to the given value.
func (b *AddOnParameterBuilder) ValidationErrMsg(value string) *AddOnParameterBuilder {
	b.validationErrMsg = value
	b.bitmap_ |= 16384
	return b
}

// ValueType sets the value of the 'value_type' attribute to the given value.
func (b *AddOnParameterBuilder) ValueType(value string) *AddOnParameterBuilder {
	b.valueType = value
	b.bitmap_ |= 32768
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnParameterBuilder) Copy(object *AddOnParameter) *AddOnParameterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.addon != nil {
		b.addon = NewAddOn().Copy(object.addon)
	} else {
		b.addon = nil
	}
	if object.conditions != nil {
		b.conditions = make([]*AddOnRequirementBuilder, len(object.conditions))
		for i, v := range object.conditions {
			b.conditions[i] = NewAddOnRequirement().Copy(v)
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
		b.options = make([]*AddOnParameterOptionBuilder, len(object.options))
		for i, v := range object.options {
			b.options[i] = NewAddOnParameterOption().Copy(v)
		}
	} else {
		b.options = nil
	}
	b.required = object.required
	b.validation = object.validation
	b.validationErrMsg = object.validationErrMsg
	b.valueType = object.valueType
	return b
}

// Build creates a 'add_on_parameter' object using the configuration stored in the builder.
func (b *AddOnParameterBuilder) Build() (object *AddOnParameter, err error) {
	object = new(AddOnParameter)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.addon != nil {
		object.addon, err = b.addon.Build()
		if err != nil {
			return
		}
	}
	if b.conditions != nil {
		object.conditions = make([]*AddOnRequirement, len(b.conditions))
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
		object.options = make([]*AddOnParameterOption, len(b.options))
		for i, v := range b.options {
			object.options[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.required = b.required
	object.validation = b.validation
	object.validationErrMsg = b.validationErrMsg
	object.valueType = b.valueType
	return
}
