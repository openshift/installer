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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// STSBuilder contains the data and logic needed to build 'STS' objects.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type STSBuilder struct {
	bitmap_            uint32
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

// AutoMode sets the value of the 'auto_mode' attribute to the given value.
func (b *STSBuilder) AutoMode(value bool) *STSBuilder {
	b.autoMode = value
	b.bitmap_ |= 2
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *STSBuilder) Enabled(value bool) *STSBuilder {
	b.enabled = value
	b.bitmap_ |= 4
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *STSBuilder) ExternalID(value string) *STSBuilder {
	b.externalID = value
	b.bitmap_ |= 8
	return b
}

// InstanceIAMRoles sets the value of the 'instance_IAM_roles' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *STSBuilder) InstanceIAMRoles(value *InstanceIAMRolesBuilder) *STSBuilder {
	b.instanceIAMRoles = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// ManagedPolicies sets the value of the 'managed_policies' attribute to the given value.
func (b *STSBuilder) ManagedPolicies(value bool) *STSBuilder {
	b.managedPolicies = value
	b.bitmap_ |= 32
	return b
}

// OidcConfig sets the value of the 'oidc_config' attribute to the given value.
//
// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
func (b *STSBuilder) OidcConfig(value *OidcConfigBuilder) *STSBuilder {
	b.oidcConfig = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// OperatorIAMRoles sets the value of the 'operator_IAM_roles' attribute to the given values.
func (b *STSBuilder) OperatorIAMRoles(values ...*OperatorIAMRoleBuilder) *STSBuilder {
	b.operatorIAMRoles = make([]*OperatorIAMRoleBuilder, len(values))
	copy(b.operatorIAMRoles, values)
	b.bitmap_ |= 128
	return b
}

// OperatorRolePrefix sets the value of the 'operator_role_prefix' attribute to the given value.
func (b *STSBuilder) OperatorRolePrefix(value string) *STSBuilder {
	b.operatorRolePrefix = value
	b.bitmap_ |= 256
	return b
}

// PermissionBoundary sets the value of the 'permission_boundary' attribute to the given value.
func (b *STSBuilder) PermissionBoundary(value string) *STSBuilder {
	b.permissionBoundary = value
	b.bitmap_ |= 512
	return b
}

// RoleARN sets the value of the 'role_ARN' attribute to the given value.
func (b *STSBuilder) RoleARN(value string) *STSBuilder {
	b.roleARN = value
	b.bitmap_ |= 1024
	return b
}

// SupportRoleARN sets the value of the 'support_role_ARN' attribute to the given value.
func (b *STSBuilder) SupportRoleARN(value string) *STSBuilder {
	b.supportRoleARN = value
	b.bitmap_ |= 2048
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *STSBuilder) Copy(object *STS) *STSBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
