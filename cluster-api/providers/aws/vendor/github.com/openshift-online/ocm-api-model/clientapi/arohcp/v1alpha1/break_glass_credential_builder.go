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

import (
	time "time"
)

// Representation of a break glass credential.
type BreakGlassCredentialBuilder struct {
	fieldSet_           []bool
	id                  string
	href                string
	expirationTimestamp time.Time
	kubeconfig          string
	revocationTimestamp time.Time
	status              BreakGlassCredentialStatus
	username            string
}

// NewBreakGlassCredential creates a new builder of 'break_glass_credential' objects.
func NewBreakGlassCredential() *BreakGlassCredentialBuilder {
	return &BreakGlassCredentialBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *BreakGlassCredentialBuilder) Link(value bool) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *BreakGlassCredentialBuilder) ID(value string) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *BreakGlassCredentialBuilder) HREF(value string) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *BreakGlassCredentialBuilder) Empty() bool {
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

// ExpirationTimestamp sets the value of the 'expiration_timestamp' attribute to the given value.
func (b *BreakGlassCredentialBuilder) ExpirationTimestamp(value time.Time) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.expirationTimestamp = value
	b.fieldSet_[3] = true
	return b
}

// Kubeconfig sets the value of the 'kubeconfig' attribute to the given value.
func (b *BreakGlassCredentialBuilder) Kubeconfig(value string) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.kubeconfig = value
	b.fieldSet_[4] = true
	return b
}

// RevocationTimestamp sets the value of the 'revocation_timestamp' attribute to the given value.
func (b *BreakGlassCredentialBuilder) RevocationTimestamp(value time.Time) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.revocationTimestamp = value
	b.fieldSet_[5] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Status of the break glass credential.
func (b *BreakGlassCredentialBuilder) Status(value BreakGlassCredentialStatus) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.status = value
	b.fieldSet_[6] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *BreakGlassCredentialBuilder) Username(value string) *BreakGlassCredentialBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.username = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *BreakGlassCredentialBuilder) Copy(object *BreakGlassCredential) *BreakGlassCredentialBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.expirationTimestamp = object.expirationTimestamp
	b.kubeconfig = object.kubeconfig
	b.revocationTimestamp = object.revocationTimestamp
	b.status = object.status
	b.username = object.username
	return b
}

// Build creates a 'break_glass_credential' object using the configuration stored in the builder.
func (b *BreakGlassCredentialBuilder) Build() (object *BreakGlassCredential, err error) {
	object = new(BreakGlassCredential)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.expirationTimestamp = b.expirationTimestamp
	object.kubeconfig = b.kubeconfig
	object.revocationTimestamp = b.revocationTimestamp
	object.status = b.status
	object.username = b.username
	return
}
