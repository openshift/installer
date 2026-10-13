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

type RolePolicyBuilder struct {
	fieldSet_ []bool
	arn       string
	name      string
	type_     string
}

// NewRolePolicy creates a new builder of 'role_policy' objects.
func NewRolePolicy() *RolePolicyBuilder {
	return &RolePolicyBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RolePolicyBuilder) Empty() bool {
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

// Arn sets the value of the 'arn' attribute to the given value.
func (b *RolePolicyBuilder) Arn(value string) *RolePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.arn = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RolePolicyBuilder) Name(value string) *RolePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RolePolicyBuilder) Type(value string) *RolePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.type_ = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RolePolicyBuilder) Copy(object *RolePolicy) *RolePolicyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.arn = object.arn
	b.name = object.name
	b.type_ = object.type_
	return b
}

// Build creates a 'role_policy' object using the configuration stored in the builder.
func (b *RolePolicyBuilder) Build() (object *RolePolicy, err error) {
	object = new(RolePolicy)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.arn = b.arn
	object.name = b.name
	object.type_ = b.type_
	return
}
