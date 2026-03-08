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

// Representation of an add-on requirement.
type AddOnRequirementBuilder struct {
	fieldSet_ []bool
	id        string
	data      map[string]interface{}
	resource  string
	status    *AddOnRequirementStatusBuilder
	enabled   bool
}

// NewAddOnRequirement creates a new builder of 'add_on_requirement' objects.
func NewAddOnRequirement() *AddOnRequirementBuilder {
	return &AddOnRequirementBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnRequirementBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *AddOnRequirementBuilder) ID(value string) *AddOnRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Data sets the value of the 'data' attribute to the given value.
func (b *AddOnRequirementBuilder) Data(value map[string]interface{}) *AddOnRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.data = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnRequirementBuilder) Enabled(value bool) *AddOnRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.enabled = value
	b.fieldSet_[2] = true
	return b
}

// Resource sets the value of the 'resource' attribute to the given value.
func (b *AddOnRequirementBuilder) Resource(value string) *AddOnRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.resource = value
	b.fieldSet_[3] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of an add-on requirement status.
func (b *AddOnRequirementBuilder) Status(value *AddOnRequirementStatusBuilder) *AddOnRequirementBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnRequirementBuilder) Copy(object *AddOnRequirement) *AddOnRequirementBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
		b.status = NewAddOnRequirementStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	return b
}

// Build creates a 'add_on_requirement' object using the configuration stored in the builder.
func (b *AddOnRequirementBuilder) Build() (object *AddOnRequirement, err error) {
	object = new(AddOnRequirement)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
