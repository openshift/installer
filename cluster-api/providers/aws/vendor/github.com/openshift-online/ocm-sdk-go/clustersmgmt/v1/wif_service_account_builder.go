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

// WifServiceAccountBuilder contains the data and logic needed to build 'wif_service_account' objects.
type WifServiceAccountBuilder struct {
	bitmap_           uint32
	accessMethod      WifAccessMethod
	credentialRequest *WifCredentialRequestBuilder
	osdRole           string
	roles             []*WifRoleBuilder
	serviceAccountId  string
}

// NewWifServiceAccount creates a new builder of 'wif_service_account' objects.
func NewWifServiceAccount() *WifServiceAccountBuilder {
	return &WifServiceAccountBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifServiceAccountBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccessMethod sets the value of the 'access_method' attribute to the given value.
func (b *WifServiceAccountBuilder) AccessMethod(value WifAccessMethod) *WifServiceAccountBuilder {
	b.accessMethod = value
	b.bitmap_ |= 1
	return b
}

// CredentialRequest sets the value of the 'credential_request' attribute to the given value.
func (b *WifServiceAccountBuilder) CredentialRequest(value *WifCredentialRequestBuilder) *WifServiceAccountBuilder {
	b.credentialRequest = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// OsdRole sets the value of the 'osd_role' attribute to the given value.
func (b *WifServiceAccountBuilder) OsdRole(value string) *WifServiceAccountBuilder {
	b.osdRole = value
	b.bitmap_ |= 4
	return b
}

// Roles sets the value of the 'roles' attribute to the given values.
func (b *WifServiceAccountBuilder) Roles(values ...*WifRoleBuilder) *WifServiceAccountBuilder {
	b.roles = make([]*WifRoleBuilder, len(values))
	copy(b.roles, values)
	b.bitmap_ |= 8
	return b
}

// ServiceAccountId sets the value of the 'service_account_id' attribute to the given value.
func (b *WifServiceAccountBuilder) ServiceAccountId(value string) *WifServiceAccountBuilder {
	b.serviceAccountId = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifServiceAccountBuilder) Copy(object *WifServiceAccount) *WifServiceAccountBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accessMethod = object.accessMethod
	if object.credentialRequest != nil {
		b.credentialRequest = NewWifCredentialRequest().Copy(object.credentialRequest)
	} else {
		b.credentialRequest = nil
	}
	b.osdRole = object.osdRole
	if object.roles != nil {
		b.roles = make([]*WifRoleBuilder, len(object.roles))
		for i, v := range object.roles {
			b.roles[i] = NewWifRole().Copy(v)
		}
	} else {
		b.roles = nil
	}
	b.serviceAccountId = object.serviceAccountId
	return b
}

// Build creates a 'wif_service_account' object using the configuration stored in the builder.
func (b *WifServiceAccountBuilder) Build() (object *WifServiceAccount, err error) {
	object = new(WifServiceAccount)
	object.bitmap_ = b.bitmap_
	object.accessMethod = b.accessMethod
	if b.credentialRequest != nil {
		object.credentialRequest, err = b.credentialRequest.Build()
		if err != nil {
			return
		}
	}
	object.osdRole = b.osdRole
	if b.roles != nil {
		object.roles = make([]*WifRole, len(b.roles))
		for i, v := range b.roles {
			object.roles[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.serviceAccountId = b.serviceAccountId
	return
}
