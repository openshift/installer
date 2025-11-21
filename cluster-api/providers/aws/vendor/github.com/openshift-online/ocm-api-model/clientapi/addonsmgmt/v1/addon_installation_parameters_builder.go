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
type AddonInstallationParametersBuilder struct {
	fieldSet_ []bool
	items     []*AddonInstallationParameterBuilder
}

// NewAddonInstallationParameters creates a new builder of 'addon_installation_parameters' objects.
func NewAddonInstallationParameters() *AddonInstallationParametersBuilder {
	return &AddonInstallationParametersBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationParametersBuilder) Empty() bool {
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

// Items sets the value of the 'items' attribute to the given values.
func (b *AddonInstallationParametersBuilder) Items(values ...*AddonInstallationParameterBuilder) *AddonInstallationParametersBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.items = make([]*AddonInstallationParameterBuilder, len(values))
	copy(b.items, values)
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationParametersBuilder) Copy(object *AddonInstallationParameters) *AddonInstallationParametersBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.items != nil {
		b.items = make([]*AddonInstallationParameterBuilder, len(object.items))
		for i, v := range object.items {
			b.items[i] = NewAddonInstallationParameter().Copy(v)
		}
	} else {
		b.items = nil
	}
	return b
}

// Build creates a 'addon_installation_parameters' object using the configuration stored in the builder.
func (b *AddonInstallationParametersBuilder) Build() (object *AddonInstallationParameters, err error) {
	object = new(AddonInstallationParameters)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.items != nil {
		object.items = make([]*AddonInstallationParameter, len(b.items))
		for i, v := range b.items {
			object.items[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
