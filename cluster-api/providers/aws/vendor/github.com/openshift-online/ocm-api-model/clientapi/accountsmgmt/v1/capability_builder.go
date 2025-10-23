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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// Capability model that represents internal labels with a key that matches a set list defined in AMS (defined in pkg/api/capability_types.go).
type CapabilityBuilder struct {
	fieldSet_ []bool
	name      string
	value     string
	inherited bool
}

// NewCapability creates a new builder of 'capability' objects.
func NewCapability() *CapabilityBuilder {
	return &CapabilityBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CapabilityBuilder) Empty() bool {
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

// Inherited sets the value of the 'inherited' attribute to the given value.
func (b *CapabilityBuilder) Inherited(value bool) *CapabilityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.inherited = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CapabilityBuilder) Name(value string) *CapabilityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *CapabilityBuilder) Value(value string) *CapabilityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.value = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CapabilityBuilder) Copy(object *Capability) *CapabilityBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.inherited = object.inherited
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'capability' object using the configuration stored in the builder.
func (b *CapabilityBuilder) Build() (object *Capability, err error) {
	object = new(Capability)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.inherited = b.inherited
	object.name = b.name
	object.value = b.value
	return
}
