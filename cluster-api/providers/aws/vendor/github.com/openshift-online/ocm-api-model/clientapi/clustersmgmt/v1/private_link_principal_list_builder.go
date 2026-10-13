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

// PrivateLinkPrincipalListBuilder contains the data and logic needed to build
// 'private_link_principal' objects.
type PrivateLinkPrincipalListBuilder struct {
	items []*PrivateLinkPrincipalBuilder
}

// NewPrivateLinkPrincipalList creates a new builder of 'private_link_principal' objects.
func NewPrivateLinkPrincipalList() *PrivateLinkPrincipalListBuilder {
	return new(PrivateLinkPrincipalListBuilder)
}

// Items sets the items of the list.
func (b *PrivateLinkPrincipalListBuilder) Items(values ...*PrivateLinkPrincipalBuilder) *PrivateLinkPrincipalListBuilder {
	b.items = make([]*PrivateLinkPrincipalBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *PrivateLinkPrincipalListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *PrivateLinkPrincipalListBuilder) Copy(list *PrivateLinkPrincipalList) *PrivateLinkPrincipalListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*PrivateLinkPrincipalBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewPrivateLinkPrincipal().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'private_link_principal' objects using the
// configuration stored in the builder.
func (b *PrivateLinkPrincipalListBuilder) Build() (list *PrivateLinkPrincipalList, err error) {
	items := make([]*PrivateLinkPrincipal, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(PrivateLinkPrincipalList)
	list.items = items
	return
}
