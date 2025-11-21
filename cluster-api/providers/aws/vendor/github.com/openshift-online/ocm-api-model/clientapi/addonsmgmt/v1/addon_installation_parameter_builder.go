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

// representation of addon installation parameter
type AddonInstallationParameterBuilder struct {
	fieldSet_ []bool
	href      string
	id        string
	kind      string
	value     string
}

// NewAddonInstallationParameter creates a new builder of 'addon_installation_parameter' objects.
func NewAddonInstallationParameter() *AddonInstallationParameterBuilder {
	return &AddonInstallationParameterBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationParameterBuilder) Empty() bool {
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

// Href sets the value of the 'href' attribute to the given value.
func (b *AddonInstallationParameterBuilder) Href(value string) *AddonInstallationParameterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.href = value
	b.fieldSet_[0] = true
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *AddonInstallationParameterBuilder) Id(value string) *AddonInstallationParameterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *AddonInstallationParameterBuilder) Kind(value string) *AddonInstallationParameterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.kind = value
	b.fieldSet_[2] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *AddonInstallationParameterBuilder) Value(value string) *AddonInstallationParameterBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.value = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationParameterBuilder) Copy(object *AddonInstallationParameter) *AddonInstallationParameterBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.href = object.href
	b.id = object.id
	b.kind = object.kind
	b.value = object.value
	return b
}

// Build creates a 'addon_installation_parameter' object using the configuration stored in the builder.
func (b *AddonInstallationParameterBuilder) Build() (object *AddonInstallationParameter, err error) {
	object = new(AddonInstallationParameter)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.href = b.href
	object.id = b.id
	object.kind = b.kind
	object.value = b.value
	return
}
