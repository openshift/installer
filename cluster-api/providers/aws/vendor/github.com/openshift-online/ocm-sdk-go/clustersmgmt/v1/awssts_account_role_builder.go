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

// AWSSTSAccountRoleBuilder contains the data and logic needed to build 'AWSSTS_account_role' objects.
//
// Representation of an sts account role for a rosa cluster
type AWSSTSAccountRoleBuilder struct {
	bitmap_ uint32
	items   []*AWSSTSRoleBuilder
	prefix  string
}

// NewAWSSTSAccountRole creates a new builder of 'AWSSTS_account_role' objects.
func NewAWSSTSAccountRole() *AWSSTSAccountRoleBuilder {
	return &AWSSTSAccountRoleBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSAccountRoleBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Items sets the value of the 'items' attribute to the given values.
func (b *AWSSTSAccountRoleBuilder) Items(values ...*AWSSTSRoleBuilder) *AWSSTSAccountRoleBuilder {
	b.items = make([]*AWSSTSRoleBuilder, len(values))
	copy(b.items, values)
	b.bitmap_ |= 1
	return b
}

// Prefix sets the value of the 'prefix' attribute to the given value.
func (b *AWSSTSAccountRoleBuilder) Prefix(value string) *AWSSTSAccountRoleBuilder {
	b.prefix = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSAccountRoleBuilder) Copy(object *AWSSTSAccountRole) *AWSSTSAccountRoleBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.items != nil {
		b.items = make([]*AWSSTSRoleBuilder, len(object.items))
		for i, v := range object.items {
			b.items[i] = NewAWSSTSRole().Copy(v)
		}
	} else {
		b.items = nil
	}
	b.prefix = object.prefix
	return b
}

// Build creates a 'AWSSTS_account_role' object using the configuration stored in the builder.
func (b *AWSSTSAccountRoleBuilder) Build() (object *AWSSTSAccountRole, err error) {
	object = new(AWSSTSAccountRole)
	object.bitmap_ = b.bitmap_
	if b.items != nil {
		object.items = make([]*AWSSTSRole, len(b.items))
		for i, v := range b.items {
			object.items[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.prefix = b.prefix
	return
}
