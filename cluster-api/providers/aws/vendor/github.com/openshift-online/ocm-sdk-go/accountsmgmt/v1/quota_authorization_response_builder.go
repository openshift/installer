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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// QuotaAuthorizationResponseBuilder contains the data and logic needed to build 'quota_authorization_response' objects.
type QuotaAuthorizationResponseBuilder struct {
	bitmap_         uint32
	excessResources []*ReservedResourceBuilder
	subscription    *SubscriptionBuilder
	allowed         bool
}

// NewQuotaAuthorizationResponse creates a new builder of 'quota_authorization_response' objects.
func NewQuotaAuthorizationResponse() *QuotaAuthorizationResponseBuilder {
	return &QuotaAuthorizationResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *QuotaAuthorizationResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *QuotaAuthorizationResponseBuilder) Allowed(value bool) *QuotaAuthorizationResponseBuilder {
	b.allowed = value
	b.bitmap_ |= 1
	return b
}

// ExcessResources sets the value of the 'excess_resources' attribute to the given values.
func (b *QuotaAuthorizationResponseBuilder) ExcessResources(values ...*ReservedResourceBuilder) *QuotaAuthorizationResponseBuilder {
	b.excessResources = make([]*ReservedResourceBuilder, len(values))
	copy(b.excessResources, values)
	b.bitmap_ |= 2
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
func (b *QuotaAuthorizationResponseBuilder) Subscription(value *SubscriptionBuilder) *QuotaAuthorizationResponseBuilder {
	b.subscription = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *QuotaAuthorizationResponseBuilder) Copy(object *QuotaAuthorizationResponse) *QuotaAuthorizationResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.allowed = object.allowed
	if object.excessResources != nil {
		b.excessResources = make([]*ReservedResourceBuilder, len(object.excessResources))
		for i, v := range object.excessResources {
			b.excessResources[i] = NewReservedResource().Copy(v)
		}
	} else {
		b.excessResources = nil
	}
	if object.subscription != nil {
		b.subscription = NewSubscription().Copy(object.subscription)
	} else {
		b.subscription = nil
	}
	return b
}

// Build creates a 'quota_authorization_response' object using the configuration stored in the builder.
func (b *QuotaAuthorizationResponseBuilder) Build() (object *QuotaAuthorizationResponse, err error) {
	object = new(QuotaAuthorizationResponse)
	object.bitmap_ = b.bitmap_
	object.allowed = b.allowed
	if b.excessResources != nil {
		object.excessResources = make([]*ReservedResource, len(b.excessResources))
		for i, v := range b.excessResources {
			object.excessResources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.subscription != nil {
		object.subscription, err = b.subscription.Build()
		if err != nil {
			return
		}
	}
	return
}
