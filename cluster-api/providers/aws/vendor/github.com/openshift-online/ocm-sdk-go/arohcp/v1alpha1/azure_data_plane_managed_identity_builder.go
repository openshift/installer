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

// AzureDataPlaneManagedIdentityBuilder contains the data and logic needed to build 'azure_data_plane_managed_identity' objects.
//
// Represents the information associated to an Azure User-Assigned
// Managed Identity belonging to the Data Plane of the cluster.
type AzureDataPlaneManagedIdentityBuilder struct {
	bitmap_    uint32
	resourceID string
}

// NewAzureDataPlaneManagedIdentity creates a new builder of 'azure_data_plane_managed_identity' objects.
func NewAzureDataPlaneManagedIdentity() *AzureDataPlaneManagedIdentityBuilder {
	return &AzureDataPlaneManagedIdentityBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureDataPlaneManagedIdentityBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ResourceID sets the value of the 'resource_ID' attribute to the given value.
func (b *AzureDataPlaneManagedIdentityBuilder) ResourceID(value string) *AzureDataPlaneManagedIdentityBuilder {
	b.resourceID = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureDataPlaneManagedIdentityBuilder) Copy(object *AzureDataPlaneManagedIdentity) *AzureDataPlaneManagedIdentityBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.resourceID = object.resourceID
	return b
}

// Build creates a 'azure_data_plane_managed_identity' object using the configuration stored in the builder.
func (b *AzureDataPlaneManagedIdentityBuilder) Build() (object *AzureDataPlaneManagedIdentity, err error) {
	object = new(AzureDataPlaneManagedIdentity)
	object.bitmap_ = b.bitmap_
	object.resourceID = b.resourceID
	return
}
