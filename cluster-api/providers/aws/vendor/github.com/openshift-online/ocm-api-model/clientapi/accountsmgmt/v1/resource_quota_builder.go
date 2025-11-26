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

type ResourceQuotaBuilder struct {
	fieldSet_      []bool
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
	return &ResourceQuotaBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ResourceQuotaBuilder) Link(value bool) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ResourceQuotaBuilder) ID(value string) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ResourceQuotaBuilder) HREF(value string) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceQuotaBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// SKU sets the value of the 'SKU' attribute to the given value.
func (b *ResourceQuotaBuilder) SKU(value string) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.sku = value
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ResourceQuotaBuilder) CreatedAt(value time.Time) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *ResourceQuotaBuilder) OrganizationID(value string) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.organizationID = value
	b.fieldSet_[5] = true
	return b
}

// SkuCount sets the value of the 'sku_count' attribute to the given value.
func (b *ResourceQuotaBuilder) SkuCount(value int) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.skuCount = value
	b.fieldSet_[6] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *ResourceQuotaBuilder) Type(value string) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.type_ = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ResourceQuotaBuilder) UpdatedAt(value time.Time) *ResourceQuotaBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceQuotaBuilder) Copy(object *ResourceQuota) *ResourceQuotaBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.sku = b.sku
	object.createdAt = b.createdAt
	object.organizationID = b.organizationID
	object.skuCount = b.skuCount
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	return
}
