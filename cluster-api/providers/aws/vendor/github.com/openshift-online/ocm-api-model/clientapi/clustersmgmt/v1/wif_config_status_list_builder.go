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

// WifConfigStatusListBuilder contains the data and logic needed to build
// 'wif_config_status' objects.
type WifConfigStatusListBuilder struct {
	items []*WifConfigStatusBuilder
}

// NewWifConfigStatusList creates a new builder of 'wif_config_status' objects.
func NewWifConfigStatusList() *WifConfigStatusListBuilder {
	return new(WifConfigStatusListBuilder)
}

// Items sets the items of the list.
func (b *WifConfigStatusListBuilder) Items(values ...*WifConfigStatusBuilder) *WifConfigStatusListBuilder {
	b.items = make([]*WifConfigStatusBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *WifConfigStatusListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *WifConfigStatusListBuilder) Copy(list *WifConfigStatusList) *WifConfigStatusListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*WifConfigStatusBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewWifConfigStatus().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'wif_config_status' objects using the
// configuration stored in the builder.
func (b *WifConfigStatusListBuilder) Build() (list *WifConfigStatusList, err error) {
	items := make([]*WifConfigStatus, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(WifConfigStatusList)
	list.items = items
	return
}
