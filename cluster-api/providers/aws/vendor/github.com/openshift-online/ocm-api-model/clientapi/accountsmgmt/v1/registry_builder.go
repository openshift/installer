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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	time "time"
)

type RegistryBuilder struct {
	fieldSet_  []bool
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
	return &RegistryBuilder{
		fieldSet_: make([]bool, 11),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *RegistryBuilder) Link(value bool) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *RegistryBuilder) ID(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *RegistryBuilder) HREF(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryBuilder) Empty() bool {
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

// URL sets the value of the 'URL' attribute to the given value.
func (b *RegistryBuilder) URL(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.url = value
	b.fieldSet_[3] = true
	return b
}

// CloudAlias sets the value of the 'cloud_alias' attribute to the given value.
func (b *RegistryBuilder) CloudAlias(value bool) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.cloudAlias = value
	b.fieldSet_[4] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RegistryBuilder) CreatedAt(value time.Time) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.createdAt = value
	b.fieldSet_[5] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *RegistryBuilder) Name(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.name = value
	b.fieldSet_[6] = true
	return b
}

// OrgName sets the value of the 'org_name' attribute to the given value.
func (b *RegistryBuilder) OrgName(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.orgName = value
	b.fieldSet_[7] = true
	return b
}

// TeamName sets the value of the 'team_name' attribute to the given value.
func (b *RegistryBuilder) TeamName(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.teamName = value
	b.fieldSet_[8] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RegistryBuilder) Type(value string) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.type_ = value
	b.fieldSet_[9] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RegistryBuilder) UpdatedAt(value time.Time) *RegistryBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.updatedAt = value
	b.fieldSet_[10] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryBuilder) Copy(object *Registry) *RegistryBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
