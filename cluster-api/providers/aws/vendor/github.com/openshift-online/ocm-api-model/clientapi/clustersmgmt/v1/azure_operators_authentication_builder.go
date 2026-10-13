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

// The configuration that the operators of the
// cluster have to authenticate to Azure.
type AzureOperatorsAuthenticationBuilder struct {
	fieldSet_         []bool
	managedIdentities *AzureOperatorsAuthenticationManagedIdentitiesBuilder
}

// NewAzureOperatorsAuthentication creates a new builder of 'azure_operators_authentication' objects.
func NewAzureOperatorsAuthentication() *AzureOperatorsAuthenticationBuilder {
	return &AzureOperatorsAuthenticationBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureOperatorsAuthenticationBuilder) Empty() bool {
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

// ManagedIdentities sets the value of the 'managed_identities' attribute to the given value.
//
// Represents the information related to Azure User-Assigned managed identities
// needed to perform Operators authentication based on Azure User-Assigned
// Managed Identities
func (b *AzureOperatorsAuthenticationBuilder) ManagedIdentities(value *AzureOperatorsAuthenticationManagedIdentitiesBuilder) *AzureOperatorsAuthenticationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.managedIdentities = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureOperatorsAuthenticationBuilder) Copy(object *AzureOperatorsAuthentication) *AzureOperatorsAuthenticationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.managedIdentities != nil {
		object.managedIdentities, err = b.managedIdentities.Build()
		if err != nil {
			return
		}
	}
	return
}
