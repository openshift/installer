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

// ClusterRegistrationListBuilder contains the data and logic needed to build
// 'cluster_registration' objects.
type ClusterRegistrationListBuilder struct {
	items []*ClusterRegistrationBuilder
}

// NewClusterRegistrationList creates a new builder of 'cluster_registration' objects.
func NewClusterRegistrationList() *ClusterRegistrationListBuilder {
	return new(ClusterRegistrationListBuilder)
}

// Items sets the items of the list.
func (b *ClusterRegistrationListBuilder) Items(values ...*ClusterRegistrationBuilder) *ClusterRegistrationListBuilder {
	b.items = make([]*ClusterRegistrationBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *ClusterRegistrationListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *ClusterRegistrationListBuilder) Copy(list *ClusterRegistrationList) *ClusterRegistrationListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*ClusterRegistrationBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewClusterRegistration().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'cluster_registration' objects using the
// configuration stored in the builder.
func (b *ClusterRegistrationListBuilder) Build() (list *ClusterRegistrationList, err error) {
	items := make([]*ClusterRegistration, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(ClusterRegistrationList)
	list.items = items
	return
}
