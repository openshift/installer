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

// AzureControlPlaneManagedIdentityBuilder contains the data and logic needed to build 'azure_control_plane_managed_identity' objects.
//
// Represents the information associated to an Azure User-Assigned
// Managed Identity belonging to the Control Plane of the cluster.
type AzureControlPlaneManagedIdentityBuilder struct {
	bitmap_     uint32
	clientID    string
	principalID string
	resourceID  string
}

// NewAzureControlPlaneManagedIdentity creates a new builder of 'azure_control_plane_managed_identity' objects.
func NewAzureControlPlaneManagedIdentity() *AzureControlPlaneManagedIdentityBuilder {
	return &AzureControlPlaneManagedIdentityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureControlPlaneManagedIdentityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) ClientID(value string) *AzureControlPlaneManagedIdentityBuilder {
	b.clientID = value
	b.bitmap_ |= 1
	return b
}

// PrincipalID sets the value of the 'principal_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) PrincipalID(value string) *AzureControlPlaneManagedIdentityBuilder {
	b.principalID = value
	b.bitmap_ |= 2
	return b
}

// ResourceID sets the value of the 'resource_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) ResourceID(value string) *AzureControlPlaneManagedIdentityBuilder {
	b.resourceID = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureControlPlaneManagedIdentityBuilder) Copy(object *AzureControlPlaneManagedIdentity) *AzureControlPlaneManagedIdentityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clientID = object.clientID
	b.principalID = object.principalID
	b.resourceID = object.resourceID
	return b
}

// Build creates a 'azure_control_plane_managed_identity' object using the configuration stored in the builder.
func (b *AzureControlPlaneManagedIdentityBuilder) Build() (object *AzureControlPlaneManagedIdentity, err error) {
	object = new(AzureControlPlaneManagedIdentity)
	object.bitmap_ = b.bitmap_
	object.clientID = b.clientID
	object.principalID = b.principalID
	object.resourceID = b.resourceID
	return
}
