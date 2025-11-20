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

// ExternalConfigurationListBuilder contains the data and logic needed to build
// 'external_configuration' objects.
type ExternalConfigurationListBuilder struct {
	items []*ExternalConfigurationBuilder
}

// NewExternalConfigurationList creates a new builder of 'external_configuration' objects.
func NewExternalConfigurationList() *ExternalConfigurationListBuilder {
	return new(ExternalConfigurationListBuilder)
}

// Items sets the items of the list.
func (b *ExternalConfigurationListBuilder) Items(values ...*ExternalConfigurationBuilder) *ExternalConfigurationListBuilder {
	b.items = make([]*ExternalConfigurationBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *ExternalConfigurationListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *ExternalConfigurationListBuilder) Copy(list *ExternalConfigurationList) *ExternalConfigurationListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*ExternalConfigurationBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewExternalConfiguration().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'external_configuration' objects using the
// configuration stored in the builder.
func (b *ExternalConfigurationListBuilder) Build() (list *ExternalConfigurationList, err error) {
	items := make([]*ExternalConfiguration, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(ExternalConfigurationList)
	list.items = items
	return
}
