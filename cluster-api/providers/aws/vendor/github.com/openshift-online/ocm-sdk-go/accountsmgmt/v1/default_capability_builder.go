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

// DefaultCapabilityBuilder contains the data and logic needed to build 'default_capability' objects.
type DefaultCapabilityBuilder struct {
	bitmap_ uint32
	name    string
	value   string
}

// NewDefaultCapability creates a new builder of 'default_capability' objects.
func NewDefaultCapability() *DefaultCapabilityBuilder {
	return &DefaultCapabilityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DefaultCapabilityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *DefaultCapabilityBuilder) Name(value string) *DefaultCapabilityBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *DefaultCapabilityBuilder) Value(value string) *DefaultCapabilityBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DefaultCapabilityBuilder) Copy(object *DefaultCapability) *DefaultCapabilityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'default_capability' object using the configuration stored in the builder.
func (b *DefaultCapabilityBuilder) Build() (object *DefaultCapability, err error) {
	object = new(DefaultCapability)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.value = b.value
	return
}
