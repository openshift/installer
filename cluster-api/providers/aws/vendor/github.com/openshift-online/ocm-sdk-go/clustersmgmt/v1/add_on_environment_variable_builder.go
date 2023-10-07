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

// AddOnEnvironmentVariableBuilder contains the data and logic needed to build 'add_on_environment_variable' objects.
//
// Representation of an add-on env object.
type AddOnEnvironmentVariableBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	name    string
	value   string
}

// NewAddOnEnvironmentVariable creates a new builder of 'add_on_environment_variable' objects.
func NewAddOnEnvironmentVariable() *AddOnEnvironmentVariableBuilder {
	return &AddOnEnvironmentVariableBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnEnvironmentVariableBuilder) Link(value bool) *AddOnEnvironmentVariableBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnEnvironmentVariableBuilder) ID(value string) *AddOnEnvironmentVariableBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnEnvironmentVariableBuilder) HREF(value string) *AddOnEnvironmentVariableBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnEnvironmentVariableBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AddOnEnvironmentVariableBuilder) Name(value string) *AddOnEnvironmentVariableBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddOnEnvironmentVariableBuilder) Value(value string) *AddOnEnvironmentVariableBuilder {
	b.value = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnEnvironmentVariableBuilder) Copy(object *AddOnEnvironmentVariable) *AddOnEnvironmentVariableBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.name = object.name
	b.value = object.value
	return b
}

// Build creates a 'add_on_environment_variable' object using the configuration stored in the builder.
func (b *AddOnEnvironmentVariableBuilder) Build() (object *AddOnEnvironmentVariable, err error) {
	object = new(AddOnEnvironmentVariable)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.value = b.value
	return
}
