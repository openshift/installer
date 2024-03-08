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

// ExternalAuthBuilder contains the data and logic needed to build 'external_auth' objects.
//
// Representation of an external authentication provider.
type ExternalAuthBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	claim   *ExternalAuthClaimBuilder
	clients []*ExternalAuthClientConfigBuilder
	issuer  *TokenIssuerBuilder
}

// NewExternalAuth creates a new builder of 'external_auth' objects.
func NewExternalAuth() *ExternalAuthBuilder {
	return &ExternalAuthBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ExternalAuthBuilder) Link(value bool) *ExternalAuthBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ExternalAuthBuilder) ID(value string) *ExternalAuthBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ExternalAuthBuilder) HREF(value string) *ExternalAuthBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Claim sets the value of the 'claim' attribute to the given value.
//
// The claims and validation rules used in the configuration of the external authentication.
func (b *ExternalAuthBuilder) Claim(value *ExternalAuthClaimBuilder) *ExternalAuthBuilder {
	b.claim = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Clients sets the value of the 'clients' attribute to the given values.
func (b *ExternalAuthBuilder) Clients(values ...*ExternalAuthClientConfigBuilder) *ExternalAuthBuilder {
	b.clients = make([]*ExternalAuthClientConfigBuilder, len(values))
	copy(b.clients, values)
	b.bitmap_ |= 16
	return b
}

// Issuer sets the value of the 'issuer' attribute to the given value.
//
// Representation of a token issuer used in an external authentication.
func (b *ExternalAuthBuilder) Issuer(value *TokenIssuerBuilder) *ExternalAuthBuilder {
	b.issuer = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthBuilder) Copy(object *ExternalAuth) *ExternalAuthBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.claim != nil {
		b.claim = NewExternalAuthClaim().Copy(object.claim)
	} else {
		b.claim = nil
	}
	if object.clients != nil {
		b.clients = make([]*ExternalAuthClientConfigBuilder, len(object.clients))
		for i, v := range object.clients {
			b.clients[i] = NewExternalAuthClientConfig().Copy(v)
		}
	} else {
		b.clients = nil
	}
	if object.issuer != nil {
		b.issuer = NewTokenIssuer().Copy(object.issuer)
	} else {
		b.issuer = nil
	}
	return b
}

// Build creates a 'external_auth' object using the configuration stored in the builder.
func (b *ExternalAuthBuilder) Build() (object *ExternalAuth, err error) {
	object = new(ExternalAuth)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.claim != nil {
		object.claim, err = b.claim.Build()
		if err != nil {
			return
		}
	}
	if b.clients != nil {
		object.clients = make([]*ExternalAuthClientConfig, len(b.clients))
		for i, v := range b.clients {
			object.clients[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.issuer != nil {
		object.issuer, err = b.issuer.Build()
		if err != nil {
			return
		}
	}
	return
}
