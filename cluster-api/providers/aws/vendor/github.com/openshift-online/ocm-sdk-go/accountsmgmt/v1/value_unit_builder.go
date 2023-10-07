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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// ValueUnitBuilder contains the data and logic needed to build 'value_unit' objects.
type ValueUnitBuilder struct {
	bitmap_ uint32
	unit    string
	value   float64
}

// NewValueUnit creates a new builder of 'value_unit' objects.
func NewValueUnit() *ValueUnitBuilder {
	return &ValueUnitBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ValueUnitBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Unit sets the value of the 'unit' attribute to the given value.
func (b *ValueUnitBuilder) Unit(value string) *ValueUnitBuilder {
	b.unit = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *ValueUnitBuilder) Value(value float64) *ValueUnitBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ValueUnitBuilder) Copy(object *ValueUnit) *ValueUnitBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.unit = object.unit
	b.value = object.value
	return b
}

// Build creates a 'value_unit' object using the configuration stored in the builder.
func (b *ValueUnitBuilder) Build() (object *ValueUnit, err error) {
	object = new(ValueUnit)
	object.bitmap_ = b.bitmap_
	object.unit = b.unit
	object.value = b.value
	return
}
