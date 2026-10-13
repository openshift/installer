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

// Representation of an AWS infrastructure access role grant.
type AWSInfrastructureAccessRoleGrantBuilder struct {
	fieldSet_        []bool
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
	return &AWSInfrastructureAccessRoleGrantBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Link(value bool) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AWSInfrastructureAccessRoleGrantBuilder) ID(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AWSInfrastructureAccessRoleGrantBuilder) HREF(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Empty() bool {
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

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) ConsoleURL(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.consoleURL = value
	b.fieldSet_[3] = true
	return b
}

// Role sets the value of the 'role' attribute to the given value.
//
// A set of acces permissions for AWS resources
func (b *AWSInfrastructureAccessRoleGrantBuilder) Role(value *AWSInfrastructureAccessRoleBuilder) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.role = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// State of an AWS infrastructure access role grant.
func (b *AWSInfrastructureAccessRoleGrantBuilder) State(value AWSInfrastructureAccessRoleGrantState) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.state = value
	b.fieldSet_[5] = true
	return b
}

// StateDescription sets the value of the 'state_description' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) StateDescription(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.stateDescription = value
	b.fieldSet_[6] = true
	return b
}

// UserARN sets the value of the 'user_ARN' attribute to the given value.
func (b *AWSInfrastructureAccessRoleGrantBuilder) UserARN(value string) *AWSInfrastructureAccessRoleGrantBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.userARN = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSInfrastructureAccessRoleGrantBuilder) Copy(object *AWSInfrastructureAccessRoleGrant) *AWSInfrastructureAccessRoleGrantBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
