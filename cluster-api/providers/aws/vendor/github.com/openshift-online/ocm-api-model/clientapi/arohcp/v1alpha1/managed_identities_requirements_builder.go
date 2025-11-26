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

// Representation of managed identities requirements.
// When creating ARO-HCP Clusters, the end-users will need to pre-create the set of Managed Identities
// required by the clusters.
// The set of Managed Identities that the end-users need to precreate is not static and depends on
// several factors:
// (1) The OpenShift version of the cluster being created.
// (2) The functionalities that are being enabled for the cluster. Some Managed Identities are not
// always required but become required if a given functionality is enabled.
// Additionally, the Managed Identities that the end-users will need to precreate will have to have a
// set of required permissions assigned to them which also have to be returned to the end users.
type ManagedIdentitiesRequirementsBuilder struct {
	fieldSet_                       []bool
	id                              string
	href                            string
	controlPlaneOperatorsIdentities []*ControlPlaneOperatorIdentityRequirementBuilder
	dataPlaneOperatorsIdentities    []*DataPlaneOperatorIdentityRequirementBuilder
}

// NewManagedIdentitiesRequirements creates a new builder of 'managed_identities_requirements' objects.
func NewManagedIdentitiesRequirements() *ManagedIdentitiesRequirementsBuilder {
	return &ManagedIdentitiesRequirementsBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ManagedIdentitiesRequirementsBuilder) Link(value bool) *ManagedIdentitiesRequirementsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ManagedIdentitiesRequirementsBuilder) ID(value string) *ManagedIdentitiesRequirementsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ManagedIdentitiesRequirementsBuilder) HREF(value string) *ManagedIdentitiesRequirementsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagedIdentitiesRequirementsBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// ControlPlaneOperatorsIdentities sets the value of the 'control_plane_operators_identities' attribute to the given values.
func (b *ManagedIdentitiesRequirementsBuilder) ControlPlaneOperatorsIdentities(values ...*ControlPlaneOperatorIdentityRequirementBuilder) *ManagedIdentitiesRequirementsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.controlPlaneOperatorsIdentities = make([]*ControlPlaneOperatorIdentityRequirementBuilder, len(values))
	copy(b.controlPlaneOperatorsIdentities, values)
	b.fieldSet_[3] = true
	return b
}

// DataPlaneOperatorsIdentities sets the value of the 'data_plane_operators_identities' attribute to the given values.
func (b *ManagedIdentitiesRequirementsBuilder) DataPlaneOperatorsIdentities(values ...*DataPlaneOperatorIdentityRequirementBuilder) *ManagedIdentitiesRequirementsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.dataPlaneOperatorsIdentities = make([]*DataPlaneOperatorIdentityRequirementBuilder, len(values))
	copy(b.dataPlaneOperatorsIdentities, values)
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagedIdentitiesRequirementsBuilder) Copy(object *ManagedIdentitiesRequirements) *ManagedIdentitiesRequirementsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.controlPlaneOperatorsIdentities != nil {
		b.controlPlaneOperatorsIdentities = make([]*ControlPlaneOperatorIdentityRequirementBuilder, len(object.controlPlaneOperatorsIdentities))
		for i, v := range object.controlPlaneOperatorsIdentities {
			b.controlPlaneOperatorsIdentities[i] = NewControlPlaneOperatorIdentityRequirement().Copy(v)
		}
	} else {
		b.controlPlaneOperatorsIdentities = nil
	}
	if object.dataPlaneOperatorsIdentities != nil {
		b.dataPlaneOperatorsIdentities = make([]*DataPlaneOperatorIdentityRequirementBuilder, len(object.dataPlaneOperatorsIdentities))
		for i, v := range object.dataPlaneOperatorsIdentities {
			b.dataPlaneOperatorsIdentities[i] = NewDataPlaneOperatorIdentityRequirement().Copy(v)
		}
	} else {
		b.dataPlaneOperatorsIdentities = nil
	}
	return b
}

// Build creates a 'managed_identities_requirements' object using the configuration stored in the builder.
func (b *ManagedIdentitiesRequirementsBuilder) Build() (object *ManagedIdentitiesRequirements, err error) {
	object = new(ManagedIdentitiesRequirements)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.controlPlaneOperatorsIdentities != nil {
		object.controlPlaneOperatorsIdentities = make([]*ControlPlaneOperatorIdentityRequirement, len(b.controlPlaneOperatorsIdentities))
		for i, v := range b.controlPlaneOperatorsIdentities {
			object.controlPlaneOperatorsIdentities[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.dataPlaneOperatorsIdentities != nil {
		object.dataPlaneOperatorsIdentities = make([]*DataPlaneOperatorIdentityRequirement, len(b.dataPlaneOperatorsIdentities))
		for i, v := range b.dataPlaneOperatorsIdentities {
			object.dataPlaneOperatorsIdentities[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
