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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	time "time"
)

// OidcConfigBuilder contains the data and logic needed to build 'oidc_config' objects.
//
// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
type OidcConfigBuilder struct {
	bitmap_             uint32
	href                string
	id                  string
	creationTimestamp   time.Time
	installerRoleArn    string
	issuerUrl           string
	lastUpdateTimestamp time.Time
	lastUsedTimestamp   time.Time
	organizationId      string
	secretArn           string
	managed             bool
	reusable            bool
}

// NewOidcConfig creates a new builder of 'oidc_config' objects.
func NewOidcConfig() *OidcConfigBuilder {
	return &OidcConfigBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OidcConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// HREF sets the value of the 'HREF' attribute to the given value.
func (b *OidcConfigBuilder) HREF(value string) *OidcConfigBuilder {
	b.href = value
	b.bitmap_ |= 1
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *OidcConfigBuilder) ID(value string) *OidcConfigBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) CreationTimestamp(value time.Time) *OidcConfigBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 4
	return b
}

// InstallerRoleArn sets the value of the 'installer_role_arn' attribute to the given value.
func (b *OidcConfigBuilder) InstallerRoleArn(value string) *OidcConfigBuilder {
	b.installerRoleArn = value
	b.bitmap_ |= 8
	return b
}

// IssuerUrl sets the value of the 'issuer_url' attribute to the given value.
func (b *OidcConfigBuilder) IssuerUrl(value string) *OidcConfigBuilder {
	b.issuerUrl = value
	b.bitmap_ |= 16
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) LastUpdateTimestamp(value time.Time) *OidcConfigBuilder {
	b.lastUpdateTimestamp = value
	b.bitmap_ |= 32
	return b
}

// LastUsedTimestamp sets the value of the 'last_used_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) LastUsedTimestamp(value time.Time) *OidcConfigBuilder {
	b.lastUsedTimestamp = value
	b.bitmap_ |= 64
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *OidcConfigBuilder) Managed(value bool) *OidcConfigBuilder {
	b.managed = value
	b.bitmap_ |= 128
	return b
}

// OrganizationId sets the value of the 'organization_id' attribute to the given value.
func (b *OidcConfigBuilder) OrganizationId(value string) *OidcConfigBuilder {
	b.organizationId = value
	b.bitmap_ |= 256
	return b
}

// Reusable sets the value of the 'reusable' attribute to the given value.
func (b *OidcConfigBuilder) Reusable(value bool) *OidcConfigBuilder {
	b.reusable = value
	b.bitmap_ |= 512
	return b
}

// SecretArn sets the value of the 'secret_arn' attribute to the given value.
func (b *OidcConfigBuilder) SecretArn(value string) *OidcConfigBuilder {
	b.secretArn = value
	b.bitmap_ |= 1024
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OidcConfigBuilder) Copy(object *OidcConfig) *OidcConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.href = object.href
	b.id = object.id
	b.creationTimestamp = object.creationTimestamp
	b.installerRoleArn = object.installerRoleArn
	b.issuerUrl = object.issuerUrl
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	b.lastUsedTimestamp = object.lastUsedTimestamp
	b.managed = object.managed
	b.organizationId = object.organizationId
	b.reusable = object.reusable
	b.secretArn = object.secretArn
	return b
}

// Build creates a 'oidc_config' object using the configuration stored in the builder.
func (b *OidcConfigBuilder) Build() (object *OidcConfig, err error) {
	object = new(OidcConfig)
	object.bitmap_ = b.bitmap_
	object.href = b.href
	object.id = b.id
	object.creationTimestamp = b.creationTimestamp
	object.installerRoleArn = b.installerRoleArn
	object.issuerUrl = b.issuerUrl
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	object.lastUsedTimestamp = b.lastUsedTimestamp
	object.managed = b.managed
	object.organizationId = b.organizationId
	object.reusable = b.reusable
	object.secretArn = b.secretArn
	return
}
