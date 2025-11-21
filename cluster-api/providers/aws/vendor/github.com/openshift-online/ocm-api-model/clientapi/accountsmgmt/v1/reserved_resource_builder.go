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

import (
	time "time"
)

type ReservedResourceBuilder struct {
	fieldSet_                 []bool
	availabilityZoneType      string
	billingMarketplaceAccount string
	billingModel              BillingModel
	count                     int
	createdAt                 time.Time
	resourceName              string
	resourceType              string
	scope                     string
	updatedAt                 time.Time
	byoc                      bool
}

// NewReservedResource creates a new builder of 'reserved_resource' objects.
func NewReservedResource() *ReservedResourceBuilder {
	return &ReservedResourceBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReservedResourceBuilder) Empty() bool {
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

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ReservedResourceBuilder) BYOC(value bool) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.byoc = value
	b.fieldSet_[0] = true
	return b
}

// AvailabilityZoneType sets the value of the 'availability_zone_type' attribute to the given value.
func (b *ReservedResourceBuilder) AvailabilityZoneType(value string) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.availabilityZoneType = value
	b.fieldSet_[1] = true
	return b
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *ReservedResourceBuilder) BillingMarketplaceAccount(value string) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.billingMarketplaceAccount = value
	b.fieldSet_[2] = true
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Billing model for subscripiton and reserved_resource resources.
func (b *ReservedResourceBuilder) BillingModel(value BillingModel) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.billingModel = value
	b.fieldSet_[3] = true
	return b
}

// Count sets the value of the 'count' attribute to the given value.
func (b *ReservedResourceBuilder) Count(value int) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.count = value
	b.fieldSet_[4] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ReservedResourceBuilder) CreatedAt(value time.Time) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[5] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *ReservedResourceBuilder) ResourceName(value string) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.resourceName = value
	b.fieldSet_[6] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *ReservedResourceBuilder) ResourceType(value string) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.resourceType = value
	b.fieldSet_[7] = true
	return b
}

// Scope sets the value of the 'scope' attribute to the given value.
func (b *ReservedResourceBuilder) Scope(value string) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.scope = value
	b.fieldSet_[8] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ReservedResourceBuilder) UpdatedAt(value time.Time) *ReservedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReservedResourceBuilder) Copy(object *ReservedResource) *ReservedResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.byoc = object.byoc
	b.availabilityZoneType = object.availabilityZoneType
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	b.billingModel = object.billingModel
	b.count = object.count
	b.createdAt = object.createdAt
	b.resourceName = object.resourceName
	b.resourceType = object.resourceType
	b.scope = object.scope
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'reserved_resource' object using the configuration stored in the builder.
func (b *ReservedResourceBuilder) Build() (object *ReservedResource, err error) {
	object = new(ReservedResource)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.byoc = b.byoc
	object.availabilityZoneType = b.availabilityZoneType
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.billingModel = b.billingModel
	object.count = b.count
	object.createdAt = b.createdAt
	object.resourceName = b.resourceName
	object.resourceType = b.resourceType
	object.scope = b.scope
	object.updatedAt = b.updatedAt
	return
}
