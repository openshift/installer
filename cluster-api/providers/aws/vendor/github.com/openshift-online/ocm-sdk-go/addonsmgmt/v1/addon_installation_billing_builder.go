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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonInstallationBillingBuilder contains the data and logic needed to build 'addon_installation_billing' objects.
//
// Representation of an add-on installation billing.
type AddonInstallationBillingBuilder struct {
	bitmap_                   uint32
	billingMarketplaceAccount string
	billingModel              BillingModel
	href                      string
	id                        string
	kind                      string
}

// NewAddonInstallationBilling creates a new builder of 'addon_installation_billing' objects.
func NewAddonInstallationBilling() *AddonInstallationBillingBuilder {
	return &AddonInstallationBillingBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonInstallationBillingBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *AddonInstallationBillingBuilder) BillingMarketplaceAccount(value string) *AddonInstallationBillingBuilder {
	b.billingMarketplaceAccount = value
	b.bitmap_ |= 1
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Representation of an billing model field.
func (b *AddonInstallationBillingBuilder) BillingModel(value BillingModel) *AddonInstallationBillingBuilder {
	b.billingModel = value
	b.bitmap_ |= 2
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Href(value string) *AddonInstallationBillingBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Id(value string) *AddonInstallationBillingBuilder {
	b.id = value
	b.bitmap_ |= 8
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *AddonInstallationBillingBuilder) Kind(value string) *AddonInstallationBillingBuilder {
	b.kind = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonInstallationBillingBuilder) Copy(object *AddonInstallationBilling) *AddonInstallationBillingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.billingModel = b.billingModel
	object.href = b.href
	object.id = b.id
	object.kind = b.kind
	return
}
