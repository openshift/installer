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

// BillingModelItemBuilder contains the data and logic needed to build 'billing_model_item' objects.
//
// BillingModelItem represents a billing model
type BillingModelItemBuilder struct {
	bitmap_          uint32
	id               string
	href             string
	billingModelType string
	description      string
	displayName      string
	marketplace      string
}

// NewBillingModelItem creates a new builder of 'billing_model_item' objects.
func NewBillingModelItem() *BillingModelItemBuilder {
	return &BillingModelItemBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *BillingModelItemBuilder) Link(value bool) *BillingModelItemBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *BillingModelItemBuilder) ID(value string) *BillingModelItemBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *BillingModelItemBuilder) HREF(value string) *BillingModelItemBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *BillingModelItemBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BillingModelType sets the value of the 'billing_model_type' attribute to the given value.
func (b *BillingModelItemBuilder) BillingModelType(value string) *BillingModelItemBuilder {
	b.billingModelType = value
	b.bitmap_ |= 8
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *BillingModelItemBuilder) Description(value string) *BillingModelItemBuilder {
	b.description = value
	b.bitmap_ |= 16
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *BillingModelItemBuilder) DisplayName(value string) *BillingModelItemBuilder {
	b.displayName = value
	b.bitmap_ |= 32
	return b
}

// Marketplace sets the value of the 'marketplace' attribute to the given value.
func (b *BillingModelItemBuilder) Marketplace(value string) *BillingModelItemBuilder {
	b.marketplace = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *BillingModelItemBuilder) Copy(object *BillingModelItem) *BillingModelItemBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.billingModelType = b.billingModelType
	object.description = b.description
	object.displayName = b.displayName
	object.marketplace = b.marketplace
	return
}
