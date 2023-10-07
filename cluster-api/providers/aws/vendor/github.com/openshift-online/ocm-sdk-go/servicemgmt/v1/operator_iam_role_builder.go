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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// OperatorIAMRoleBuilder contains the data and logic needed to build 'operator_IAM_role' objects.
//
// Contains the necessary attributes to allow each operator to access the necessary AWS resources
type OperatorIAMRoleBuilder struct {
	bitmap_   uint32
	name      string
	namespace string
	roleARN   string
}

// NewOperatorIAMRole creates a new builder of 'operator_IAM_role' objects.
func NewOperatorIAMRole() *OperatorIAMRoleBuilder {
	return &OperatorIAMRoleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OperatorIAMRoleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *OperatorIAMRoleBuilder) Name(value string) *OperatorIAMRoleBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *OperatorIAMRoleBuilder) Namespace(value string) *OperatorIAMRoleBuilder {
	b.namespace = value
	b.bitmap_ |= 2
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *OperatorIAMRoleBuilder) RoleARN(value string) *OperatorIAMRoleBuilder {
	b.roleARN = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OperatorIAMRoleBuilder) Copy(object *OperatorIAMRole) *OperatorIAMRoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	b.namespace = object.namespace
	b.roleARN = object.roleARN
	return b
}

// Build creates a 'operator_IAM_role' object using the configuration stored in the builder.
func (b *OperatorIAMRoleBuilder) Build() (object *OperatorIAMRole, err error) {
	object = new(OperatorIAMRole)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.namespace = b.namespace
	object.roleARN = b.roleARN
	return
}
