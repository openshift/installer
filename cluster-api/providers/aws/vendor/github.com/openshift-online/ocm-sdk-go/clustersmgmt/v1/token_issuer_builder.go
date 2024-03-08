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

// TokenIssuerBuilder contains the data and logic needed to build 'token_issuer' objects.
//
// Representation of a token issuer used in an external authentication.
type TokenIssuerBuilder struct {
	bitmap_   uint32
	ca        string
	url       string
	audiences []string
}

// NewTokenIssuer creates a new builder of 'token_issuer' objects.
func NewTokenIssuer() *TokenIssuerBuilder {
	return &TokenIssuerBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TokenIssuerBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CA sets the value of the 'CA' attribute to the given value.
func (b *TokenIssuerBuilder) CA(value string) *TokenIssuerBuilder {
	b.ca = value
	b.bitmap_ |= 1
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *TokenIssuerBuilder) URL(value string) *TokenIssuerBuilder {
	b.url = value
	b.bitmap_ |= 2
	return b
}

// Audiences sets the value of the 'audiences' attribute to the given values.
func (b *TokenIssuerBuilder) Audiences(values ...string) *TokenIssuerBuilder {
	b.audiences = make([]string, len(values))
	copy(b.audiences, values)
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TokenIssuerBuilder) Copy(object *TokenIssuer) *TokenIssuerBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.ca = object.ca
	b.url = object.url
	if object.audiences != nil {
		b.audiences = make([]string, len(object.audiences))
		copy(b.audiences, object.audiences)
	} else {
		b.audiences = nil
	}
	return b
}

// Build creates a 'token_issuer' object using the configuration stored in the builder.
func (b *TokenIssuerBuilder) Build() (object *TokenIssuer, err error) {
	object = new(TokenIssuer)
	object.bitmap_ = b.bitmap_
	object.ca = b.ca
	object.url = b.url
	if b.audiences != nil {
		object.audiences = make([]string, len(b.audiences))
		copy(object.audiences, b.audiences)
	}
	return
}
