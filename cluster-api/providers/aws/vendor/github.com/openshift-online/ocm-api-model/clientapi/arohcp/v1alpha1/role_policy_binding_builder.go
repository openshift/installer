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

import (
	time "time"
)

type RolePolicyBindingBuilder struct {
	fieldSet_           []bool
	arn                 string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	name                string
	policies            []*RolePolicyBuilder
	status              *RolePolicyBindingStatusBuilder
	type_               string
}

// NewRolePolicyBinding creates a new builder of 'role_policy_binding' objects.
func NewRolePolicyBinding() *RolePolicyBindingBuilder {
	return &RolePolicyBindingBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RolePolicyBindingBuilder) Empty() bool {
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

// Arn sets the value of the 'arn' attribute to the given value.
func (b *RolePolicyBindingBuilder) Arn(value string) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.arn = value
	b.fieldSet_[0] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *RolePolicyBindingBuilder) CreationTimestamp(value time.Time) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.creationTimestamp = value
	b.fieldSet_[1] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *RolePolicyBindingBuilder) LastUpdateTimestamp(value time.Time) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[2] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RolePolicyBindingBuilder) Name(value string) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.name = value
	b.fieldSet_[3] = true
	return b
}

// Policies sets the value of the 'policies' attribute to the given values.
func (b *RolePolicyBindingBuilder) Policies(values ...*RolePolicyBuilder) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.policies = make([]*RolePolicyBuilder, len(values))
	copy(b.policies, values)
	b.fieldSet_[4] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *RolePolicyBindingBuilder) Status(value *RolePolicyBindingStatusBuilder) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RolePolicyBindingBuilder) Type(value string) *RolePolicyBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.type_ = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RolePolicyBindingBuilder) Copy(object *RolePolicyBinding) *RolePolicyBindingBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.arn = object.arn
	b.creationTimestamp = object.creationTimestamp
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	b.name = object.name
	if object.policies != nil {
		b.policies = make([]*RolePolicyBuilder, len(object.policies))
		for i, v := range object.policies {
			b.policies[i] = NewRolePolicy().Copy(v)
		}
	} else {
		b.policies = nil
	}
	if object.status != nil {
		b.status = NewRolePolicyBindingStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	b.type_ = object.type_
	return b
}

// Build creates a 'role_policy_binding' object using the configuration stored in the builder.
func (b *RolePolicyBindingBuilder) Build() (object *RolePolicyBinding, err error) {
	object = new(RolePolicyBinding)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.arn = b.arn
	object.creationTimestamp = b.creationTimestamp
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	object.name = b.name
	if b.policies != nil {
		object.policies = make([]*RolePolicy, len(b.policies))
		for i, v := range b.policies {
			object.policies[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	return
}
