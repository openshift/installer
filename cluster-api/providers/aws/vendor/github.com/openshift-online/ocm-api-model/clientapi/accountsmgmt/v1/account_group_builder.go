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

type AccountGroupBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	createdAt      time.Time
	description    string
	externalID     string
	managedBy      AccountGroupManagedBy
	name           string
	organizationID string
	updatedAt      time.Time
}

// NewAccountGroup creates a new builder of 'account_group' objects.
func NewAccountGroup() *AccountGroupBuilder {
	return &AccountGroupBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AccountGroupBuilder) Link(value bool) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AccountGroupBuilder) ID(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AccountGroupBuilder) HREF(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccountGroupBuilder) Empty() bool {
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AccountGroupBuilder) CreatedAt(value time.Time) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AccountGroupBuilder) Description(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.description = value
	b.fieldSet_[4] = true
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *AccountGroupBuilder) ExternalID(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.externalID = value
	b.fieldSet_[5] = true
	return b
}

// ManagedBy sets the value of the 'managed_by' attribute to the given value.
func (b *AccountGroupBuilder) ManagedBy(value AccountGroupManagedBy) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.managedBy = value
	b.fieldSet_[6] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AccountGroupBuilder) Name(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.name = value
	b.fieldSet_[7] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *AccountGroupBuilder) OrganizationID(value string) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.organizationID = value
	b.fieldSet_[8] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AccountGroupBuilder) UpdatedAt(value time.Time) *AccountGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccountGroupBuilder) Copy(object *AccountGroup) *AccountGroupBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.description = object.description
	b.externalID = object.externalID
	b.managedBy = object.managedBy
	b.name = object.name
	b.organizationID = object.organizationID
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'account_group' object using the configuration stored in the builder.
func (b *AccountGroupBuilder) Build() (object *AccountGroup, err error) {
	object = new(AccountGroup)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.description = b.description
	object.externalID = b.externalID
	object.managedBy = b.managedBy
	object.name = b.name
	object.organizationID = b.organizationID
	object.updatedAt = b.updatedAt
	return
}
