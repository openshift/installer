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

// ResourceQuotaBuilder contains the data and logic needed to build 'resource_quota' objects.
type ResourceQuotaBuilder struct {
	bitmap_        uint32
	id             string
	href           string
	sku            string
	createdAt      time.Time
	organizationID string
	skuCount       int
	type_          string
	updatedAt      time.Time
}

// NewResourceQuota creates a new builder of 'resource_quota' objects.
func NewResourceQuota() *ResourceQuotaBuilder {
	return &ResourceQuotaBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ResourceQuotaBuilder) Link(value bool) *ResourceQuotaBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ResourceQuotaBuilder) ID(value string) *ResourceQuotaBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ResourceQuotaBuilder) HREF(value string) *ResourceQuotaBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceQuotaBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// SKU sets the value of the 'SKU' attribute to the given value.
func (b *ResourceQuotaBuilder) SKU(value string) *ResourceQuotaBuilder {
	b.sku = value
	b.bitmap_ |= 8
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ResourceQuotaBuilder) CreatedAt(value time.Time) *ResourceQuotaBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *ResourceQuotaBuilder) OrganizationID(value string) *ResourceQuotaBuilder {
	b.organizationID = value
	b.bitmap_ |= 32
	return b
}

// SkuCount sets the value of the 'sku_count' attribute to the given value.
func (b *ResourceQuotaBuilder) SkuCount(value int) *ResourceQuotaBuilder {
	b.skuCount = value
	b.bitmap_ |= 64
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *ResourceQuotaBuilder) Type(value string) *ResourceQuotaBuilder {
	b.type_ = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ResourceQuotaBuilder) UpdatedAt(value time.Time) *ResourceQuotaBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceQuotaBuilder) Copy(object *ResourceQuota) *ResourceQuotaBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.sku = object.sku
	b.createdAt = object.createdAt
	b.organizationID = object.organizationID
	b.skuCount = object.skuCount
	b.type_ = object.type_
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'resource_quota' object using the configuration stored in the builder.
func (b *ResourceQuotaBuilder) Build() (object *ResourceQuota, err error) {
	object = new(ResourceQuota)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.sku = b.sku
	object.createdAt = b.createdAt
	object.organizationID = b.organizationID
	object.skuCount = b.skuCount
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	return
}
