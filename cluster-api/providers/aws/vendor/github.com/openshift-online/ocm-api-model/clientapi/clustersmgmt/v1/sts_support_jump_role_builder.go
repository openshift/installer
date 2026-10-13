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

// Isolated STS support role created per organization.
type StsSupportJumpRoleBuilder struct {
	fieldSet_ []bool
	roleArn   string
}

// NewStsSupportJumpRole creates a new builder of 'sts_support_jump_role' objects.
func NewStsSupportJumpRole() *StsSupportJumpRoleBuilder {
	return &StsSupportJumpRoleBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StsSupportJumpRoleBuilder) Empty() bool {
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

// RoleArn sets the value of the 'role_arn' attribute to the given value.
func (b *StsSupportJumpRoleBuilder) RoleArn(value string) *StsSupportJumpRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.roleArn = value
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StsSupportJumpRoleBuilder) Copy(object *StsSupportJumpRole) *StsSupportJumpRoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.roleArn = object.roleArn
	return b
}

// Build creates a 'sts_support_jump_role' object using the configuration stored in the builder.
func (b *StsSupportJumpRoleBuilder) Build() (object *StsSupportJumpRole, err error) {
	object = new(StsSupportJumpRole)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.roleArn = b.roleArn
	return
}
