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

// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
type OidcConfigBuilder struct {
	fieldSet_           []bool
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
	return &OidcConfigBuilder{
		fieldSet_: make([]bool, 11),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OidcConfigBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// HREF sets the value of the 'HREF' attribute to the given value.
func (b *OidcConfigBuilder) HREF(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.href = value
	b.fieldSet_[0] = true
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *OidcConfigBuilder) ID(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) CreationTimestamp(value time.Time) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.creationTimestamp = value
	b.fieldSet_[2] = true
	return b
}

// InstallerRoleArn sets the value of the 'installer_role_arn' attribute to the given value.
func (b *OidcConfigBuilder) InstallerRoleArn(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.installerRoleArn = value
	b.fieldSet_[3] = true
	return b
}

// IssuerUrl sets the value of the 'issuer_url' attribute to the given value.
func (b *OidcConfigBuilder) IssuerUrl(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.issuerUrl = value
	b.fieldSet_[4] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) LastUpdateTimestamp(value time.Time) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[5] = true
	return b
}

// LastUsedTimestamp sets the value of the 'last_used_timestamp' attribute to the given value.
func (b *OidcConfigBuilder) LastUsedTimestamp(value time.Time) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.lastUsedTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *OidcConfigBuilder) Managed(value bool) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.managed = value
	b.fieldSet_[7] = true
	return b
}

// OrganizationId sets the value of the 'organization_id' attribute to the given value.
func (b *OidcConfigBuilder) OrganizationId(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.organizationId = value
	b.fieldSet_[8] = true
	return b
}

// Reusable sets the value of the 'reusable' attribute to the given value.
func (b *OidcConfigBuilder) Reusable(value bool) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.reusable = value
	b.fieldSet_[9] = true
	return b
}

// SecretArn sets the value of the 'secret_arn' attribute to the given value.
func (b *OidcConfigBuilder) SecretArn(value string) *OidcConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.secretArn = value
	b.fieldSet_[10] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OidcConfigBuilder) Copy(object *OidcConfig) *OidcConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
