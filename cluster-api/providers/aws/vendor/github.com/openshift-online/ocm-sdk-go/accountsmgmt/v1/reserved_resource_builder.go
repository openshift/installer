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

import (
	time "time"
)

// ReservedResourceBuilder contains the data and logic needed to build 'reserved_resource' objects.
type ReservedResourceBuilder struct {
	bitmap_                   uint32
	availabilityZoneType      string
	billingMarketplaceAccount string
	billingModel              BillingModel
	count                     int
	createdAt                 time.Time
	resourceName              string
	resourceType              string
	updatedAt                 time.Time
	byoc                      bool
}

// NewReservedResource creates a new builder of 'reserved_resource' objects.
func NewReservedResource() *ReservedResourceBuilder {
	return &ReservedResourceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReservedResourceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ReservedResourceBuilder) BYOC(value bool) *ReservedResourceBuilder {
	b.byoc = value
	b.bitmap_ |= 1
	return b
}

// AvailabilityZoneType sets the value of the 'availability_zone_type' attribute to the given value.
func (b *ReservedResourceBuilder) AvailabilityZoneType(value string) *ReservedResourceBuilder {
	b.availabilityZoneType = value
	b.bitmap_ |= 2
	return b
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *ReservedResourceBuilder) BillingMarketplaceAccount(value string) *ReservedResourceBuilder {
	b.billingMarketplaceAccount = value
	b.bitmap_ |= 4
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Billing model for subscripiton and reserved_resource resources.
func (b *ReservedResourceBuilder) BillingModel(value BillingModel) *ReservedResourceBuilder {
	b.billingModel = value
	b.bitmap_ |= 8
	return b
}

// Count sets the value of the 'count' attribute to the given value.
func (b *ReservedResourceBuilder) Count(value int) *ReservedResourceBuilder {
	b.count = value
	b.bitmap_ |= 16
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ReservedResourceBuilder) CreatedAt(value time.Time) *ReservedResourceBuilder {
	b.createdAt = value
	b.bitmap_ |= 32
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *ReservedResourceBuilder) ResourceName(value string) *ReservedResourceBuilder {
	b.resourceName = value
	b.bitmap_ |= 64
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *ReservedResourceBuilder) ResourceType(value string) *ReservedResourceBuilder {
	b.resourceType = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ReservedResourceBuilder) UpdatedAt(value time.Time) *ReservedResourceBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReservedResourceBuilder) Copy(object *ReservedResource) *ReservedResourceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.byoc = object.byoc
	b.availabilityZoneType = object.availabilityZoneType
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	b.billingModel = object.billingModel
	b.count = object.count
	b.createdAt = object.createdAt
	b.resourceName = object.resourceName
	b.resourceType = object.resourceType
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'reserved_resource' object using the configuration stored in the builder.
func (b *ReservedResourceBuilder) Build() (object *ReservedResource, err error) {
	object = new(ReservedResource)
	object.bitmap_ = b.bitmap_
	object.byoc = b.byoc
	object.availabilityZoneType = b.availabilityZoneType
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.billingModel = b.billingModel
	object.count = b.count
	object.createdAt = b.createdAt
	object.resourceName = b.resourceName
	object.resourceType = b.resourceType
	object.updatedAt = b.updatedAt
	return
}
