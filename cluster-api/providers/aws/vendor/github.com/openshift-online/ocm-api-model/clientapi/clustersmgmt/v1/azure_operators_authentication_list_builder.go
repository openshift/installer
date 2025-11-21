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

// AzureOperatorsAuthenticationListBuilder contains the data and logic needed to build
// 'azure_operators_authentication' objects.
type AzureOperatorsAuthenticationListBuilder struct {
	items []*AzureOperatorsAuthenticationBuilder
}

// NewAzureOperatorsAuthenticationList creates a new builder of 'azure_operators_authentication' objects.
func NewAzureOperatorsAuthenticationList() *AzureOperatorsAuthenticationListBuilder {
	return new(AzureOperatorsAuthenticationListBuilder)
}

// Items sets the items of the list.
func (b *AzureOperatorsAuthenticationListBuilder) Items(values ...*AzureOperatorsAuthenticationBuilder) *AzureOperatorsAuthenticationListBuilder {
	b.items = make([]*AzureOperatorsAuthenticationBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *AzureOperatorsAuthenticationListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *AzureOperatorsAuthenticationListBuilder) Copy(list *AzureOperatorsAuthenticationList) *AzureOperatorsAuthenticationListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*AzureOperatorsAuthenticationBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewAzureOperatorsAuthentication().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'azure_operators_authentication' objects using the
// configuration stored in the builder.
func (b *AzureOperatorsAuthenticationListBuilder) Build() (list *AzureOperatorsAuthenticationList, err error) {
	items := make([]*AzureOperatorsAuthentication, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(AzureOperatorsAuthenticationList)
	list.items = items
	return
}
