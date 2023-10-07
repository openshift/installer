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

// OpenIDIdentityProviderBuilder contains the data and logic needed to build 'open_ID_identity_provider' objects.
//
// Details for `openid` identity providers.
type OpenIDIdentityProviderBuilder struct {
	bitmap_                  uint32
	ca                       string
	claims                   *OpenIDClaimsBuilder
	clientID                 string
	clientSecret             string
	extraAuthorizeParameters map[string]string
	extraScopes              []string
	issuer                   string
}

// NewOpenIDIdentityProvider creates a new builder of 'open_ID_identity_provider' objects.
func NewOpenIDIdentityProvider() *OpenIDIdentityProviderBuilder {
	return &OpenIDIdentityProviderBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OpenIDIdentityProviderBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CA sets the value of the 'CA' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) CA(value string) *OpenIDIdentityProviderBuilder {
	b.ca = value
	b.bitmap_ |= 1
	return b
}

// Claims sets the value of the 'claims' attribute to the given value.
//
// _OpenID_ identity provider claims.
func (b *OpenIDIdentityProviderBuilder) Claims(value *OpenIDClaimsBuilder) *OpenIDIdentityProviderBuilder {
	b.claims = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ClientID(value string) *OpenIDIdentityProviderBuilder {
	b.clientID = value
	b.bitmap_ |= 4
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ClientSecret(value string) *OpenIDIdentityProviderBuilder {
	b.clientSecret = value
	b.bitmap_ |= 8
	return b
}

// ExtraAuthorizeParameters sets the value of the 'extra_authorize_parameters' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ExtraAuthorizeParameters(value map[string]string) *OpenIDIdentityProviderBuilder {
	b.extraAuthorizeParameters = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// ExtraScopes sets the value of the 'extra_scopes' attribute to the given values.
func (b *OpenIDIdentityProviderBuilder) ExtraScopes(values ...string) *OpenIDIdentityProviderBuilder {
	b.extraScopes = make([]string, len(values))
	copy(b.extraScopes, values)
	b.bitmap_ |= 32
	return b
}

// Issuer sets the value of the 'issuer' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) Issuer(value string) *OpenIDIdentityProviderBuilder {
	b.issuer = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OpenIDIdentityProviderBuilder) Copy(object *OpenIDIdentityProvider) *OpenIDIdentityProviderBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.ca = object.ca
	if object.claims != nil {
		b.claims = NewOpenIDClaims().Copy(object.claims)
	} else {
		b.claims = nil
	}
	b.clientID = object.clientID
	b.clientSecret = object.clientSecret
	if len(object.extraAuthorizeParameters) > 0 {
		b.extraAuthorizeParameters = map[string]string{}
		for k, v := range object.extraAuthorizeParameters {
			b.extraAuthorizeParameters[k] = v
		}
	} else {
		b.extraAuthorizeParameters = nil
	}
	if object.extraScopes != nil {
		b.extraScopes = make([]string, len(object.extraScopes))
		copy(b.extraScopes, object.extraScopes)
	} else {
		b.extraScopes = nil
	}
	b.issuer = object.issuer
	return b
}

// Build creates a 'open_ID_identity_provider' object using the configuration stored in the builder.
func (b *OpenIDIdentityProviderBuilder) Build() (object *OpenIDIdentityProvider, err error) {
	object = new(OpenIDIdentityProvider)
	object.bitmap_ = b.bitmap_
	object.ca = b.ca
	if b.claims != nil {
		object.claims, err = b.claims.Build()
		if err != nil {
			return
		}
	}
	object.clientID = b.clientID
	object.clientSecret = b.clientSecret
	if b.extraAuthorizeParameters != nil {
		object.extraAuthorizeParameters = make(map[string]string)
		for k, v := range b.extraAuthorizeParameters {
			object.extraAuthorizeParameters[k] = v
		}
	}
	if b.extraScopes != nil {
		object.extraScopes = make([]string, len(b.extraScopes))
		copy(object.extraScopes, b.extraScopes)
	}
	object.issuer = b.issuer
	return
}
