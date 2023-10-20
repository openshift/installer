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

// AddonRequirementBuilder contains the data and logic needed to build 'addon_requirement' objects.
//
// Representation of an addon requirement.
type AddonRequirementBuilder struct {
	bitmap_  uint32
	id       string
	data     map[string]interface{}
	resource AddonRequirementResource
	status   *AddonRequirementStatusBuilder
	enabled  bool
}

// NewAddonRequirement creates a new builder of 'addon_requirement' objects.
func NewAddonRequirement() *AddonRequirementBuilder {
	return &AddonRequirementBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonRequirementBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AddonRequirementBuilder) ID(value string) *AddonRequirementBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Data sets the value of the 'data' attribute to the given value.
func (b *AddonRequirementBuilder) Data(value map[string]interface{}) *AddonRequirementBuilder {
	b.data = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonRequirementBuilder) Enabled(value bool) *AddonRequirementBuilder {
	b.enabled = value
	b.bitmap_ |= 4
	return b
}

// Resource sets the value of the 'resource' attribute to the given value.
//
// Addon requirement resource type
func (b *AddonRequirementBuilder) Resource(value AddonRequirementResource) *AddonRequirementBuilder {
	b.resource = value
	b.bitmap_ |= 8
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of an addon requirement status.
func (b *AddonRequirementBuilder) Status(value *AddonRequirementStatusBuilder) *AddonRequirementBuilder {
	b.status = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonRequirementBuilder) Copy(object *AddonRequirement) *AddonRequirementBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	if len(object.data) > 0 {
		b.data = map[string]interface{}{}
		for k, v := range object.data {
			b.data[k] = v
		}
	} else {
		b.data = nil
	}
	b.enabled = object.enabled
	b.resource = object.resource
	if object.status != nil {
		b.status = NewAddonRequirementStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	return b
}

// Build creates a 'addon_requirement' object using the configuration stored in the builder.
func (b *AddonRequirementBuilder) Build() (object *AddonRequirement, err error) {
	object = new(AddonRequirement)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	if b.data != nil {
		object.data = make(map[string]interface{})
		for k, v := range b.data {
			object.data[k] = v
		}
	}
	object.enabled = b.enabled
	object.resource = b.resource
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	return
}
