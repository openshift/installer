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

// Configuration status of a WifConfig.
type WifConfigStatusBuilder struct {
	fieldSet_   []bool
	description string
	configured  bool
}

// NewWifConfigStatus creates a new builder of 'wif_config_status' objects.
func NewWifConfigStatus() *WifConfigStatusBuilder {
	return &WifConfigStatusBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifConfigStatusBuilder) Empty() bool {
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

// Configured sets the value of the 'configured' attribute to the given value.
func (b *WifConfigStatusBuilder) Configured(value bool) *WifConfigStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.configured = value
	b.fieldSet_[0] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *WifConfigStatusBuilder) Description(value string) *WifConfigStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.description = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifConfigStatusBuilder) Copy(object *WifConfigStatus) *WifConfigStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.configured = object.configured
	b.description = object.description
	return b
}

// Build creates a 'wif_config_status' object using the configuration stored in the builder.
func (b *WifConfigStatusBuilder) Build() (object *WifConfigStatus, err error) {
	object = new(WifConfigStatus)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.configured = b.configured
	object.description = b.description
	return
}
