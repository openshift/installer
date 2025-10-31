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

// Contains the necessary attributes to support role-based authentication on AWS.
type STSBuilder struct {
	fieldSet_          []bool
	oidcEndpointURL    string
	externalID         string
	instanceIAMRoles   *InstanceIAMRolesBuilder
	oidcConfig         *OidcConfigBuilder
	operatorIAMRoles   []*OperatorIAMRoleBuilder
	operatorRolePrefix string
	permissionBoundary string
	roleARN            string
	supportRoleARN     string
	autoMode           bool
	enabled            bool
	managedPolicies    bool
}

// NewSTS creates a new builder of 'STS' objects.
func NewSTS() *STSBuilder {
	return &STSBuilder{
		fieldSet_: make([]bool, 12),
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
		b.fieldSet_ = make([]bool, 12)
	}
	b.oidcEndpointURL = value
	b.fieldSet_[0] = true
	return b
}

// AutoMode sets the value of the 'auto_mode' attribute to the given value.
func (b *STSBuilder) AutoMode(value bool) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.autoMode = value
	b.fieldSet_[1] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *STSBuilder) Enabled(value bool) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.enabled = value
	b.fieldSet_[2] = true
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *STSBuilder) ExternalID(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.externalID = value
	b.fieldSet_[3] = true
	return b
}

// InstanceIAMRoles sets the value of the 'instance_IAM_roles' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *STSBuilder) InstanceIAMRoles(value *InstanceIAMRolesBuilder) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.instanceIAMRoles = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// ManagedPolicies sets the value of the 'managed_policies' attribute to the given value.
func (b *STSBuilder) ManagedPolicies(value bool) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.managedPolicies = value
	b.fieldSet_[5] = true
	return b
}

// OidcConfig sets the value of the 'oidc_config' attribute to the given value.
//
// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
func (b *STSBuilder) OidcConfig(value *OidcConfigBuilder) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.oidcConfig = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// OperatorIAMRoles sets the value of the 'operator_IAM_roles' attribute to the given values.
func (b *STSBuilder) OperatorIAMRoles(values ...*OperatorIAMRoleBuilder) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.operatorIAMRoles = make([]*OperatorIAMRoleBuilder, len(values))
	copy(b.operatorIAMRoles, values)
	b.fieldSet_[7] = true
	return b
}

// OperatorRolePrefix sets the value of the 'operator_role_prefix' attribute to the given value.
func (b *STSBuilder) OperatorRolePrefix(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.operatorRolePrefix = value
	b.fieldSet_[8] = true
	return b
}

// PermissionBoundary sets the value of the 'permission_boundary' attribute to the given value.
func (b *STSBuilder) PermissionBoundary(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.permissionBoundary = value
	b.fieldSet_[9] = true
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *STSBuilder) RoleARN(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.roleARN = value
	b.fieldSet_[10] = true
	return b
}

// SupportRoleARN sets the value of the 'support_role_ARN' attribute to the given value.
func (b *STSBuilder) SupportRoleARN(value string) *STSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.supportRoleARN = value
	b.fieldSet_[11] = true
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
	b.autoMode = object.autoMode
	b.enabled = object.enabled
	b.externalID = object.externalID
	if object.instanceIAMRoles != nil {
		b.instanceIAMRoles = NewInstanceIAMRoles().Copy(object.instanceIAMRoles)
	} else {
		b.instanceIAMRoles = nil
	}
	b.managedPolicies = object.managedPolicies
	if object.oidcConfig != nil {
		b.oidcConfig = NewOidcConfig().Copy(object.oidcConfig)
	} else {
		b.oidcConfig = nil
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
	b.permissionBoundary = object.permissionBoundary
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
	object.autoMode = b.autoMode
	object.enabled = b.enabled
	object.externalID = b.externalID
	if b.instanceIAMRoles != nil {
		object.instanceIAMRoles, err = b.instanceIAMRoles.Build()
		if err != nil {
			return
		}
	}
	object.managedPolicies = b.managedPolicies
	if b.oidcConfig != nil {
		object.oidcConfig, err = b.oidcConfig.Build()
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
	object.permissionBoundary = b.permissionBoundary
	object.roleARN = b.roleARN
	object.supportRoleARN = b.supportRoleARN
	return
}
