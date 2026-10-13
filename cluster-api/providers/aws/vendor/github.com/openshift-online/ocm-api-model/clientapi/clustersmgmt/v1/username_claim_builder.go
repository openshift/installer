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

// The username claim mapping.
type UsernameClaimBuilder struct {
	fieldSet_    []bool
	claim        string
	prefix       string
	prefixPolicy string
}

// NewUsernameClaim creates a new builder of 'username_claim' objects.
func NewUsernameClaim() *UsernameClaimBuilder {
	return &UsernameClaimBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *UsernameClaimBuilder) Empty() bool {
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

// Claim sets the value of the 'claim' attribute to the given value.
func (b *UsernameClaimBuilder) Claim(value string) *UsernameClaimBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.claim = value
	b.fieldSet_[0] = true
	return b
}

// Prefix sets the value of the 'prefix' attribute to the given value.
func (b *UsernameClaimBuilder) Prefix(value string) *UsernameClaimBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.prefix = value
	b.fieldSet_[1] = true
	return b
}

// PrefixPolicy sets the value of the 'prefix_policy' attribute to the given value.
func (b *UsernameClaimBuilder) PrefixPolicy(value string) *UsernameClaimBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.prefixPolicy = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UsernameClaimBuilder) Copy(object *UsernameClaim) *UsernameClaimBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.claim = object.claim
	b.prefix = object.prefix
	b.prefixPolicy = object.prefixPolicy
	return b
}

// Build creates a 'username_claim' object using the configuration stored in the builder.
func (b *UsernameClaimBuilder) Build() (object *UsernameClaim, err error) {
	object = new(UsernameClaim)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.claim = b.claim
	object.prefix = b.prefix
	object.prefixPolicy = b.prefixPolicy
	return
}
