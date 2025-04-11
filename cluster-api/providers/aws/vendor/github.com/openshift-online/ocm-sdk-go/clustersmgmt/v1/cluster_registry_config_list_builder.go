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

// ClusterRegistryConfigListBuilder contains the data and logic needed to build
// 'cluster_registry_config' objects.
type ClusterRegistryConfigListBuilder struct {
	items []*ClusterRegistryConfigBuilder
}

// NewClusterRegistryConfigList creates a new builder of 'cluster_registry_config' objects.
func NewClusterRegistryConfigList() *ClusterRegistryConfigListBuilder {
	return new(ClusterRegistryConfigListBuilder)
}

// Items sets the items of the list.
func (b *ClusterRegistryConfigListBuilder) Items(values ...*ClusterRegistryConfigBuilder) *ClusterRegistryConfigListBuilder {
	b.items = make([]*ClusterRegistryConfigBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *ClusterRegistryConfigListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *ClusterRegistryConfigListBuilder) Copy(list *ClusterRegistryConfigList) *ClusterRegistryConfigListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*ClusterRegistryConfigBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewClusterRegistryConfig().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'cluster_registry_config' objects using the
// configuration stored in the builder.
func (b *ClusterRegistryConfigListBuilder) Build() (list *ClusterRegistryConfigList, err error) {
	items := make([]*ClusterRegistryConfig, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(ClusterRegistryConfigList)
	list.items = items
	return
}
