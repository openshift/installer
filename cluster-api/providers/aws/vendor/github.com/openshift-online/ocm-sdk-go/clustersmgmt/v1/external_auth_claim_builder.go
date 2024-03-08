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

// ExternalAuthClaimBuilder contains the data and logic needed to build 'external_auth_claim' objects.
//
// The claims and validation rules used in the configuration of the external authentication.
type ExternalAuthClaimBuilder struct {
	bitmap_         uint32
	mappings        *TokenClaimMappingsBuilder
	validationRules []*TokenClaimValidationRuleBuilder
}

// NewExternalAuthClaim creates a new builder of 'external_auth_claim' objects.
func NewExternalAuthClaim() *ExternalAuthClaimBuilder {
	return &ExternalAuthClaimBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ExternalAuthClaimBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Mappings sets the value of the 'mappings' attribute to the given value.
//
// The claim mappings defined for users and groups.
func (b *ExternalAuthClaimBuilder) Mappings(value *TokenClaimMappingsBuilder) *ExternalAuthClaimBuilder {
	b.mappings = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// ValidationRules sets the value of the 'validation_rules' attribute to the given values.
func (b *ExternalAuthClaimBuilder) ValidationRules(values ...*TokenClaimValidationRuleBuilder) *ExternalAuthClaimBuilder {
	b.validationRules = make([]*TokenClaimValidationRuleBuilder, len(values))
	copy(b.validationRules, values)
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ExternalAuthClaimBuilder) Copy(object *ExternalAuthClaim) *ExternalAuthClaimBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.mappings != nil {
		b.mappings = NewTokenClaimMappings().Copy(object.mappings)
	} else {
		b.mappings = nil
	}
	if object.validationRules != nil {
		b.validationRules = make([]*TokenClaimValidationRuleBuilder, len(object.validationRules))
		for i, v := range object.validationRules {
			b.validationRules[i] = NewTokenClaimValidationRule().Copy(v)
		}
	} else {
		b.validationRules = nil
	}
	return b
}

// Build creates a 'external_auth_claim' object using the configuration stored in the builder.
func (b *ExternalAuthClaimBuilder) Build() (object *ExternalAuthClaim, err error) {
	object = new(ExternalAuthClaim)
	object.bitmap_ = b.bitmap_
	if b.mappings != nil {
		object.mappings, err = b.mappings.Build()
		if err != nil {
			return
		}
	}
	if b.validationRules != nil {
		object.validationRules = make([]*TokenClaimValidationRule, len(b.validationRules))
		for i, v := range b.validationRules {
			object.validationRules[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
