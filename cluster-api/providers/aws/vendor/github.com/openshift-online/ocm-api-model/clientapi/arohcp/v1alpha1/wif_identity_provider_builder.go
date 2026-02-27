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

type WifIdentityProviderBuilder struct {
	fieldSet_          []bool
	allowedAudiences   []string
	identityProviderId string
	issuerUrl          string
	jwks               string
}

// NewWifIdentityProvider creates a new builder of 'wif_identity_provider' objects.
func NewWifIdentityProvider() *WifIdentityProviderBuilder {
	return &WifIdentityProviderBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifIdentityProviderBuilder) Empty() bool {
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

// AllowedAudiences sets the value of the 'allowed_audiences' attribute to the given values.
func (b *WifIdentityProviderBuilder) AllowedAudiences(values ...string) *WifIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.allowedAudiences = make([]string, len(values))
	copy(b.allowedAudiences, values)
	b.fieldSet_[0] = true
	return b
}

// IdentityProviderId sets the value of the 'identity_provider_id' attribute to the given value.
func (b *WifIdentityProviderBuilder) IdentityProviderId(value string) *WifIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.identityProviderId = value
	b.fieldSet_[1] = true
	return b
}

// IssuerUrl sets the value of the 'issuer_url' attribute to the given value.
func (b *WifIdentityProviderBuilder) IssuerUrl(value string) *WifIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.issuerUrl = value
	b.fieldSet_[2] = true
	return b
}

// Jwks sets the value of the 'jwks' attribute to the given value.
func (b *WifIdentityProviderBuilder) Jwks(value string) *WifIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.jwks = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifIdentityProviderBuilder) Copy(object *WifIdentityProvider) *WifIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.allowedAudiences != nil {
		b.allowedAudiences = make([]string, len(object.allowedAudiences))
		copy(b.allowedAudiences, object.allowedAudiences)
	} else {
		b.allowedAudiences = nil
	}
	b.identityProviderId = object.identityProviderId
	b.issuerUrl = object.issuerUrl
	b.jwks = object.jwks
	return b
}

// Build creates a 'wif_identity_provider' object using the configuration stored in the builder.
func (b *WifIdentityProviderBuilder) Build() (object *WifIdentityProvider, err error) {
	object = new(WifIdentityProvider)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.allowedAudiences != nil {
		object.allowedAudiences = make([]string, len(b.allowedAudiences))
		copy(object.allowedAudiences, b.allowedAudiences)
	}
	object.identityProviderId = b.identityProviderId
	object.issuerUrl = b.issuerUrl
	object.jwks = b.jwks
	return
}
