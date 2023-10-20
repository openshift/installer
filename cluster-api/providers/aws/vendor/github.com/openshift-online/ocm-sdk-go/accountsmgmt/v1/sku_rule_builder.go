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

// SkuRuleBuilder contains the data and logic needed to build 'sku_rule' objects.
//
// Identifies sku rule
type SkuRuleBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	allowed int
	quotaId string
	sku     string
}

// NewSkuRule creates a new builder of 'sku_rule' objects.
func NewSkuRule() *SkuRuleBuilder {
	return &SkuRuleBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SkuRuleBuilder) Link(value bool) *SkuRuleBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SkuRuleBuilder) ID(value string) *SkuRuleBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SkuRuleBuilder) HREF(value string) *SkuRuleBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SkuRuleBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *SkuRuleBuilder) Allowed(value int) *SkuRuleBuilder {
	b.allowed = value
	b.bitmap_ |= 8
	return b
}

// QuotaId sets the value of the 'quota_id' attribute to the given value.
func (b *SkuRuleBuilder) QuotaId(value string) *SkuRuleBuilder {
	b.quotaId = value
	b.bitmap_ |= 16
	return b
}

// Sku sets the value of the 'sku' attribute to the given value.
func (b *SkuRuleBuilder) Sku(value string) *SkuRuleBuilder {
	b.sku = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SkuRuleBuilder) Copy(object *SkuRule) *SkuRuleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.allowed = object.allowed
	b.quotaId = object.quotaId
	b.sku = object.sku
	return b
}

// Build creates a 'sku_rule' object using the configuration stored in the builder.
func (b *SkuRuleBuilder) Build() (object *SkuRule, err error) {
	object = new(SkuRule)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.allowed = b.allowed
	object.quotaId = b.quotaId
	object.sku = b.sku
	return
}
