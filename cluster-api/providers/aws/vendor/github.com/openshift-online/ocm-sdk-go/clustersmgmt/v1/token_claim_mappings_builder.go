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

// TokenClaimMappingsBuilder contains the data and logic needed to build 'token_claim_mappings' objects.
//
// The claim mappings defined for users and groups.
type TokenClaimMappingsBuilder struct {
	bitmap_  uint32
	groups   *GroupsClaimBuilder
	userName *UsernameClaimBuilder
}

// NewTokenClaimMappings creates a new builder of 'token_claim_mappings' objects.
func NewTokenClaimMappings() *TokenClaimMappingsBuilder {
	return &TokenClaimMappingsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TokenClaimMappingsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Groups sets the value of the 'groups' attribute to the given value.
func (b *TokenClaimMappingsBuilder) Groups(value *GroupsClaimBuilder) *TokenClaimMappingsBuilder {
	b.groups = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// UserName sets the value of the 'user_name' attribute to the given value.
//
// The username claim mapping.
func (b *TokenClaimMappingsBuilder) UserName(value *UsernameClaimBuilder) *TokenClaimMappingsBuilder {
	b.userName = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TokenClaimMappingsBuilder) Copy(object *TokenClaimMappings) *TokenClaimMappingsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.groups != nil {
		b.groups = NewGroupsClaim().Copy(object.groups)
	} else {
		b.groups = nil
	}
	if object.userName != nil {
		b.userName = NewUsernameClaim().Copy(object.userName)
	} else {
		b.userName = nil
	}
	return b
}

// Build creates a 'token_claim_mappings' object using the configuration stored in the builder.
func (b *TokenClaimMappingsBuilder) Build() (object *TokenClaimMappings, err error) {
	object = new(TokenClaimMappings)
	object.bitmap_ = b.bitmap_
	if b.groups != nil {
		object.groups, err = b.groups.Build()
		if err != nil {
			return
		}
	}
	if b.userName != nil {
		object.userName, err = b.userName.Build()
		if err != nil {
			return
		}
	}
	return
}
