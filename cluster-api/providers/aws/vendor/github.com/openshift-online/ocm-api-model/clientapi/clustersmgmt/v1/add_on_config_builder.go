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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of an add-on config.
// The attributes under it are to be used by the addon once its installed in the cluster.
type AddOnConfigBuilder struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	addOnEnvironmentVariables []*AddOnEnvironmentVariableBuilder
	secretPropagations        []*AddOnSecretPropagationBuilder
}

// NewAddOnConfig creates a new builder of 'add_on_config' objects.
func NewAddOnConfig() *AddOnConfigBuilder {
	return &AddOnConfigBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnConfigBuilder) Link(value bool) *AddOnConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddOnConfigBuilder) ID(value string) *AddOnConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddOnConfigBuilder) HREF(value string) *AddOnConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnConfigBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AddOnEnvironmentVariables sets the value of the 'add_on_environment_variables' attribute to the given values.
func (b *AddOnConfigBuilder) AddOnEnvironmentVariables(values ...*AddOnEnvironmentVariableBuilder) *AddOnConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.addOnEnvironmentVariables = make([]*AddOnEnvironmentVariableBuilder, len(values))
	copy(b.addOnEnvironmentVariables, values)
	b.fieldSet_[3] = true
	return b
}

// SecretPropagations sets the value of the 'secret_propagations' attribute to the given values.
func (b *AddOnConfigBuilder) SecretPropagations(values ...*AddOnSecretPropagationBuilder) *AddOnConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.secretPropagations = make([]*AddOnSecretPropagationBuilder, len(values))
	copy(b.secretPropagations, values)
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnConfigBuilder) Copy(object *AddOnConfig) *AddOnConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.addOnEnvironmentVariables != nil {
		b.addOnEnvironmentVariables = make([]*AddOnEnvironmentVariableBuilder, len(object.addOnEnvironmentVariables))
		for i, v := range object.addOnEnvironmentVariables {
			b.addOnEnvironmentVariables[i] = NewAddOnEnvironmentVariable().Copy(v)
		}
	} else {
		b.addOnEnvironmentVariables = nil
	}
	if object.secretPropagations != nil {
		b.secretPropagations = make([]*AddOnSecretPropagationBuilder, len(object.secretPropagations))
		for i, v := range object.secretPropagations {
			b.secretPropagations[i] = NewAddOnSecretPropagation().Copy(v)
		}
	} else {
		b.secretPropagations = nil
	}
	return b
}

// Build creates a 'add_on_config' object using the configuration stored in the builder.
func (b *AddOnConfigBuilder) Build() (object *AddOnConfig, err error) {
	object = new(AddOnConfig)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.addOnEnvironmentVariables != nil {
		object.addOnEnvironmentVariables = make([]*AddOnEnvironmentVariable, len(b.addOnEnvironmentVariables))
		for i, v := range b.addOnEnvironmentVariables {
			object.addOnEnvironmentVariables[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.secretPropagations != nil {
		object.secretPropagations = make([]*AddOnSecretPropagation, len(b.secretPropagations))
		for i, v := range b.secretPropagations {
			object.secretPropagations[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
