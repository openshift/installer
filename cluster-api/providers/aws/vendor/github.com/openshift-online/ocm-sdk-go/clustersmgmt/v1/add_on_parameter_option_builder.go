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

// AddOnParameterOptionBuilder contains the data and logic needed to build 'add_on_parameter_option' objects.
//
// Representation of an add-on parameter option.
type AddOnParameterOptionBuilder struct {
	bitmap_      uint32
	name         string
	rank         int
	requirements []*AddOnRequirementBuilder
	value        string
}

// NewAddOnParameterOption creates a new builder of 'add_on_parameter_option' objects.
func NewAddOnParameterOption() *AddOnParameterOptionBuilder {
	return &AddOnParameterOptionBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnParameterOptionBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Name(value string) *AddOnParameterOptionBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Rank sets the value of the 'rank' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Rank(value int) *AddOnParameterOptionBuilder {
	b.rank = value
	b.bitmap_ |= 2
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddOnParameterOptionBuilder) Requirements(values ...*AddOnRequirementBuilder) *AddOnParameterOptionBuilder {
	b.requirements = make([]*AddOnRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.bitmap_ |= 4
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddOnParameterOptionBuilder) Value(value string) *AddOnParameterOptionBuilder {
	b.value = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnParameterOptionBuilder) Copy(object *AddOnParameterOption) *AddOnParameterOptionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
