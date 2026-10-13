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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// AzureServiceManagedIdentityListBuilder contains the data and logic needed to build
// 'azure_service_managed_identity' objects.
type AzureServiceManagedIdentityListBuilder struct {
	items []*AzureServiceManagedIdentityBuilder
}

// NewAzureServiceManagedIdentityList creates a new builder of 'azure_service_managed_identity' objects.
func NewAzureServiceManagedIdentityList() *AzureServiceManagedIdentityListBuilder {
	return new(AzureServiceManagedIdentityListBuilder)
}

// Items sets the items of the list.
func (b *AzureServiceManagedIdentityListBuilder) Items(values ...*AzureServiceManagedIdentityBuilder) *AzureServiceManagedIdentityListBuilder {
	b.items = make([]*AzureServiceManagedIdentityBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *AzureServiceManagedIdentityListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *AzureServiceManagedIdentityListBuilder) Copy(list *AzureServiceManagedIdentityList) *AzureServiceManagedIdentityListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*AzureServiceManagedIdentityBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewAzureServiceManagedIdentity().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'azure_service_managed_identity' objects using the
// configuration stored in the builder.
func (b *AzureServiceManagedIdentityListBuilder) Build() (list *AzureServiceManagedIdentityList, err error) {
	items := make([]*AzureServiceManagedIdentity, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(AzureServiceManagedIdentityList)
	list.items = items
	return
}
