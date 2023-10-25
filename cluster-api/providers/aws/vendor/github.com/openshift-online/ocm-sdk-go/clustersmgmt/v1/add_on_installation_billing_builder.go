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

// AddOnInstallationBillingBuilder contains the data and logic needed to build 'add_on_installation_billing' objects.
//
// Representation of an add-on installation billing.
type AddOnInstallationBillingBuilder struct {
	bitmap_                   uint32
	id                        string
	href                      string
	billingMarketplaceAccount string
	billingModel              BillingModel
}

// NewAddOnInstallationBilling creates a new builder of 'add_on_installation_billing' objects.
func NewAddOnInstallationBilling() *AddOnInstallationBillingBuilder {
	return &AddOnInstallationBillingBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddOnInstallationBillingBuilder) Link(value bool) *AddOnInstallationBillingBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddOnInstallationBillingBuilder) ID(value string) *AddOnInstallationBillingBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddOnInstallationBillingBuilder) HREF(value string) *AddOnInstallationBillingBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnInstallationBillingBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BillingMarketplaceAccount sets the value of the 'billing_marketplace_account' attribute to the given value.
func (b *AddOnInstallationBillingBuilder) BillingMarketplaceAccount(value string) *AddOnInstallationBillingBuilder {
	b.billingMarketplaceAccount = value
	b.bitmap_ |= 8
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Billing model for cluster resources.
func (b *AddOnInstallationBillingBuilder) BillingModel(value BillingModel) *AddOnInstallationBillingBuilder {
	b.billingModel = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnInstallationBillingBuilder) Copy(object *AddOnInstallationBilling) *AddOnInstallationBillingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.billingMarketplaceAccount = object.billingMarketplaceAccount
	b.billingModel = object.billingModel
	return b
}

// Build creates a 'add_on_installation_billing' object using the configuration stored in the builder.
func (b *AddOnInstallationBillingBuilder) Build() (object *AddOnInstallationBilling, err error) {
	object = new(AddOnInstallationBilling)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.billingMarketplaceAccount = b.billingMarketplaceAccount
	object.billingModel = b.billingModel
	return
}
