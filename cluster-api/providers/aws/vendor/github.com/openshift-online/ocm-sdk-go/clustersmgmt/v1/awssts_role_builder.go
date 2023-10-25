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

// AWSSTSRoleBuilder contains the data and logic needed to build 'AWSSTS_role' objects.
//
// Representation of an sts role for a rosa cluster
type AWSSTSRoleBuilder struct {
	bitmap_            uint32
	roleARN            string
	roleType           string
	roleVersion        string
	hcpManagedPolicies bool
	isAdmin            bool
	managedPolicies    bool
}

// NewAWSSTSRole creates a new builder of 'AWSSTS_role' objects.
func NewAWSSTSRole() *AWSSTSRoleBuilder {
	return &AWSSTSRoleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSRoleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// HcpManagedPolicies sets the value of the 'hcp_managed_policies' attribute to the given value.
func (b *AWSSTSRoleBuilder) HcpManagedPolicies(value bool) *AWSSTSRoleBuilder {
	b.hcpManagedPolicies = value
	b.bitmap_ |= 1
	return b
}

// IsAdmin sets the value of the 'is_admin' attribute to the given value.
func (b *AWSSTSRoleBuilder) IsAdmin(value bool) *AWSSTSRoleBuilder {
	b.isAdmin = value
	b.bitmap_ |= 2
	return b
}

// ManagedPolicies sets the value of the 'managed_policies' attribute to the given value.
func (b *AWSSTSRoleBuilder) ManagedPolicies(value bool) *AWSSTSRoleBuilder {
	b.managedPolicies = value
	b.bitmap_ |= 4
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleARN(value string) *AWSSTSRoleBuilder {
	b.roleARN = value
	b.bitmap_ |= 8
	return b
}

// RoleType sets the value of the 'role_type' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleType(value string) *AWSSTSRoleBuilder {
	b.roleType = value
	b.bitmap_ |= 16
	return b
}

// RoleVersion sets the value of the 'role_version' attribute to the given value.
func (b *AWSSTSRoleBuilder) RoleVersion(value string) *AWSSTSRoleBuilder {
	b.roleVersion = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSRoleBuilder) Copy(object *AWSSTSRole) *AWSSTSRoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.hcpManagedPolicies = b.hcpManagedPolicies
	object.isAdmin = b.isAdmin
	object.managedPolicies = b.managedPolicies
	object.roleARN = b.roleARN
	object.roleType = b.roleType
	object.roleVersion = b.roleVersion
	return
}
