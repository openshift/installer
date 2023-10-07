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

// AddOnInstallationParameterBuilder contains the data and logic needed to build 'add_on_installation_parameter' objects.
//
// Representation of an add-on installation parameter.
type AddOnInstallationParameterBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	value   string
}

// NewAddOnInstallationParameter creates a new builder of 'add_on_installation_parameter' objects.
func NewAddOnInstallationParameter() *AddOnInstallationParameterBuilder {
	return &AddOnInstallationParameterBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnInstallationParameterBuilder) Link(value bool) *AddOnInstallationParameterBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnInstallationParameterBuilder) ID(value string) *AddOnInstallationParameterBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnInstallationParameterBuilder) HREF(value string) *AddOnInstallationParameterBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnInstallationParameterBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddOnInstallationParameterBuilder) Value(value string) *AddOnInstallationParameterBuilder {
	b.value = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnInstallationParameterBuilder) Copy(object *AddOnInstallationParameter) *AddOnInstallationParameterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.value = object.value
	return b
}

// Build creates a 'add_on_installation_parameter' object using the configuration stored in the builder.
func (b *AddOnInstallationParameterBuilder) Build() (object *AddOnInstallationParameter, err error) {
	object = new(AddOnInstallationParameter)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.value = b.value
	return
}
