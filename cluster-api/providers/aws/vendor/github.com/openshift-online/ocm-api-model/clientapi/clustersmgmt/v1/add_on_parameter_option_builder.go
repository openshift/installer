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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of an add-on parameter option.
type AddOnParameterOptionBuilder struct {
	fieldSet_    []bool
	name         string
	rank         int
	requirements []*AddOnRequirementBuilder
	value        string
}

// NewAddOnParameterOption creates a new builder of 'add_on_parameter_option' objects.
func NewAddOnParameterOption() *AddOnParameterOptionBuilder {
	return &AddOnParameterOptionBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnParameterOptionBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Name(value string) *AddOnParameterOptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.name = value
	b.fieldSet_[0] = true
	return b
}

// Rank sets the value of the 'rank' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Rank(value int) *AddOnParameterOptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.rank = value
	b.fieldSet_[1] = true
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddOnParameterOptionBuilder) Requirements(values ...*AddOnRequirementBuilder) *AddOnParameterOptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.requirements = make([]*AddOnRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.fieldSet_[2] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Value(value string) *AddOnParameterOptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.value = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnParameterOptionBuilder) Copy(object *AddOnParameterOption) *AddOnParameterOptionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.name = object.name
	b.rank = object.rank
	if object.requirements != nil {
		b.requirements = make([]*AddOnRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddOnRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.value = object.value
	return b
}

// Build creates a 'add_on_parameter_option' object using the configuration stored in the builder.
func (b *AddOnParameterOptionBuilder) Build() (object *AddOnParameterOption, err error) {
	object = new(AddOnParameterOption)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.name = b.name
	object.rank = b.rank
	if b.requirements != nil {
		object.requirements = make([]*AddOnRequirement, len(b.requirements))
		for i, v := range b.requirements {
			object.requirements[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.value = b.value
	return
}
