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

// RegistryCredentialBuilder contains the data and logic needed to build 'registry_credential' objects.
type RegistryCredentialBuilder struct {
	bitmap_            uint32
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
	return &RegistryCredentialBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *RegistryCredentialBuilder) Link(value bool) *RegistryCredentialBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *RegistryCredentialBuilder) ID(value string) *RegistryCredentialBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *RegistryCredentialBuilder) HREF(value string) *RegistryCredentialBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryCredentialBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Account sets the value of the 'account' attribute to the given value.
func (b *RegistryCredentialBuilder) Account(value *AccountBuilder) *RegistryCredentialBuilder {
	b.account = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RegistryCredentialBuilder) CreatedAt(value time.Time) *RegistryCredentialBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// ExternalResourceID sets the value of the 'external_resource_ID' attribute to the given value.
func (b *RegistryCredentialBuilder) ExternalResourceID(value string) *RegistryCredentialBuilder {
	b.externalResourceID = value
	b.bitmap_ |= 32
	return b
}

// Registry sets the value of the 'registry' attribute to the given value.
func (b *RegistryCredentialBuilder) Registry(value *RegistryBuilder) *RegistryCredentialBuilder {
	b.registry = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Token sets the value of the 'token' attribute to the given value.
func (b *RegistryCredentialBuilder) Token(value string) *RegistryCredentialBuilder {
	b.token = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RegistryCredentialBuilder) UpdatedAt(value time.Time) *RegistryCredentialBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *RegistryCredentialBuilder) Username(value string) *RegistryCredentialBuilder {
	b.username = value
	b.bitmap_ |= 512
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryCredentialBuilder) Copy(object *RegistryCredential) *RegistryCredentialBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
