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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	time "time"
)

// RegistryBuilder contains the data and logic needed to build 'registry' objects.
type RegistryBuilder struct {
	bitmap_    uint32
	id         string
	href       string
	url        string
	createdAt  time.Time
	name       string
	orgName    string
	teamName   string
	type_      string
	updatedAt  time.Time
	cloudAlias bool
}

// NewRegistry creates a new builder of 'registry' objects.
func NewRegistry() *RegistryBuilder {
	return &RegistryBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *RegistryBuilder) Link(value bool) *RegistryBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *RegistryBuilder) ID(value string) *RegistryBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *RegistryBuilder) HREF(value string) *RegistryBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *RegistryBuilder) URL(value string) *RegistryBuilder {
	b.url = value
	b.bitmap_ |= 8
	return b
}

// CloudAlias sets the value of the 'cloud_alias' attribute to the given value.
func (b *RegistryBuilder) CloudAlias(value bool) *RegistryBuilder {
	b.cloudAlias = value
	b.bitmap_ |= 16
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RegistryBuilder) CreatedAt(value time.Time) *RegistryBuilder {
	b.createdAt = value
	b.bitmap_ |= 32
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RegistryBuilder) Name(value string) *RegistryBuilder {
	b.name = value
	b.bitmap_ |= 64
	return b
}

// OrgName sets the value of the 'org_name' attribute to the given value.
func (b *RegistryBuilder) OrgName(value string) *RegistryBuilder {
	b.orgName = value
	b.bitmap_ |= 128
	return b
}

// TeamName sets the value of the 'team_name' attribute to the given value.
func (b *RegistryBuilder) TeamName(value string) *RegistryBuilder {
	b.teamName = value
	b.bitmap_ |= 256
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RegistryBuilder) Type(value string) *RegistryBuilder {
	b.type_ = value
	b.bitmap_ |= 512
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RegistryBuilder) UpdatedAt(value time.Time) *RegistryBuilder {
	b.updatedAt = value
	b.bitmap_ |= 1024
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryBuilder) Copy(object *Registry) *RegistryBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.url = object.url
	b.cloudAlias = object.cloudAlias
	b.createdAt = object.createdAt
	b.name = object.name
	b.orgName = object.orgName
	b.teamName = object.teamName
	b.type_ = object.type_
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'registry' object using the configuration stored in the builder.
func (b *RegistryBuilder) Build() (object *Registry, err error) {
	object = new(Registry)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.url = b.url
	object.cloudAlias = b.cloudAlias
	object.createdAt = b.createdAt
	object.name = b.name
	object.orgName = b.orgName
	object.teamName = b.teamName
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	return
}
