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

// StorageQuotaListBuilder contains the data and logic needed to build
// 'storage_quota' objects.
type StorageQuotaListBuilder struct {
	items []*StorageQuotaBuilder
}

// NewStorageQuotaList creates a new builder of 'storage_quota' objects.
func NewStorageQuotaList() *StorageQuotaListBuilder {
	return new(StorageQuotaListBuilder)
}

// Items sets the items of the list.
func (b *StorageQuotaListBuilder) Items(values ...*StorageQuotaBuilder) *StorageQuotaListBuilder {
	b.items = make([]*StorageQuotaBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *StorageQuotaListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *StorageQuotaListBuilder) Copy(list *StorageQuotaList) *StorageQuotaListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*StorageQuotaBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewStorageQuota().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'storage_quota' objects using the
// configuration stored in the builder.
func (b *StorageQuotaListBuilder) Build() (list *StorageQuotaList, err error) {
	items := make([]*StorageQuota, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(StorageQuotaList)
	list.items = items
	return
}
