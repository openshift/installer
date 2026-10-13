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

// Details for `google` identity providers.
type GoogleIdentityProviderBuilder struct {
	fieldSet_    []bool
	clientID     string
	clientSecret string
	hostedDomain string
}

// NewGoogleIdentityProvider creates a new builder of 'google_identity_provider' objects.
func NewGoogleIdentityProvider() *GoogleIdentityProviderBuilder {
	return &GoogleIdentityProviderBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GoogleIdentityProviderBuilder) Empty() bool {
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

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) ClientID(value string) *GoogleIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.clientID = value
	b.fieldSet_[0] = true
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) ClientSecret(value string) *GoogleIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.clientSecret = value
	b.fieldSet_[1] = true
	return b
}

// HostedDomain sets the value of the 'hosted_domain' attribute to the given value.
func (b *GoogleIdentityProviderBuilder) HostedDomain(value string) *GoogleIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.hostedDomain = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GoogleIdentityProviderBuilder) Copy(object *GoogleIdentityProvider) *GoogleIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	b.hostedDomain = object.hostedDomain
	return b
}

// Build creates a 'google_identity_provider' object using the configuration stored in the builder.
func (b *GoogleIdentityProviderBuilder) Build() (object *GoogleIdentityProvider, err error) {
	object = new(GoogleIdentityProvider)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	object.hostedDomain = b.hostedDomain
	return
}
