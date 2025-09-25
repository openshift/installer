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

// AzureOperatorsAuthenticationManagedIdentitiesBuilder contains the data and logic needed to build 'azure_operators_authentication_managed_identities' objects.
//
// Represents the information related to Azure User-Assigned managed identities
// needed to perform Operators authentication based on Azure User-Assigned
// Managed Identities
type AzureOperatorsAuthenticationManagedIdentitiesBuilder struct {
	bitmap_                                uint32
	controlPlaneOperatorsManagedIdentities map[string]*AzureControlPlaneManagedIdentityBuilder
	dataPlaneOperatorsManagedIdentities    map[string]*AzureDataPlaneManagedIdentityBuilder
	managedIdentitiesDataPlaneIdentityUrl  string
	serviceManagedIdentity                 *AzureServiceManagedIdentityBuilder
}

// NewAzureOperatorsAuthenticationManagedIdentities creates a new builder of 'azure_operators_authentication_managed_identities' objects.
func NewAzureOperatorsAuthenticationManagedIdentities() *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	return &AzureOperatorsAuthenticationManagedIdentitiesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ControlPlaneOperatorsManagedIdentities sets the value of the 'control_plane_operators_managed_identities' attribute to the given value.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) ControlPlaneOperatorsManagedIdentities(value map[string]*AzureControlPlaneManagedIdentityBuilder) *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	b.controlPlaneOperatorsManagedIdentities = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// DataPlaneOperatorsManagedIdentities sets the value of the 'data_plane_operators_managed_identities' attribute to the given value.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) DataPlaneOperatorsManagedIdentities(value map[string]*AzureDataPlaneManagedIdentityBuilder) *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	b.dataPlaneOperatorsManagedIdentities = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// ManagedIdentitiesDataPlaneIdentityUrl sets the value of the 'managed_identities_data_plane_identity_url' attribute to the given value.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) ManagedIdentitiesDataPlaneIdentityUrl(value string) *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	b.managedIdentitiesDataPlaneIdentityUrl = value
	b.bitmap_ |= 4
	return b
}

// ServiceManagedIdentity sets the value of the 'service_managed_identity' attribute to the given value.
//
// Represents the information associated to an Azure User-Assigned
// Managed Identity whose purpose is to perform service level actions.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) ServiceManagedIdentity(value *AzureServiceManagedIdentityBuilder) *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	b.serviceManagedIdentity = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) Copy(object *AzureOperatorsAuthenticationManagedIdentities) *AzureOperatorsAuthenticationManagedIdentitiesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if len(object.controlPlaneOperatorsManagedIdentities) > 0 {
		b.controlPlaneOperatorsManagedIdentities = map[string]*AzureControlPlaneManagedIdentityBuilder{}
		for k, v := range object.controlPlaneOperatorsManagedIdentities {
			b.controlPlaneOperatorsManagedIdentities[k] = NewAzureControlPlaneManagedIdentity().Copy(v)
		}
	} else {
		b.controlPlaneOperatorsManagedIdentities = nil
	}
	if len(object.dataPlaneOperatorsManagedIdentities) > 0 {
		b.dataPlaneOperatorsManagedIdentities = map[string]*AzureDataPlaneManagedIdentityBuilder{}
		for k, v := range object.dataPlaneOperatorsManagedIdentities {
			b.dataPlaneOperatorsManagedIdentities[k] = NewAzureDataPlaneManagedIdentity().Copy(v)
		}
	} else {
		b.dataPlaneOperatorsManagedIdentities = nil
	}
	b.managedIdentitiesDataPlaneIdentityUrl = object.managedIdentitiesDataPlaneIdentityUrl
	if object.serviceManagedIdentity != nil {
		b.serviceManagedIdentity = NewAzureServiceManagedIdentity().Copy(object.serviceManagedIdentity)
	} else {
		b.serviceManagedIdentity = nil
	}
	return b
}

// Build creates a 'azure_operators_authentication_managed_identities' object using the configuration stored in the builder.
func (b *AzureOperatorsAuthenticationManagedIdentitiesBuilder) Build() (object *AzureOperatorsAuthenticationManagedIdentities, err error) {
	object = new(AzureOperatorsAuthenticationManagedIdentities)
	object.bitmap_ = b.bitmap_
	if b.controlPlaneOperatorsManagedIdentities != nil {
		object.controlPlaneOperatorsManagedIdentities = make(map[string]*AzureControlPlaneManagedIdentity)
		for k, v := range b.controlPlaneOperatorsManagedIdentities {
			object.controlPlaneOperatorsManagedIdentities[k], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.dataPlaneOperatorsManagedIdentities != nil {
		object.dataPlaneOperatorsManagedIdentities = make(map[string]*AzureDataPlaneManagedIdentity)
		for k, v := range b.dataPlaneOperatorsManagedIdentities {
			object.dataPlaneOperatorsManagedIdentities[k], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.managedIdentitiesDataPlaneIdentityUrl = b.managedIdentitiesDataPlaneIdentityUrl
	if b.serviceManagedIdentity != nil {
		object.serviceManagedIdentity, err = b.serviceManagedIdentity.Build()
		if err != nil {
			return
		}
	}
	return
}
