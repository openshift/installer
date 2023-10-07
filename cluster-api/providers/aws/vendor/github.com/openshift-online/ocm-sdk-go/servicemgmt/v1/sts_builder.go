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

// STSBuilder contains the data and logic needed to build 'STS' objects.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type STSBuilder struct {
	bitmap_            uint32
	oidcEndpointURL    string
	instanceIAMRoles   *InstanceIAMRolesBuilder
	operatorIAMRoles   []*OperatorIAMRoleBuilder
	operatorRolePrefix string
	roleARN            string
	supportRoleARN     string
}

// NewSTS creates a new builder of 'STS' objects.
func NewSTS() *STSBuilder {
	return &STSBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *STSBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// OIDCEndpointURL sets the value of the 'OIDC_endpoint_URL' attribute to the given value.
func (b *STSBuilder) OIDCEndpointURL(value string) *STSBuilder {
	b.oidcEndpointURL = value
	b.bitmap_ |= 1
	return b
}

// InstanceIAMRoles sets the value of the 'instance_IAM_roles' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *STSBuilder) InstanceIAMRoles(value *InstanceIAMRolesBuilder) *STSBuilder {
	b.instanceIAMRoles = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// OperatorIAMRoles sets the value of the 'operator_IAM_roles' attribute to the given values.
func (b *STSBuilder) OperatorIAMRoles(values ...*OperatorIAMRoleBuilder) *STSBuilder {
	b.operatorIAMRoles = make([]*OperatorIAMRoleBuilder, len(values))
	copy(b.operatorIAMRoles, values)
	b.bitmap_ |= 4
	return b
}

// OperatorRolePrefix sets the value of the 'operator_role_prefix' attribute to the given value.
func (b *STSBuilder) OperatorRolePrefix(value string) *STSBuilder {
	b.operatorRolePrefix = value
	b.bitmap_ |= 8
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *STSBuilder) RoleARN(value string) *STSBuilder {
	b.roleARN = value
	b.bitmap_ |= 16
	return b
}

// SupportRoleARN sets the value of the 'support_role_ARN' attribute to the given value.
func (b *STSBuilder) SupportRoleARN(value string) *STSBuilder {
	b.supportRoleARN = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSBuilder) Copy(object *STS) *STSBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.oidcEndpointURL = object.oidcEndpointURL
	if object.instanceIAMRoles != nil {
		b.instanceIAMRoles = NewInstanceIAMRoles().Copy(object.instanceIAMRoles)
	} else {
		b.instanceIAMRoles = nil
	}
	if object.operatorIAMRoles != nil {
		b.operatorIAMRoles = make([]*OperatorIAMRoleBuilder, len(object.operatorIAMRoles))
		for i, v := range object.operatorIAMRoles {
			b.operatorIAMRoles[i] = NewOperatorIAMRole().Copy(v)
		}
	} else {
		b.operatorIAMRoles = nil
	}
	b.operatorRolePrefix = object.operatorRolePrefix
	b.roleARN = object.roleARN
	b.supportRoleARN = object.supportRoleARN
	return b
}

// Build creates a 'STS' object using the configuration stored in the builder.
func (b *STSBuilder) Build() (object *STS, err error) {
	object = new(STS)
	object.bitmap_ = b.bitmap_
	object.oidcEndpointURL = b.oidcEndpointURL
	if b.instanceIAMRoles != nil {
		object.instanceIAMRoles, err = b.instanceIAMRoles.Build()
		if err != nil {
			return
		}
	}
	if b.operatorIAMRoles != nil {
		object.operatorIAMRoles = make([]*OperatorIAMRole, len(b.operatorIAMRoles))
		for i, v := range b.operatorIAMRoles {
			object.operatorIAMRoles[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.operatorRolePrefix = b.operatorRolePrefix
	object.roleARN = b.roleARN
	object.supportRoleARN = b.supportRoleARN
	return
}
