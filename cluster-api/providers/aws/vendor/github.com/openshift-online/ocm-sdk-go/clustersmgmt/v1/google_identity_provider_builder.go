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

// GoogleIdentityProviderBuilder contains the data and logic needed to build 'google_identity_provider' objects.
//
// Details for `google` identity providers.
type GoogleIdentityProviderBuilder struct {
	bitmap_      uint32
	clientID     string
	clientSecret string
	hostedDomain string
}

// NewGoogleIdentityProvider creates a new builder of 'google_identity_provider' objects.
func NewGoogleIdentityProvider() *GoogleIdentityProviderBuilder {
	return &GoogleIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GoogleIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) ClientID(value string) *GoogleIdentityProviderBuilder {
	b.clientID = value
	b.bitmap_ |= 1
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) ClientSecret(value string) *GoogleIdentityProviderBuilder {
	b.clientSecret = value
	b.bitmap_ |= 2
	return b
}

// HostedDomain sets the value of the 'hosted_domain' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) HostedDomain(value string) *GoogleIdentityProviderBuilder {
	b.hostedDomain = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GoogleIdentityProviderBuilder) Copy(object *GoogleIdentityProvider) *GoogleIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	b.hostedDomain = object.hostedDomain
	return b
}

// Build creates a 'google_identity_provider' object using the configuration stored in the builder.
func (b *GoogleIdentityProviderBuilder) Build() (object *GoogleIdentityProvider, err error) {
	object = new(GoogleIdentityProvider)
	object.bitmap_ = b.bitmap_
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	object.hostedDomain = b.hostedDomain
	return
}
