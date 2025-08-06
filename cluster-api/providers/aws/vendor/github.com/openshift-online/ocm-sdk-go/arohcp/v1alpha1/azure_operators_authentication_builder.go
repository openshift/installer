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

// AzureOperatorsAuthenticationBuilder contains the data and logic needed to build 'azure_operators_authentication' objects.
//
// The configuration that the operators of the
// cluster have to authenticate to Azure.
type AzureOperatorsAuthenticationBuilder struct {
	bitmap_           uint32
	managedIdentities *AzureOperatorsAuthenticationManagedIdentitiesBuilder
}

// NewAzureOperatorsAuthentication creates a new builder of 'azure_operators_authentication' objects.
func NewAzureOperatorsAuthentication() *AzureOperatorsAuthenticationBuilder {
	return &AzureOperatorsAuthenticationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureOperatorsAuthenticationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ManagedIdentities sets the value of the 'managed_identities' attribute to the given value.
//
// Represents the information related to Azure User-Assigned managed identities
// needed to perform Operators authentication based on Azure User-Assigned
// Managed Identities
func (b *AzureOperatorsAuthenticationBuilder) ManagedIdentities(value *AzureOperatorsAuthenticationManagedIdentitiesBuilder) *AzureOperatorsAuthenticationBuilder {
	b.managedIdentities = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureOperatorsAuthenticationBuilder) Copy(object *AzureOperatorsAuthentication) *AzureOperatorsAuthenticationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.managedIdentities != nil {
		b.managedIdentities = NewAzureOperatorsAuthenticationManagedIdentities().Copy(object.managedIdentities)
	} else {
		b.managedIdentities = nil
	}
	return b
}

// Build creates a 'azure_operators_authentication' object using the configuration stored in the builder.
func (b *AzureOperatorsAuthenticationBuilder) Build() (object *AzureOperatorsAuthentication, err error) {
	object = new(AzureOperatorsAuthentication)
	object.bitmap_ = b.bitmap_
	if b.managedIdentities != nil {
		object.managedIdentities, err = b.managedIdentities.Build()
		if err != nil {
			return
		}
	}
	return
}
