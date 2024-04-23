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

import (
	time "time"
)

// BreakGlassCredentialBuilder contains the data and logic needed to build 'break_glass_credential' objects.
//
// Representation of a break glass credential.
type BreakGlassCredentialBuilder struct {
	bitmap_             uint32
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
	return &BreakGlassCredentialBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *BreakGlassCredentialBuilder) Link(value bool) *BreakGlassCredentialBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *BreakGlassCredentialBuilder) ID(value string) *BreakGlassCredentialBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *BreakGlassCredentialBuilder) HREF(value string) *BreakGlassCredentialBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *BreakGlassCredentialBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ExpirationTimestamp sets the value of the 'expiration_timestamp' attribute to the given value.
func (b *BreakGlassCredentialBuilder) ExpirationTimestamp(value time.Time) *BreakGlassCredentialBuilder {
	b.expirationTimestamp = value
	b.bitmap_ |= 8
	return b
}

// Kubeconfig sets the value of the 'kubeconfig' attribute to the given value.
func (b *BreakGlassCredentialBuilder) Kubeconfig(value string) *BreakGlassCredentialBuilder {
	b.kubeconfig = value
	b.bitmap_ |= 16
	return b
}

// RevocationTimestamp sets the value of the 'revocation_timestamp' attribute to the given value.
func (b *BreakGlassCredentialBuilder) RevocationTimestamp(value time.Time) *BreakGlassCredentialBuilder {
	b.revocationTimestamp = value
	b.bitmap_ |= 32
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Status of the break glass credential.
func (b *BreakGlassCredentialBuilder) Status(value BreakGlassCredentialStatus) *BreakGlassCredentialBuilder {
	b.status = value
	b.bitmap_ |= 64
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *BreakGlassCredentialBuilder) Username(value string) *BreakGlassCredentialBuilder {
	b.username = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *BreakGlassCredentialBuilder) Copy(object *BreakGlassCredential) *BreakGlassCredentialBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.expirationTimestamp = b.expirationTimestamp
	object.kubeconfig = b.kubeconfig
	object.revocationTimestamp = b.revocationTimestamp
	object.status = b.status
	object.username = b.username
	return
}
