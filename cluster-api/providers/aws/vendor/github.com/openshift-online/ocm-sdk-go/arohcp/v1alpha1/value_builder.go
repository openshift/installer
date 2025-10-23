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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// ValueBuilder contains the data and logic needed to build 'value' objects.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
type ValueBuilder struct {
	bitmap_ uint32
	unit    string
	value   float64
}

// NewValue creates a new builder of 'value' objects.
func NewValue() *ValueBuilder {
	return &ValueBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ValueBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Unit sets the value of the 'unit' attribute to the given value.
func (b *ValueBuilder) Unit(value string) *ValueBuilder {
	b.unit = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *ValueBuilder) Value(value float64) *ValueBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ValueBuilder) Copy(object *Value) *ValueBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.unit = object.unit
	b.value = object.value
	return b
}

// Build creates a 'value' object using the configuration stored in the builder.
func (b *ValueBuilder) Build() (object *Value, err error) {
	object = new(Value)
	object.bitmap_ = b.bitmap_
	object.unit = b.unit
	object.value = b.value
	return
}
