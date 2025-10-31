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

// Representation of an addon config.
// The attributes under it are to be used by the addon once its installed in the cluster.
type AddonConfigBuilder struct {
	fieldSet_                 []bool
	addOnEnvironmentVariables []*AddonEnvironmentVariableBuilder
	addOnSecretPropagations   []*AddonSecretPropagationBuilder
}

// NewAddonConfig creates a new builder of 'addon_config' objects.
func NewAddonConfig() *AddonConfigBuilder {
	return &AddonConfigBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonConfigBuilder) Empty() bool {
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

// AddOnEnvironmentVariables sets the value of the 'add_on_environment_variables' attribute to the given values.
func (b *AddonConfigBuilder) AddOnEnvironmentVariables(values ...*AddonEnvironmentVariableBuilder) *AddonConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.addOnEnvironmentVariables = make([]*AddonEnvironmentVariableBuilder, len(values))
	copy(b.addOnEnvironmentVariables, values)
	b.fieldSet_[0] = true
	return b
}

// AddOnSecretPropagations sets the value of the 'add_on_secret_propagations' attribute to the given values.
func (b *AddonConfigBuilder) AddOnSecretPropagations(values ...*AddonSecretPropagationBuilder) *AddonConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.addOnSecretPropagations = make([]*AddonSecretPropagationBuilder, len(values))
	copy(b.addOnSecretPropagations, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonConfigBuilder) Copy(object *AddonConfig) *AddonConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.addOnEnvironmentVariables != nil {
		b.addOnEnvironmentVariables = make([]*AddonEnvironmentVariableBuilder, len(object.addOnEnvironmentVariables))
		for i, v := range object.addOnEnvironmentVariables {
			b.addOnEnvironmentVariables[i] = NewAddonEnvironmentVariable().Copy(v)
		}
	} else {
		b.addOnEnvironmentVariables = nil
	}
	if object.addOnSecretPropagations != nil {
		b.addOnSecretPropagations = make([]*AddonSecretPropagationBuilder, len(object.addOnSecretPropagations))
		for i, v := range object.addOnSecretPropagations {
			b.addOnSecretPropagations[i] = NewAddonSecretPropagation().Copy(v)
		}
	} else {
		b.addOnSecretPropagations = nil
	}
	return b
}

// Build creates a 'addon_config' object using the configuration stored in the builder.
func (b *AddonConfigBuilder) Build() (object *AddonConfig, err error) {
	object = new(AddonConfig)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.addOnEnvironmentVariables != nil {
		object.addOnEnvironmentVariables = make([]*AddonEnvironmentVariable, len(b.addOnEnvironmentVariables))
		for i, v := range b.addOnEnvironmentVariables {
			object.addOnEnvironmentVariables[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.addOnSecretPropagations != nil {
		object.addOnSecretPropagations = make([]*AddonSecretPropagation, len(b.addOnSecretPropagations))
		for i, v := range b.addOnSecretPropagations {
			object.addOnSecretPropagations[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
