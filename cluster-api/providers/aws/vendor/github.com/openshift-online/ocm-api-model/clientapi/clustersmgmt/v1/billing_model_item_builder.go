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

// BillingModelItem represents a billing model
type BillingModelItemBuilder struct {
	fieldSet_        []bool
	id               string
	href             string
	billingModelType string
	description      string
	displayName      string
	marketplace      string
}

// NewBillingModelItem creates a new builder of 'billing_model_item' objects.
func NewBillingModelItem() *BillingModelItemBuilder {
	return &BillingModelItemBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *BillingModelItemBuilder) Link(value bool) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *BillingModelItemBuilder) ID(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *BillingModelItemBuilder) HREF(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *BillingModelItemBuilder) Empty() bool {
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

// BillingModelType sets the value of the 'billing_model_type' attribute to the given value.
func (b *BillingModelItemBuilder) BillingModelType(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.billingModelType = value
	b.fieldSet_[3] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *BillingModelItemBuilder) Description(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.description = value
	b.fieldSet_[4] = true
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *BillingModelItemBuilder) DisplayName(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.displayName = value
	b.fieldSet_[5] = true
	return b
}

// Marketplace sets the value of the 'marketplace' attribute to the given value.
func (b *BillingModelItemBuilder) Marketplace(value string) *BillingModelItemBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.marketplace = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *BillingModelItemBuilder) Copy(object *BillingModelItem) *BillingModelItemBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.billingModelType = object.billingModelType
	b.description = object.description
	b.displayName = object.displayName
	b.marketplace = object.marketplace
	return b
}

// Build creates a 'billing_model_item' object using the configuration stored in the builder.
func (b *BillingModelItemBuilder) Build() (object *BillingModelItem, err error) {
	object = new(BillingModelItem)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.billingModelType = b.billingModelType
	object.description = b.description
	object.displayName = b.displayName
	object.marketplace = b.marketplace
	return
}
