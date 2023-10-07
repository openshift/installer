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

// QuotaRulesBuilder contains the data and logic needed to build 'quota_rules' objects.
type QuotaRulesBuilder struct {
	bitmap_          uint32
	availabilityZone string
	billingModel     string
	byoc             string
	cloud            string
	cost             int
	name             string
	product          string
	quotaId          string
}

// NewQuotaRules creates a new builder of 'quota_rules' objects.
func NewQuotaRules() *QuotaRulesBuilder {
	return &QuotaRulesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *QuotaRulesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *QuotaRulesBuilder) AvailabilityZone(value string) *QuotaRulesBuilder {
	b.availabilityZone = value
	b.bitmap_ |= 1
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
func (b *QuotaRulesBuilder) BillingModel(value string) *QuotaRulesBuilder {
	b.billingModel = value
	b.bitmap_ |= 2
	return b
}

// Byoc sets the value of the 'byoc' attribute to the given value.
func (b *QuotaRulesBuilder) Byoc(value string) *QuotaRulesBuilder {
	b.byoc = value
	b.bitmap_ |= 4
	return b
}

// Cloud sets the value of the 'cloud' attribute to the given value.
func (b *QuotaRulesBuilder) Cloud(value string) *QuotaRulesBuilder {
	b.cloud = value
	b.bitmap_ |= 8
	return b
}

// Cost sets the value of the 'cost' attribute to the given value.
func (b *QuotaRulesBuilder) Cost(value int) *QuotaRulesBuilder {
	b.cost = value
	b.bitmap_ |= 16
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *QuotaRulesBuilder) Name(value string) *QuotaRulesBuilder {
	b.name = value
	b.bitmap_ |= 32
	return b
}

// Product sets the value of the 'product' attribute to the given value.
func (b *QuotaRulesBuilder) Product(value string) *QuotaRulesBuilder {
	b.product = value
	b.bitmap_ |= 64
	return b
}

// QuotaId sets the value of the 'quota_id' attribute to the given value.
func (b *QuotaRulesBuilder) QuotaId(value string) *QuotaRulesBuilder {
	b.quotaId = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *QuotaRulesBuilder) Copy(object *QuotaRules) *QuotaRulesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.availabilityZone = object.availabilityZone
	b.billingModel = object.billingModel
	b.byoc = object.byoc
	b.cloud = object.cloud
	b.cost = object.cost
	b.name = object.name
	b.product = object.product
	b.quotaId = object.quotaId
	return b
}

// Build creates a 'quota_rules' object using the configuration stored in the builder.
func (b *QuotaRulesBuilder) Build() (object *QuotaRules, err error) {
	object = new(QuotaRules)
	object.bitmap_ = b.bitmap_
	object.availabilityZone = b.availabilityZone
	object.billingModel = b.billingModel
	object.byoc = b.byoc
	object.cloud = b.cloud
	object.cost = b.cost
	object.name = b.name
	object.product = b.product
	object.quotaId = b.quotaId
	return
}
