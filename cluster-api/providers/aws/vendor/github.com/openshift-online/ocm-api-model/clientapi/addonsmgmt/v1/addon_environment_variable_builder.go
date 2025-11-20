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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// Representation of an addon env object.
type AddonEnvironmentVariableBuilder struct {
	fieldSet_ []bool
	id        string
	name      string
	value     string
	enabled   bool
}

// NewAddonEnvironmentVariable creates a new builder of 'addon_environment_variable' objects.
func NewAddonEnvironmentVariable() *AddonEnvironmentVariableBuilder {
	return &AddonEnvironmentVariableBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonEnvironmentVariableBuilder) Empty() bool {
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
func (b *AddonEnvironmentVariableBuilder) ID(value string) *AddonEnvironmentVariableBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Enabled(value bool) *AddonEnvironmentVariableBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.enabled = value
	b.fieldSet_[1] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Name(value string) *AddonEnvironmentVariableBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.name = value
	b.fieldSet_[2] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddonEnvironmentVariableBuilder) Value(value string) *AddonEnvironmentVariableBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.value = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonEnvironmentVariableBuilder) Copy(object *AddonEnvironmentVariable) *AddonEnvironmentVariableBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.enabled = object.enabled
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'addon_environment_variable' object using the configuration stored in the builder.
func (b *AddonEnvironmentVariableBuilder) Build() (object *AddonEnvironmentVariable, err error) {
	object = new(AddonEnvironmentVariable)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.enabled = b.enabled
	object.name = b.name
	object.value = b.value
	return
}
