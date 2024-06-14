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

// RolePolicyBuilder contains the data and logic needed to build 'role_policy' objects.
type RolePolicyBuilder struct {
	bitmap_ uint32
	arn     string
	name    string
	type_   string
}

// NewRolePolicy creates a new builder of 'role_policy' objects.
func NewRolePolicy() *RolePolicyBuilder {
	return &RolePolicyBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RolePolicyBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Arn sets the value of the 'arn' attribute to the given value.
func (b *RolePolicyBuilder) Arn(value string) *RolePolicyBuilder {
	b.arn = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RolePolicyBuilder) Name(value string) *RolePolicyBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RolePolicyBuilder) Type(value string) *RolePolicyBuilder {
	b.type_ = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RolePolicyBuilder) Copy(object *RolePolicy) *RolePolicyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.arn = object.arn
	b.name = object.name
	b.type_ = object.type_
	return b
}

// Build creates a 'role_policy' object using the configuration stored in the builder.
func (b *RolePolicyBuilder) Build() (object *RolePolicy, err error) {
	object = new(RolePolicy)
	object.bitmap_ = b.bitmap_
	object.arn = b.arn
	object.name = b.name
	object.type_ = b.type_
	return
}
