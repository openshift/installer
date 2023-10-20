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

// AddonParameterOptionBuilder contains the data and logic needed to build 'addon_parameter_option' objects.
//
// Representation of an addon parameter option.
type AddonParameterOptionBuilder struct {
	bitmap_      uint32
	name         string
	rank         int
	requirements []*AddonRequirementBuilder
	value        string
}

// NewAddonParameterOption creates a new builder of 'addon_parameter_option' objects.
func NewAddonParameterOption() *AddonParameterOptionBuilder {
	return &AddonParameterOptionBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonParameterOptionBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonParameterOptionBuilder) Name(value string) *AddonParameterOptionBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Rank sets the value of the 'rank' attribute to the given value.
func (b *AddonParameterOptionBuilder) Rank(value int) *AddonParameterOptionBuilder {
	b.rank = value
	b.bitmap_ |= 2
	return b
}

// Requirements sets the value of the 'requirements' attribute to the given values.
func (b *AddonParameterOptionBuilder) Requirements(values ...*AddonRequirementBuilder) *AddonParameterOptionBuilder {
	b.requirements = make([]*AddonRequirementBuilder, len(values))
	copy(b.requirements, values)
	b.bitmap_ |= 4
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddonParameterOptionBuilder) Value(value string) *AddonParameterOptionBuilder {
	b.value = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonParameterOptionBuilder) Copy(object *AddonParameterOption) *AddonParameterOptionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	b.rank = object.rank
	if object.requirements != nil {
		b.requirements = make([]*AddonRequirementBuilder, len(object.requirements))
		for i, v := range object.requirements {
			b.requirements[i] = NewAddonRequirement().Copy(v)
		}
	} else {
		b.requirements = nil
	}
	b.value = object.value
	return b
}

// Build creates a 'addon_parameter_option' object using the configuration stored in the builder.
func (b *AddonParameterOptionBuilder) Build() (object *AddonParameterOption, err error) {
	object = new(AddonParameterOption)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.rank = b.rank
	if b.requirements != nil {
		object.requirements = make([]*AddonRequirement, len(b.requirements))
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
