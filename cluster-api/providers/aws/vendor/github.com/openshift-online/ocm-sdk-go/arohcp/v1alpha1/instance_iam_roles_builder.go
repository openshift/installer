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

// InstanceIAMRolesBuilder contains the data and logic needed to build 'instance_IAM_roles' objects.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type InstanceIAMRolesBuilder struct {
	bitmap_       uint32
	masterRoleARN string
	workerRoleARN string
}

// NewInstanceIAMRoles creates a new builder of 'instance_IAM_roles' objects.
func NewInstanceIAMRoles() *InstanceIAMRolesBuilder {
	return &InstanceIAMRolesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *InstanceIAMRolesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// MasterRoleARN sets the value of the 'master_role_ARN' attribute to the given value.
func (b *InstanceIAMRolesBuilder) MasterRoleARN(value string) *InstanceIAMRolesBuilder {
	b.masterRoleARN = value
	b.bitmap_ |= 1
	return b
}

// WorkerRoleARN sets the value of the 'worker_role_ARN' attribute to the given value.
func (b *InstanceIAMRolesBuilder) WorkerRoleARN(value string) *InstanceIAMRolesBuilder {
	b.workerRoleARN = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *InstanceIAMRolesBuilder) Copy(object *InstanceIAMRoles) *InstanceIAMRolesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.masterRoleARN = object.masterRoleARN
	b.workerRoleARN = object.workerRoleARN
	return b
}

// Build creates a 'instance_IAM_roles' object using the configuration stored in the builder.
func (b *InstanceIAMRolesBuilder) Build() (object *InstanceIAMRoles, err error) {
	object = new(InstanceIAMRoles)
	object.bitmap_ = b.bitmap_
	object.masterRoleARN = b.masterRoleARN
	object.workerRoleARN = b.workerRoleARN
	return
}
