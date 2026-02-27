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

// AWSSTSRoleListBuilder contains the data and logic needed to build
// 'AWSSTS_role' objects.
type AWSSTSRoleListBuilder struct {
	items []*AWSSTSRoleBuilder
}

// NewAWSSTSRoleList creates a new builder of 'AWSSTS_role' objects.
func NewAWSSTSRoleList() *AWSSTSRoleListBuilder {
	return new(AWSSTSRoleListBuilder)
}

// Items sets the items of the list.
func (b *AWSSTSRoleListBuilder) Items(values ...*AWSSTSRoleBuilder) *AWSSTSRoleListBuilder {
	b.items = make([]*AWSSTSRoleBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *AWSSTSRoleListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *AWSSTSRoleListBuilder) Copy(list *AWSSTSRoleList) *AWSSTSRoleListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*AWSSTSRoleBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewAWSSTSRole().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'AWSSTS_role' objects using the
// configuration stored in the builder.
func (b *AWSSTSRoleListBuilder) Build() (list *AWSSTSRoleList, err error) {
	items := make([]*AWSSTSRole, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(AWSSTSRoleList)
	list.items = items
	return
}
