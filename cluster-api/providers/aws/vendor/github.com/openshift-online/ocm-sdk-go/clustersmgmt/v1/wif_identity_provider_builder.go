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

// WifIdentityProviderBuilder contains the data and logic needed to build 'wif_identity_provider' objects.
type WifIdentityProviderBuilder struct {
	bitmap_            uint32
	allowedAudiences   []string
	identityProviderId string
	issuerUrl          string
	jwks               string
}

// NewWifIdentityProvider creates a new builder of 'wif_identity_provider' objects.
func NewWifIdentityProvider() *WifIdentityProviderBuilder {
	return &WifIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AllowedAudiences sets the value of the 'allowed_audiences' attribute to the given values.
func (b *WifIdentityProviderBuilder) AllowedAudiences(values ...string) *WifIdentityProviderBuilder {
	b.allowedAudiences = make([]string, len(values))
	copy(b.allowedAudiences, values)
	b.bitmap_ |= 1
	return b
}

// IdentityProviderId sets the value of the 'identity_provider_id' attribute to the given value.
func (b *WifIdentityProviderBuilder) IdentityProviderId(value string) *WifIdentityProviderBuilder {
	b.identityProviderId = value
	b.bitmap_ |= 2
	return b
}

// IssuerUrl sets the value of the 'issuer_url' attribute to the given value.
func (b *WifIdentityProviderBuilder) IssuerUrl(value string) *WifIdentityProviderBuilder {
	b.issuerUrl = value
	b.bitmap_ |= 4
	return b
}

// Jwks sets the value of the 'jwks' attribute to the given value.
func (b *WifIdentityProviderBuilder) Jwks(value string) *WifIdentityProviderBuilder {
	b.jwks = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifIdentityProviderBuilder) Copy(object *WifIdentityProvider) *WifIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	if b.allowedAudiences != nil {
		object.allowedAudiences = make([]string, len(b.allowedAudiences))
		copy(object.allowedAudiences, b.allowedAudiences)
	}
	object.identityProviderId = b.identityProviderId
	object.issuerUrl = b.issuerUrl
	object.jwks = b.jwks
	return
}
