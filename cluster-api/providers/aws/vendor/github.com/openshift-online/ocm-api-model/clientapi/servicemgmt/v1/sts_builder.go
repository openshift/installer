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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

// Contains the necessary attributes to support role-based authentication on AWS.
type STSBuilder struct {
	fieldSet_          []bool
	oidcEndpointURL    string
	instanceIAMRoles   *InstanceIAMRolesBuilder
	operatorIAMRoles   []*OperatorIAMRoleBuilder
	operatorRolePrefix string
	roleARN            string
	supportRoleARN     string
}

// NewSTS creates a new builder of 'STS' objects.
func NewSTS() *STSBuilder {
	return &STSBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *STSBuilder) Empty() bool {
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

// OIDCEndpointURL sets the value of the 'OIDC_endpoint_URL' attribute to the given value.
func (b *STSBuilder) OIDCEndpointURL(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.oidcEndpointURL = value
	b.fieldSet_[0] = true
	return b
}

// InstanceIAMRoles sets the value of the 'instance_IAM_roles' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *STSBuilder) InstanceIAMRoles(value *InstanceIAMRolesBuilder) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.instanceIAMRoles = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// OperatorIAMRoles sets the value of the 'operator_IAM_roles' attribute to the given values.
func (b *STSBuilder) OperatorIAMRoles(values ...*OperatorIAMRoleBuilder) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.operatorIAMRoles = make([]*OperatorIAMRoleBuilder, len(values))
	copy(b.operatorIAMRoles, values)
	b.fieldSet_[2] = true
	return b
}

// OperatorRolePrefix sets the value of the 'operator_role_prefix' attribute to the given value.
func (b *STSBuilder) OperatorRolePrefix(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.operatorRolePrefix = value
	b.fieldSet_[3] = true
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *STSBuilder) RoleARN(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.roleARN = value
	b.fieldSet_[4] = true
	return b
}

// SupportRoleARN sets the value of the 'support_role_ARN' attribute to the given value.
func (b *STSBuilder) SupportRoleARN(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.supportRoleARN = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSBuilder) Copy(object *STS) *STSBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
