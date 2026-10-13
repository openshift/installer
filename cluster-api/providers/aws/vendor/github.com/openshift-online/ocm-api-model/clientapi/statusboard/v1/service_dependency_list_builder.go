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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

// ServiceDependencyListBuilder contains the data and logic needed to build
// 'service_dependency' objects.
type ServiceDependencyListBuilder struct {
	items []*ServiceDependencyBuilder
}

// NewServiceDependencyList creates a new builder of 'service_dependency' objects.
func NewServiceDependencyList() *ServiceDependencyListBuilder {
	return new(ServiceDependencyListBuilder)
}

// Items sets the items of the list.
func (b *ServiceDependencyListBuilder) Items(values ...*ServiceDependencyBuilder) *ServiceDependencyListBuilder {
	b.items = make([]*ServiceDependencyBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *ServiceDependencyListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *ServiceDependencyListBuilder) Copy(list *ServiceDependencyList) *ServiceDependencyListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*ServiceDependencyBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewServiceDependency().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'service_dependency' objects using the
// configuration stored in the builder.
func (b *ServiceDependencyListBuilder) Build() (list *ServiceDependencyList, err error) {
	items := make([]*ServiceDependency, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(ServiceDependencyList)
	list.items = items
	return
}
