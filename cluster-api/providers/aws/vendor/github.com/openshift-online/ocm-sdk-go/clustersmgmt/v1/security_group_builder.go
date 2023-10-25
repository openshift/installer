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

// SecurityGroupBuilder contains the data and logic needed to build 'security_group' objects.
//
// AWS security group object
type SecurityGroupBuilder struct {
	bitmap_       uint32
	id            string
	name          string
	redHatManaged bool
}

// NewSecurityGroup creates a new builder of 'security_group' objects.
func NewSecurityGroup() *SecurityGroupBuilder {
	return &SecurityGroupBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SecurityGroupBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *SecurityGroupBuilder) ID(value string) *SecurityGroupBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *SecurityGroupBuilder) Name(value string) *SecurityGroupBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// RedHatManaged sets the value of the 'red_hat_managed' attribute to the given value.
func (b *SecurityGroupBuilder) RedHatManaged(value bool) *SecurityGroupBuilder {
	b.redHatManaged = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SecurityGroupBuilder) Copy(object *SecurityGroup) *SecurityGroupBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.name = object.name
	b.redHatManaged = object.redHatManaged
	return b
}

// Build creates a 'security_group' object using the configuration stored in the builder.
func (b *SecurityGroupBuilder) Build() (object *SecurityGroup, err error) {
	object = new(SecurityGroup)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.name = b.name
	object.redHatManaged = b.redHatManaged
	return
}
