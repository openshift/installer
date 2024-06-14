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

import (
	time "time"
)

// RolePolicyBindingBuilder contains the data and logic needed to build 'role_policy_binding' objects.
type RolePolicyBindingBuilder struct {
	bitmap_             uint32
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
	return &RolePolicyBindingBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RolePolicyBindingBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Arn sets the value of the 'arn' attribute to the given value.
func (b *RolePolicyBindingBuilder) Arn(value string) *RolePolicyBindingBuilder {
	b.arn = value
	b.bitmap_ |= 1
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *RolePolicyBindingBuilder) CreationTimestamp(value time.Time) *RolePolicyBindingBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 2
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *RolePolicyBindingBuilder) LastUpdateTimestamp(value time.Time) *RolePolicyBindingBuilder {
	b.lastUpdateTimestamp = value
	b.bitmap_ |= 4
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RolePolicyBindingBuilder) Name(value string) *RolePolicyBindingBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// Policies sets the value of the 'policies' attribute to the given values.
func (b *RolePolicyBindingBuilder) Policies(values ...*RolePolicyBuilder) *RolePolicyBindingBuilder {
	b.policies = make([]*RolePolicyBuilder, len(values))
	copy(b.policies, values)
	b.bitmap_ |= 16
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *RolePolicyBindingBuilder) Status(value *RolePolicyBindingStatusBuilder) *RolePolicyBindingBuilder {
	b.status = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RolePolicyBindingBuilder) Type(value string) *RolePolicyBindingBuilder {
	b.type_ = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RolePolicyBindingBuilder) Copy(object *RolePolicyBinding) *RolePolicyBindingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
