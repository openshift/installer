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

// AddOnInstallationBillingListBuilder contains the data and logic needed to build
// 'add_on_installation_billing' objects.
type AddOnInstallationBillingListBuilder struct {
	items []*AddOnInstallationBillingBuilder
}

// NewAddOnInstallationBillingList creates a new builder of 'add_on_installation_billing' objects.
func NewAddOnInstallationBillingList() *AddOnInstallationBillingListBuilder {
	return new(AddOnInstallationBillingListBuilder)
}

// Items sets the items of the list.
func (b *AddOnInstallationBillingListBuilder) Items(values ...*AddOnInstallationBillingBuilder) *AddOnInstallationBillingListBuilder {
	b.items = make([]*AddOnInstallationBillingBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *AddOnInstallationBillingListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *AddOnInstallationBillingListBuilder) Copy(list *AddOnInstallationBillingList) *AddOnInstallationBillingListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*AddOnInstallationBillingBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewAddOnInstallationBilling().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'add_on_installation_billing' objects using the
// configuration stored in the builder.
func (b *AddOnInstallationBillingListBuilder) Build() (list *AddOnInstallationBillingList, err error) {
	items := make([]*AddOnInstallationBilling, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(AddOnInstallationBillingList)
	list.items = items
	return
}
