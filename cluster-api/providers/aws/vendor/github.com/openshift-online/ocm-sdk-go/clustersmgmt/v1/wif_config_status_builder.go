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

// WifConfigStatusBuilder contains the data and logic needed to build 'wif_config_status' objects.
//
// Configuration status of a WifConfig.
type WifConfigStatusBuilder struct {
	bitmap_     uint32
	description string
	configured  bool
}

// NewWifConfigStatus creates a new builder of 'wif_config_status' objects.
func NewWifConfigStatus() *WifConfigStatusBuilder {
	return &WifConfigStatusBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifConfigStatusBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Configured sets the value of the 'configured' attribute to the given value.
func (b *WifConfigStatusBuilder) Configured(value bool) *WifConfigStatusBuilder {
	b.configured = value
	b.bitmap_ |= 1
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *WifConfigStatusBuilder) Description(value string) *WifConfigStatusBuilder {
	b.description = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifConfigStatusBuilder) Copy(object *WifConfigStatus) *WifConfigStatusBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.configured = object.configured
	b.description = object.description
	return b
}

// Build creates a 'wif_config_status' object using the configuration stored in the builder.
func (b *WifConfigStatusBuilder) Build() (object *WifConfigStatus, err error) {
	object = new(WifConfigStatus)
	object.bitmap_ = b.bitmap_
	object.configured = b.configured
	object.description = b.description
	return
}
