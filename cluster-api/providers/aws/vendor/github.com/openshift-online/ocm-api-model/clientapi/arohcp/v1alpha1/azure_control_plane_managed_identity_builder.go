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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// Represents the information associated to an Azure User-Assigned
// Managed Identity belonging to the Control Plane of the cluster.
type AzureControlPlaneManagedIdentityBuilder struct {
	fieldSet_   []bool
	clientID    string
	principalID string
	resourceID  string
}

// NewAzureControlPlaneManagedIdentity creates a new builder of 'azure_control_plane_managed_identity' objects.
func NewAzureControlPlaneManagedIdentity() *AzureControlPlaneManagedIdentityBuilder {
	return &AzureControlPlaneManagedIdentityBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureControlPlaneManagedIdentityBuilder) Empty() bool {
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

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) ClientID(value string) *AzureControlPlaneManagedIdentityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.clientID = value
	b.fieldSet_[0] = true
	return b
}

// PrincipalID sets the value of the 'principal_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) PrincipalID(value string) *AzureControlPlaneManagedIdentityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.principalID = value
	b.fieldSet_[1] = true
	return b
}

// ResourceID sets the value of the 'resource_ID' attribute to the given value.
func (b *AzureControlPlaneManagedIdentityBuilder) ResourceID(value string) *AzureControlPlaneManagedIdentityBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.resourceID = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureControlPlaneManagedIdentityBuilder) Copy(object *AzureControlPlaneManagedIdentity) *AzureControlPlaneManagedIdentityBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.clientID = object.clientID
	b.principalID = object.principalID
	b.resourceID = object.resourceID
	return b
}

// Build creates a 'azure_control_plane_managed_identity' object using the configuration stored in the builder.
func (b *AzureControlPlaneManagedIdentityBuilder) Build() (object *AzureControlPlaneManagedIdentity, err error) {
	object = new(AzureControlPlaneManagedIdentity)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clientID = b.clientID
	object.principalID = b.principalID
	object.resourceID = b.resourceID
	return
}
