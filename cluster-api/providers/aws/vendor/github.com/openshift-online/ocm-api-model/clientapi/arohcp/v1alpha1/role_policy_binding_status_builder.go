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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

type RolePolicyBindingStatusBuilder struct {
	fieldSet_   []bool
	description string
	value       string
}

// NewRolePolicyBindingStatus creates a new builder of 'role_policy_binding_status' objects.
func NewRolePolicyBindingStatus() *RolePolicyBindingStatusBuilder {
	return &RolePolicyBindingStatusBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RolePolicyBindingStatusBuilder) Empty() bool {
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

// Description sets the value of the 'description' attribute to the given value.
func (b *RolePolicyBindingStatusBuilder) Description(value string) *RolePolicyBindingStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.description = value
	b.fieldSet_[0] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *RolePolicyBindingStatusBuilder) Value(value string) *RolePolicyBindingStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.value = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RolePolicyBindingStatusBuilder) Copy(object *RolePolicyBindingStatus) *RolePolicyBindingStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.description = object.description
	b.value = object.value
	return b
}

// Build creates a 'role_policy_binding_status' object using the configuration stored in the builder.
func (b *RolePolicyBindingStatusBuilder) Build() (object *RolePolicyBindingStatus, err error) {
	object = new(RolePolicyBindingStatus)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.description = b.description
	object.value = b.value
	return
}
