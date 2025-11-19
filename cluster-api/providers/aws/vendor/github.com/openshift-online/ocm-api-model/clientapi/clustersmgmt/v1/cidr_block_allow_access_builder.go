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

type CIDRBlockAllowAccessBuilder struct {
	fieldSet_ []bool
	mode      string
	values    []string
}

// NewCIDRBlockAllowAccess creates a new builder of 'CIDR_block_allow_access' objects.
func NewCIDRBlockAllowAccess() *CIDRBlockAllowAccessBuilder {
	return &CIDRBlockAllowAccessBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CIDRBlockAllowAccessBuilder) Empty() bool {
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

// Mode sets the value of the 'mode' attribute to the given value.
func (b *CIDRBlockAllowAccessBuilder) Mode(value string) *CIDRBlockAllowAccessBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.mode = value
	b.fieldSet_[0] = true
	return b
}

// Values sets the value of the 'values' attribute to the given values.
func (b *CIDRBlockAllowAccessBuilder) Values(values ...string) *CIDRBlockAllowAccessBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.values = make([]string, len(values))
	copy(b.values, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CIDRBlockAllowAccessBuilder) Copy(object *CIDRBlockAllowAccess) *CIDRBlockAllowAccessBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.mode = object.mode
	if object.values != nil {
		b.values = make([]string, len(object.values))
		copy(b.values, object.values)
	} else {
		b.values = nil
	}
	return b
}

// Build creates a 'CIDR_block_allow_access' object using the configuration stored in the builder.
func (b *CIDRBlockAllowAccessBuilder) Build() (object *CIDRBlockAllowAccess, err error) {
	object = new(CIDRBlockAllowAccess)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.mode = b.mode
	if b.values != nil {
		object.values = make([]string, len(b.values))
		copy(object.values, b.values)
	}
	return
}
