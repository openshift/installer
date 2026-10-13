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

type AccountBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	banCode        string
	banDescription string
	capabilities   []*CapabilityBuilder
	createdAt      time.Time
	email          string
	firstName      string
	labels         []*LabelBuilder
	lastName       string
	organization   *OrganizationBuilder
	rhitAccountID  string
	rhitWebUserId  string
	updatedAt      time.Time
	username       string
	banned         bool
	serviceAccount bool
}

// NewAccount creates a new builder of 'account' objects.
func NewAccount() *AccountBuilder {
	return &AccountBuilder{
		fieldSet_: make([]bool, 18),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AccountBuilder) Link(value bool) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AccountBuilder) ID(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AccountBuilder) HREF(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccountBuilder) Empty() bool {
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

// BanCode sets the value of the 'ban_code' attribute to the given value.
func (b *AccountBuilder) BanCode(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.banCode = value
	b.fieldSet_[3] = true
	return b
}

// BanDescription sets the value of the 'ban_description' attribute to the given value.
func (b *AccountBuilder) BanDescription(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.banDescription = value
	b.fieldSet_[4] = true
	return b
}

// Banned sets the value of the 'banned' attribute to the given value.
func (b *AccountBuilder) Banned(value bool) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.banned = value
	b.fieldSet_[5] = true
	return b
}

// Capabilities sets the value of the 'capabilities' attribute to the given values.
func (b *AccountBuilder) Capabilities(values ...*CapabilityBuilder) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.capabilities = make([]*CapabilityBuilder, len(values))
	copy(b.capabilities, values)
	b.fieldSet_[6] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AccountBuilder) CreatedAt(value time.Time) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.createdAt = value
	b.fieldSet_[7] = true
	return b
}

// Email sets the value of the 'email' attribute to the given value.
func (b *AccountBuilder) Email(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.email = value
	b.fieldSet_[8] = true
	return b
}

// FirstName sets the value of the 'first_name' attribute to the given value.
func (b *AccountBuilder) FirstName(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.firstName = value
	b.fieldSet_[9] = true
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *AccountBuilder) Labels(values ...*LabelBuilder) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.labels = make([]*LabelBuilder, len(values))
	copy(b.labels, values)
	b.fieldSet_[10] = true
	return b
}

// LastName sets the value of the 'last_name' attribute to the given value.
func (b *AccountBuilder) LastName(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.lastName = value
	b.fieldSet_[11] = true
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
func (b *AccountBuilder) Organization(value *OrganizationBuilder) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.organization = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// RhitAccountID sets the value of the 'rhit_account_ID' attribute to the given value.
func (b *AccountBuilder) RhitAccountID(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.rhitAccountID = value
	b.fieldSet_[13] = true
	return b
}

// RhitWebUserId sets the value of the 'rhit_web_user_id' attribute to the given value.
func (b *AccountBuilder) RhitWebUserId(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.rhitWebUserId = value
	b.fieldSet_[14] = true
	return b
}

// ServiceAccount sets the value of the 'service_account' attribute to the given value.
func (b *AccountBuilder) ServiceAccount(value bool) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.serviceAccount = value
	b.fieldSet_[15] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AccountBuilder) UpdatedAt(value time.Time) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.updatedAt = value
	b.fieldSet_[16] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *AccountBuilder) Username(value string) *AccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.username = value
	b.fieldSet_[17] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccountBuilder) Copy(object *Account) *AccountBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.banCode = object.banCode
	b.banDescription = object.banDescription
	b.banned = object.banned
	if object.capabilities != nil {
		b.capabilities = make([]*CapabilityBuilder, len(object.capabilities))
		for i, v := range object.capabilities {
			b.capabilities[i] = NewCapability().Copy(v)
		}
	} else {
		b.capabilities = nil
	}
	b.createdAt = object.createdAt
	b.email = object.email
	b.firstName = object.firstName
	if object.labels != nil {
		b.labels = make([]*LabelBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabel().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.lastName = object.lastName
	if object.organization != nil {
		b.organization = NewOrganization().Copy(object.organization)
	} else {
		b.organization = nil
	}
	b.rhitAccountID = object.rhitAccountID
	b.rhitWebUserId = object.rhitWebUserId
	b.serviceAccount = object.serviceAccount
	b.updatedAt = object.updatedAt
	b.username = object.username
	return b
}

// Build creates a 'account' object using the configuration stored in the builder.
func (b *AccountBuilder) Build() (object *Account, err error) {
	object = new(Account)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.banCode = b.banCode
	object.banDescription = b.banDescription
	object.banned = b.banned
	if b.capabilities != nil {
		object.capabilities = make([]*Capability, len(b.capabilities))
		for i, v := range b.capabilities {
			object.capabilities[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.createdAt = b.createdAt
	object.email = b.email
	object.firstName = b.firstName
	if b.labels != nil {
		object.labels = make([]*Label, len(b.labels))
		for i, v := range b.labels {
			object.labels[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.lastName = b.lastName
	if b.organization != nil {
		object.organization, err = b.organization.Build()
		if err != nil {
			return
		}
	}
	object.rhitAccountID = b.rhitAccountID
	object.rhitWebUserId = b.rhitWebUserId
	object.serviceAccount = b.serviceAccount
	object.updatedAt = b.updatedAt
	object.username = b.username
	return
}
