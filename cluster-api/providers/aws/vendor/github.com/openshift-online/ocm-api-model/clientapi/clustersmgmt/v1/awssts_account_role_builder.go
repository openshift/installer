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

// Representation of an sts account role for a rosa cluster
type AWSSTSAccountRoleBuilder struct {
	fieldSet_ []bool
	items     []*AWSSTSRoleBuilder
	prefix    string
}

// NewAWSSTSAccountRole creates a new builder of 'AWSSTS_account_role' objects.
func NewAWSSTSAccountRole() *AWSSTSAccountRoleBuilder {
	return &AWSSTSAccountRoleBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSAccountRoleBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Items sets the value of the 'items' attribute to the given values.
func (b *AWSSTSAccountRoleBuilder) Items(values ...*AWSSTSRoleBuilder) *AWSSTSAccountRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.items = make([]*AWSSTSRoleBuilder, len(values))
	copy(b.items, values)
	b.fieldSet_[0] = true
	return b
}

// Prefix sets the value of the 'prefix' attribute to the given value.
func (b *AWSSTSAccountRoleBuilder) Prefix(value string) *AWSSTSAccountRoleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.prefix = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSAccountRoleBuilder) Copy(object *AWSSTSAccountRole) *AWSSTSAccountRoleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
