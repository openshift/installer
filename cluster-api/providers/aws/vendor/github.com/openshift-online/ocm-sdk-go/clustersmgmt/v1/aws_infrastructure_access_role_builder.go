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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AWSInfrastructureAccessRoleBuilder contains the data and logic needed to build 'AWS_infrastructure_access_role' objects.
//
// A set of acces permissions for AWS resources
type AWSInfrastructureAccessRoleBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	description string
	displayName string
	state       AWSInfrastructureAccessRoleState
}

// NewAWSInfrastructureAccessRole creates a new builder of 'AWS_infrastructure_access_role' objects.
func NewAWSInfrastructureAccessRole() *AWSInfrastructureAccessRoleBuilder {
	return &AWSInfrastructureAccessRoleBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AWSInfrastructureAccessRoleBuilder) Link(value bool) *AWSInfrastructureAccessRoleBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AWSInfrastructureAccessRoleBuilder) ID(value string) *AWSInfrastructureAccessRoleBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AWSInfrastructureAccessRoleBuilder) HREF(value string) *AWSInfrastructureAccessRoleBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSInfrastructureAccessRoleBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AWSInfrastructureAccessRoleBuilder) Description(value string) *AWSInfrastructureAccessRoleBuilder {
	b.description = value
	b.bitmap_ |= 8
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *AWSInfrastructureAccessRoleBuilder) DisplayName(value string) *AWSInfrastructureAccessRoleBuilder {
	b.displayName = value
	b.bitmap_ |= 16
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// State of an AWS infrastructure access role.
func (b *AWSInfrastructureAccessRoleBuilder) State(value AWSInfrastructureAccessRoleState) *AWSInfrastructureAccessRoleBuilder {
	b.state = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSInfrastructureAccessRoleBuilder) Copy(object *AWSInfrastructureAccessRole) *AWSInfrastructureAccessRoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.description = object.description
	b.displayName = object.displayName
	b.state = object.state
	return b
}

// Build creates a 'AWS_infrastructure_access_role' object using the configuration stored in the builder.
func (b *AWSInfrastructureAccessRoleBuilder) Build() (object *AWSInfrastructureAccessRole, err error) {
	object = new(AWSInfrastructureAccessRole)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.description = b.description
	object.displayName = b.displayName
	object.state = b.state
	return
}
