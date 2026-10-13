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

// Representation of an sts role for a rosa cluster
type AWSSTSRoleBuilder struct {
	fieldSet_          []bool
	roleARN            string
	roleType           string
	roleVersion        string
	hcpManagedPolicies bool
	isAdmin            bool
	managedPolicies    bool
}

// NewAWSSTSRole creates a new builder of 'AWSSTS_role' objects.
func NewAWSSTSRole() *AWSSTSRoleBuilder {
	return &AWSSTSRoleBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSRoleBuilder) Empty() bool {
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

// HcpManagedPolicies sets the value of the 'hcp_managed_policies' attribute to the given value.
func (b *AWSSTSRoleBuilder) HcpManagedPolicies(value bool) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.hcpManagedPolicies = value
	b.fieldSet_[0] = true
	return b
}

// IsAdmin sets the value of the 'is_admin' attribute to the given value.
func (b *AWSSTSRoleBuilder) IsAdmin(value bool) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.isAdmin = value
	b.fieldSet_[1] = true
	return b
}

// ManagedPolicies sets the value of the 'managed_policies' attribute to the given value.
func (b *AWSSTSRoleBuilder) ManagedPolicies(value bool) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.managedPolicies = value
	b.fieldSet_[2] = true
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleARN(value string) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.roleARN = value
	b.fieldSet_[3] = true
	return b
}

// RoleType sets the value of the 'role_type' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleType(value string) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.roleType = value
	b.fieldSet_[4] = true
	return b
}

// RoleVersion sets the value of the 'role_version' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleVersion(value string) *AWSSTSRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.roleVersion = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSRoleBuilder) Copy(object *AWSSTSRole) *AWSSTSRoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.hcpManagedPolicies = object.hcpManagedPolicies
	b.isAdmin = object.isAdmin
	b.managedPolicies = object.managedPolicies
	b.roleARN = object.roleARN
	b.roleType = object.roleType
	b.roleVersion = object.roleVersion
	return b
}

// Build creates a 'AWSSTS_role' object using the configuration stored in the builder.
func (b *AWSSTSRoleBuilder) Build() (object *AWSSTSRole, err error) {
	object = new(AWSSTSRole)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.hcpManagedPolicies = b.hcpManagedPolicies
	object.isAdmin = b.isAdmin
	object.managedPolicies = b.managedPolicies
	object.roleARN = b.roleARN
	object.roleType = b.roleType
	object.roleVersion = b.roleVersion
	return
}
