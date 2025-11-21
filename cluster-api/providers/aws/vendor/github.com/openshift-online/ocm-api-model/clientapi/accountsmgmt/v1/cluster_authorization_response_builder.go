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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

type ClusterAuthorizationResponseBuilder struct {
	fieldSet_       []bool
	excessResources []*ReservedResourceBuilder
	subscription    *SubscriptionBuilder
	allowed         bool
}

// NewClusterAuthorizationResponse creates a new builder of 'cluster_authorization_response' objects.
func NewClusterAuthorizationResponse() *ClusterAuthorizationResponseBuilder {
	return &ClusterAuthorizationResponseBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAuthorizationResponseBuilder) Empty() bool {
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

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *ClusterAuthorizationResponseBuilder) Allowed(value bool) *ClusterAuthorizationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.allowed = value
	b.fieldSet_[0] = true
	return b
}

// ExcessResources sets the value of the 'excess_resources' attribute to the given values.
func (b *ClusterAuthorizationResponseBuilder) ExcessResources(values ...*ReservedResourceBuilder) *ClusterAuthorizationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.excessResources = make([]*ReservedResourceBuilder, len(values))
	copy(b.excessResources, values)
	b.fieldSet_[1] = true
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
func (b *ClusterAuthorizationResponseBuilder) Subscription(value *SubscriptionBuilder) *ClusterAuthorizationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.subscription = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAuthorizationResponseBuilder) Copy(object *ClusterAuthorizationResponse) *ClusterAuthorizationResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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

// Build creates a 'cluster_authorization_response' object using the configuration stored in the builder.
func (b *ClusterAuthorizationResponseBuilder) Build() (object *ClusterAuthorizationResponse, err error) {
	object = new(ClusterAuthorizationResponse)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
