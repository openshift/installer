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

// Identifies sku rule
type SkuRuleBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	allowed   int
	quotaId   string
	sku       string
}

// NewSkuRule creates a new builder of 'sku_rule' objects.
func NewSkuRule() *SkuRuleBuilder {
	return &SkuRuleBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *SkuRuleBuilder) Link(value bool) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *SkuRuleBuilder) ID(value string) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *SkuRuleBuilder) HREF(value string) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SkuRuleBuilder) Empty() bool {
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

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *SkuRuleBuilder) Allowed(value int) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.allowed = value
	b.fieldSet_[3] = true
	return b
}

// QuotaId sets the value of the 'quota_id' attribute to the given value.
func (b *SkuRuleBuilder) QuotaId(value string) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.quotaId = value
	b.fieldSet_[4] = true
	return b
}

// Sku sets the value of the 'sku' attribute to the given value.
func (b *SkuRuleBuilder) Sku(value string) *SkuRuleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.sku = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SkuRuleBuilder) Copy(object *SkuRule) *SkuRuleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.allowed = b.allowed
	object.quotaId = b.quotaId
	object.sku = b.sku
	return
}
