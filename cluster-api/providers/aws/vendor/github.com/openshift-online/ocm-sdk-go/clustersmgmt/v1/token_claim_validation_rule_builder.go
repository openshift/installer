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

// TokenClaimValidationRuleBuilder contains the data and logic needed to build 'token_claim_validation_rule' objects.
//
// The rule that is applied to validate token claims to authenticate users.
type TokenClaimValidationRuleBuilder struct {
	bitmap_       uint32
	claim         string
	requiredValue string
}

// NewTokenClaimValidationRule creates a new builder of 'token_claim_validation_rule' objects.
func NewTokenClaimValidationRule() *TokenClaimValidationRuleBuilder {
	return &TokenClaimValidationRuleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TokenClaimValidationRuleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Claim sets the value of the 'claim' attribute to the given value.
func (b *TokenClaimValidationRuleBuilder) Claim(value string) *TokenClaimValidationRuleBuilder {
	b.claim = value
	b.bitmap_ |= 1
	return b
}

// RequiredValue sets the value of the 'required_value' attribute to the given value.
func (b *TokenClaimValidationRuleBuilder) RequiredValue(value string) *TokenClaimValidationRuleBuilder {
	b.requiredValue = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TokenClaimValidationRuleBuilder) Copy(object *TokenClaimValidationRule) *TokenClaimValidationRuleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.claim = object.claim
	b.requiredValue = object.requiredValue
	return b
}

// Build creates a 'token_claim_validation_rule' object using the configuration stored in the builder.
func (b *TokenClaimValidationRuleBuilder) Build() (object *TokenClaimValidationRule, err error) {
	object = new(TokenClaimValidationRule)
	object.bitmap_ = b.bitmap_
	object.claim = b.claim
	object.requiredValue = b.requiredValue
	return
}
