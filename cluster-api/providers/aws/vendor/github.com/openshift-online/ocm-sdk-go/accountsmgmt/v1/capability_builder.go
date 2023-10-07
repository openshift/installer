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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// CapabilityBuilder contains the data and logic needed to build 'capability' objects.
//
// Capability model that represents internal labels with a key that matches a set list defined in AMS (defined in pkg/api/capability_types.go).
type CapabilityBuilder struct {
	bitmap_   uint32
	name      string
	value     string
	inherited bool
}

// NewCapability creates a new builder of 'capability' objects.
func NewCapability() *CapabilityBuilder {
	return &CapabilityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CapabilityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Inherited sets the value of the 'inherited' attribute to the given value.
func (b *CapabilityBuilder) Inherited(value bool) *CapabilityBuilder {
	b.inherited = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CapabilityBuilder) Name(value string) *CapabilityBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *CapabilityBuilder) Value(value string) *CapabilityBuilder {
	b.value = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CapabilityBuilder) Copy(object *Capability) *CapabilityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.inherited = object.inherited
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'capability' object using the configuration stored in the builder.
func (b *CapabilityBuilder) Build() (object *Capability, err error) {
	object = new(Capability)
	object.bitmap_ = b.bitmap_
	object.inherited = b.inherited
	object.name = b.name
	object.value = b.value
	return
}
