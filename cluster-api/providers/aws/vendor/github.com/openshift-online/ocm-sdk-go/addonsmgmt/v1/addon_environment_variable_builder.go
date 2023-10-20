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

// AddonEnvironmentVariableBuilder contains the data and logic needed to build 'addon_environment_variable' objects.
//
// Representation of an addon env object.
type AddonEnvironmentVariableBuilder struct {
	bitmap_ uint32
	id      string
	name    string
	value   string
	enabled bool
}

// NewAddonEnvironmentVariable creates a new builder of 'addon_environment_variable' objects.
func NewAddonEnvironmentVariable() *AddonEnvironmentVariableBuilder {
	return &AddonEnvironmentVariableBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonEnvironmentVariableBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) ID(value string) *AddonEnvironmentVariableBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Enabled(value bool) *AddonEnvironmentVariableBuilder {
	b.enabled = value
	b.bitmap_ |= 2
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Name(value string) *AddonEnvironmentVariableBuilder {
	b.name = value
	b.bitmap_ |= 4
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Value(value string) *AddonEnvironmentVariableBuilder {
	b.value = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonEnvironmentVariableBuilder) Copy(object *AddonEnvironmentVariable) *AddonEnvironmentVariableBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.enabled = object.enabled
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'addon_environment_variable' object using the configuration stored in the builder.
func (b *AddonEnvironmentVariableBuilder) Build() (object *AddonEnvironmentVariable, err error) {
	object = new(AddonEnvironmentVariable)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.enabled = b.enabled
	object.name = b.name
	object.value = b.value
	return
}
