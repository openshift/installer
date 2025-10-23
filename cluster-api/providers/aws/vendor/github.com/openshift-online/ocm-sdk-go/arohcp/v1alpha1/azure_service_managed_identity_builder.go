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

// AzureServiceManagedIdentityBuilder contains the data and logic needed to build 'azure_service_managed_identity' objects.
//
// Represents the information associated to an Azure User-Assigned
// Managed Identity whose purpose is to perform service level actions.
type AzureServiceManagedIdentityBuilder struct {
	bitmap_     uint32
	clientID    string
	principalID string
	resourceID  string
}

// NewAzureServiceManagedIdentity creates a new builder of 'azure_service_managed_identity' objects.
func NewAzureServiceManagedIdentity() *AzureServiceManagedIdentityBuilder {
	return &AzureServiceManagedIdentityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureServiceManagedIdentityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *AzureServiceManagedIdentityBuilder) ClientID(value string) *AzureServiceManagedIdentityBuilder {
	b.clientID = value
	b.bitmap_ |= 1
	return b
}

// PrincipalID sets the value of the 'principal_ID' attribute to the given value.
func (b *AzureServiceManagedIdentityBuilder) PrincipalID(value string) *AzureServiceManagedIdentityBuilder {
	b.principalID = value
	b.bitmap_ |= 2
	return b
}

// ResourceID sets the value of the 'resource_ID' attribute to the given value.
func (b *AzureServiceManagedIdentityBuilder) ResourceID(value string) *AzureServiceManagedIdentityBuilder {
	b.resourceID = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureServiceManagedIdentityBuilder) Copy(object *AzureServiceManagedIdentity) *AzureServiceManagedIdentityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clientID = object.clientID
	b.principalID = object.principalID
	b.resourceID = object.resourceID
	return b
}

// Build creates a 'azure_service_managed_identity' object using the configuration stored in the builder.
func (b *AzureServiceManagedIdentityBuilder) Build() (object *AzureServiceManagedIdentity, err error) {
	object = new(AzureServiceManagedIdentity)
	object.bitmap_ = b.bitmap_
	object.clientID = b.clientID
	object.principalID = b.principalID
	object.resourceID = b.resourceID
	return
}
