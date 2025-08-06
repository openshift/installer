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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// Representation of an add-on installation billing.
type AddonInstallationBillingBuilder struct {
	fieldSet_                 []bool
	billingMarketplaceAccount string
	billingModel              BillingModel
	href                      string
	id                        string
	kind                      string
}

// NewAddonInstallationBilling creates a new builder of 'addon_installation_billing' objects.
func NewAddonInstallationBilling() *AddonInstallationBillingBuilder {
	return &AddonInstallationBillingBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationBillingBuilder) Empty() bool {
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

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *AddonInstallationBillingBuilder) BillingMarketplaceAccount(value string) *AddonInstallationBillingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.billingMarketplaceAccount = value
	b.fieldSet_[0] = true
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Representation of an billing model field.
func (b *AddonInstallationBillingBuilder) BillingModel(value BillingModel) *AddonInstallationBillingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.billingModel = value
	b.fieldSet_[1] = true
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Href(value string) *AddonInstallationBillingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Id(value string) *AddonInstallationBillingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[3] = true
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Kind(value string) *AddonInstallationBillingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.kind = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationBillingBuilder) Copy(object *AddonInstallationBilling) *AddonInstallationBillingBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	b.billingModel = object.billingModel
	b.href = object.href
	b.id = object.id
	b.kind = object.kind
	return b
}

// Build creates a 'addon_installation_billing' object using the configuration stored in the builder.
func (b *AddonInstallationBillingBuilder) Build() (object *AddonInstallationBilling, err error) {
	object = new(AddonInstallationBilling)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.billingModel = b.billingModel
	object.href = b.href
	object.id = b.id
	object.kind = b.kind
	return
}
