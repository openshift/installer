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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of an external authentication provider.
type ExternalAuthBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	claim     *ExternalAuthClaimBuilder
	clients   []*ExternalAuthClientConfigBuilder
	issuer    *TokenIssuerBuilder
	status    *ExternalAuthStatusBuilder
}

// NewExternalAuth creates a new builder of 'external_auth' objects.
func NewExternalAuth() *ExternalAuthBuilder {
	return &ExternalAuthBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ExternalAuthBuilder) Link(value bool) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ExternalAuthBuilder) ID(value string) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ExternalAuthBuilder) HREF(value string) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// Claim sets the value of the 'claim' attribute to the given value.
//
// The claims and validation rules used in the configuration of the external authentication.
func (b *ExternalAuthBuilder) Claim(value *ExternalAuthClaimBuilder) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.claim = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Clients sets the value of the 'clients' attribute to the given values.
func (b *ExternalAuthBuilder) Clients(values ...*ExternalAuthClientConfigBuilder) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.clients = make([]*ExternalAuthClientConfigBuilder, len(values))
	copy(b.clients, values)
	b.fieldSet_[4] = true
	return b
}

// Issuer sets the value of the 'issuer' attribute to the given value.
//
// Representation of a token issuer used in an external authentication.
func (b *ExternalAuthBuilder) Issuer(value *TokenIssuerBuilder) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.issuer = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of the status of an external authentication provider.
func (b *ExternalAuthBuilder) Status(value *ExternalAuthStatusBuilder) *ExternalAuthBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthBuilder) Copy(object *ExternalAuth) *ExternalAuthBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if object.status != nil {
		b.status = NewExternalAuthStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	return b
}

// Build creates a 'external_auth' object using the configuration stored in the builder.
func (b *ExternalAuthBuilder) Build() (object *ExternalAuth, err error) {
	object = new(ExternalAuth)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	return
}
