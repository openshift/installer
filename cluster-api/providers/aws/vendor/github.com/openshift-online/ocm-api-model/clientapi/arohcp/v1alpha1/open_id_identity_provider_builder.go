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

// Details for `openid` identity providers.
type OpenIDIdentityProviderBuilder struct {
	fieldSet_                []bool
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
	return &OpenIDIdentityProviderBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OpenIDIdentityProviderBuilder) Empty() bool {
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

// CA sets the value of the 'CA' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) CA(value string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.ca = value
	b.fieldSet_[0] = true
	return b
}

// Claims sets the value of the 'claims' attribute to the given value.
//
// _OpenID_ identity provider claims.
func (b *OpenIDIdentityProviderBuilder) Claims(value *OpenIDClaimsBuilder) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.claims = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ClientID(value string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.clientID = value
	b.fieldSet_[2] = true
	return b
}

// ClientSecret sets the value of the 'client_secret' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ClientSecret(value string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.clientSecret = value
	b.fieldSet_[3] = true
	return b
}

// ExtraAuthorizeParameters sets the value of the 'extra_authorize_parameters' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) ExtraAuthorizeParameters(value map[string]string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.extraAuthorizeParameters = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// ExtraScopes sets the value of the 'extra_scopes' attribute to the given values.
func (b *OpenIDIdentityProviderBuilder) ExtraScopes(values ...string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.extraScopes = make([]string, len(values))
	copy(b.extraScopes, values)
	b.fieldSet_[5] = true
	return b
}

// Issuer sets the value of the 'issuer' attribute to the given value.
func (b *OpenIDIdentityProviderBuilder) Issuer(value string) *OpenIDIdentityProviderBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.issuer = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OpenIDIdentityProviderBuilder) Copy(object *OpenIDIdentityProvider) *OpenIDIdentityProviderBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
