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

// Contains the necessary attributes to support role-based authentication on AWS.
type InstanceIAMRolesBuilder struct {
	fieldSet_     []bool
	masterRoleARN string
	workerRoleARN string
}

// NewInstanceIAMRoles creates a new builder of 'instance_IAM_roles' objects.
func NewInstanceIAMRoles() *InstanceIAMRolesBuilder {
	return &InstanceIAMRolesBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *InstanceIAMRolesBuilder) Empty() bool {
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

// MasterRoleARN sets the value of the 'master_role_ARN' attribute to the given value.
func (b *InstanceIAMRolesBuilder) MasterRoleARN(value string) *InstanceIAMRolesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.masterRoleARN = value
	b.fieldSet_[0] = true
	return b
}

// WorkerRoleARN sets the value of the 'worker_role_ARN' attribute to the given value.
func (b *InstanceIAMRolesBuilder) WorkerRoleARN(value string) *InstanceIAMRolesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.workerRoleARN = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *InstanceIAMRolesBuilder) Copy(object *InstanceIAMRoles) *InstanceIAMRolesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.masterRoleARN = object.masterRoleARN
	b.workerRoleARN = object.workerRoleARN
	return b
}

// Build creates a 'instance_IAM_roles' object using the configuration stored in the builder.
func (b *InstanceIAMRolesBuilder) Build() (object *InstanceIAMRoles, err error) {
	object = new(InstanceIAMRoles)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.masterRoleARN = b.masterRoleARN
	object.workerRoleARN = b.workerRoleARN
	return
}
