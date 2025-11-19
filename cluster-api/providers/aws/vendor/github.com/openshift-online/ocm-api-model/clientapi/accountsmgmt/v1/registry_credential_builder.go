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

type RegistryCredentialBuilder struct {
	fieldSet_          []bool
	id                 string
	href               string
	account            *AccountBuilder
	createdAt          time.Time
	externalResourceID string
	registry           *RegistryBuilder
	token              string
	updatedAt          time.Time
	username           string
}

// NewRegistryCredential creates a new builder of 'registry_credential' objects.
func NewRegistryCredential() *RegistryCredentialBuilder {
	return &RegistryCredentialBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *RegistryCredentialBuilder) Link(value bool) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *RegistryCredentialBuilder) ID(value string) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *RegistryCredentialBuilder) HREF(value string) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryCredentialBuilder) Empty() bool {
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

// Account sets the value of the 'account' attribute to the given value.
func (b *RegistryCredentialBuilder) Account(value *AccountBuilder) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.account = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RegistryCredentialBuilder) CreatedAt(value time.Time) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// ExternalResourceID sets the value of the 'external_resource_ID' attribute to the given value.
func (b *RegistryCredentialBuilder) ExternalResourceID(value string) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.externalResourceID = value
	b.fieldSet_[5] = true
	return b
}

// Registry sets the value of the 'registry' attribute to the given value.
func (b *RegistryCredentialBuilder) Registry(value *RegistryBuilder) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.registry = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Token sets the value of the 'token' attribute to the given value.
func (b *RegistryCredentialBuilder) Token(value string) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.token = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RegistryCredentialBuilder) UpdatedAt(value time.Time) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *RegistryCredentialBuilder) Username(value string) *RegistryCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.username = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryCredentialBuilder) Copy(object *RegistryCredential) *RegistryCredentialBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.account != nil {
		b.account = NewAccount().Copy(object.account)
	} else {
		b.account = nil
	}
	b.createdAt = object.createdAt
	b.externalResourceID = object.externalResourceID
	if object.registry != nil {
		b.registry = NewRegistry().Copy(object.registry)
	} else {
		b.registry = nil
	}
	b.token = object.token
	b.updatedAt = object.updatedAt
	b.username = object.username
	return b
}

// Build creates a 'registry_credential' object using the configuration stored in the builder.
func (b *RegistryCredentialBuilder) Build() (object *RegistryCredential, err error) {
	object = new(RegistryCredential)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.account != nil {
		object.account, err = b.account.Build()
		if err != nil {
			return
		}
	}
	object.createdAt = b.createdAt
	object.externalResourceID = b.externalResourceID
	if b.registry != nil {
		object.registry, err = b.registry.Build()
		if err != nil {
			return
		}
	}
	object.token = b.token
	object.updatedAt = b.updatedAt
	object.username = b.username
	return
}
