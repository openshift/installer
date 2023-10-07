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

// AWSInfrastructureAccessRoleGrantBuilder contains the data and logic needed to build 'AWS_infrastructure_access_role_grant' objects.
//
// Representation of an AWS infrastructure access role grant.
type AWSInfrastructureAccessRoleGrantBuilder struct {
	bitmap_          uint32
	id               string
	href             string
	consoleURL       string
	role             *AWSInfrastructureAccessRoleBuilder
	state            AWSInfrastructureAccessRoleGrantState
	stateDescription string
	userARN          string
}

// NewAWSInfrastructureAccessRoleGrant creates a new builder of 'AWS_infrastructure_access_role_grant' objects.
func NewAWSInfrastructureAccessRoleGrant() *AWSInfrastructureAccessRoleGrantBuilder {
	return &AWSInfrastructureAccessRoleGrantBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Link(value bool) *AWSInfrastructureAccessRoleGrantBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AWSInfrastructureAccessRoleGrantBuilder) ID(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AWSInfrastructureAccessRoleGrantBuilder) HREF(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) ConsoleURL(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	b.consoleURL = value
	b.bitmap_ |= 8
	return b
}

// Role sets the value of the 'role' attribute to the given value.
//
// A set of acces permissions for AWS resources
func (b *AWSInfrastructureAccessRoleGrantBuilder) Role(value *AWSInfrastructureAccessRoleBuilder) *AWSInfrastructureAccessRoleGrantBuilder {
	b.role = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// State of an AWS infrastructure access role grant.
func (b *AWSInfrastructureAccessRoleGrantBuilder) State(value AWSInfrastructureAccessRoleGrantState) *AWSInfrastructureAccessRoleGrantBuilder {
	b.state = value
	b.bitmap_ |= 32
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) StateDescription(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	b.stateDescription = value
	b.bitmap_ |= 64
	return b
}

// UserARN sets the value of the 'user_ARN' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) UserARN(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	b.userARN = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Copy(object *AWSInfrastructureAccessRoleGrant) *AWSInfrastructureAccessRoleGrantBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.consoleURL = object.consoleURL
	if object.role != nil {
		b.role = NewAWSInfrastructureAccessRole().Copy(object.role)
	} else {
		b.role = nil
	}
	b.state = object.state
	b.stateDescription = object.stateDescription
	b.userARN = object.userARN
	return b
}

// Build creates a 'AWS_infrastructure_access_role_grant' object using the configuration stored in the builder.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Build() (object *AWSInfrastructureAccessRoleGrant, err error) {
	object = new(AWSInfrastructureAccessRoleGrant)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.consoleURL = b.consoleURL
	if b.role != nil {
		object.role, err = b.role.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	object.stateDescription = b.stateDescription
	object.userARN = b.userARN
	return
}
