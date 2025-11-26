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

// Contains the necessary attributes to allow each operator to access the necessary AWS resources
type OperatorIAMRoleBuilder struct {
	fieldSet_      []bool
	id             string
	name           string
	namespace      string
	roleARN        string
	serviceAccount string
}

// NewOperatorIAMRole creates a new builder of 'operator_IAM_role' objects.
func NewOperatorIAMRole() *OperatorIAMRoleBuilder {
	return &OperatorIAMRoleBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OperatorIAMRoleBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *OperatorIAMRoleBuilder) ID(value string) *OperatorIAMRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *OperatorIAMRoleBuilder) Name(value string) *OperatorIAMRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *OperatorIAMRoleBuilder) Namespace(value string) *OperatorIAMRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.namespace = value
	b.fieldSet_[2] = true
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *OperatorIAMRoleBuilder) RoleARN(value string) *OperatorIAMRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.roleARN = value
	b.fieldSet_[3] = true
	return b
}

// ServiceAccount sets the value of the 'service_account' attribute to the given value.
func (b *OperatorIAMRoleBuilder) ServiceAccount(value string) *OperatorIAMRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.serviceAccount = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OperatorIAMRoleBuilder) Copy(object *OperatorIAMRole) *OperatorIAMRoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.name = object.name
	b.namespace = object.namespace
	b.roleARN = object.roleARN
	b.serviceAccount = object.serviceAccount
	return b
}

// Build creates a 'operator_IAM_role' object using the configuration stored in the builder.
func (b *OperatorIAMRoleBuilder) Build() (object *OperatorIAMRole, err error) {
	object = new(OperatorIAMRole)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.name = b.name
	object.namespace = b.namespace
	object.roleARN = b.roleARN
	object.serviceAccount = b.serviceAccount
	return
}
