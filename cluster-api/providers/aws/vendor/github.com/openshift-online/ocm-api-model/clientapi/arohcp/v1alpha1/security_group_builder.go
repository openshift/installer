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

// AWS security group object
type SecurityGroupBuilder struct {
	fieldSet_     []bool
	id            string
	name          string
	redHatManaged bool
}

// NewSecurityGroup creates a new builder of 'security_group' objects.
func NewSecurityGroup() *SecurityGroupBuilder {
	return &SecurityGroupBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SecurityGroupBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *SecurityGroupBuilder) ID(value string) *SecurityGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *SecurityGroupBuilder) Name(value string) *SecurityGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// RedHatManaged sets the value of the 'red_hat_managed' attribute to the given value.
func (b *SecurityGroupBuilder) RedHatManaged(value bool) *SecurityGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.redHatManaged = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SecurityGroupBuilder) Copy(object *SecurityGroup) *SecurityGroupBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.name = object.name
	b.redHatManaged = object.redHatManaged
	return b
}

// Build creates a 'security_group' object using the configuration stored in the builder.
func (b *SecurityGroupBuilder) Build() (object *SecurityGroup, err error) {
	object = new(SecurityGroup)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.name = b.name
	object.redHatManaged = b.redHatManaged
	return
}
