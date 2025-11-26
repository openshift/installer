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

type ResourceRangeBuilder struct {
	fieldSet_ []bool
	max       int
	min       int
}

// NewResourceRange creates a new builder of 'resource_range' objects.
func NewResourceRange() *ResourceRangeBuilder {
	return &ResourceRangeBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceRangeBuilder) Empty() bool {
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

// Max sets the value of the 'max' attribute to the given value.
func (b *ResourceRangeBuilder) Max(value int) *ResourceRangeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.max = value
	b.fieldSet_[0] = true
	return b
}

// Min sets the value of the 'min' attribute to the given value.
func (b *ResourceRangeBuilder) Min(value int) *ResourceRangeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.min = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceRangeBuilder) Copy(object *ResourceRange) *ResourceRangeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.max = object.max
	b.min = object.min
	return b
}

// Build creates a 'resource_range' object using the configuration stored in the builder.
func (b *ResourceRangeBuilder) Build() (object *ResourceRange, err error) {
	object = new(ResourceRange)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.max = b.max
	object.min = b.min
	return
}
