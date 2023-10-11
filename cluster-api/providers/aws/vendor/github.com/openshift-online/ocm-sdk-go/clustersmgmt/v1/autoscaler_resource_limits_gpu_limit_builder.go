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

// AutoscalerResourceLimitsGPULimitBuilder contains the data and logic needed to build 'autoscaler_resource_limits_GPU_limit' objects.
type AutoscalerResourceLimitsGPULimitBuilder struct {
	bitmap_ uint32
	range_  *ResourceRangeBuilder
	type_   string
}

// NewAutoscalerResourceLimitsGPULimit creates a new builder of 'autoscaler_resource_limits_GPU_limit' objects.
func NewAutoscalerResourceLimitsGPULimit() *AutoscalerResourceLimitsGPULimitBuilder {
	return &AutoscalerResourceLimitsGPULimitBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Range sets the value of the 'range' attribute to the given value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Range(value *ResourceRangeBuilder) *AutoscalerResourceLimitsGPULimitBuilder {
	b.range_ = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Type(value string) *AutoscalerResourceLimitsGPULimitBuilder {
	b.type_ = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Copy(object *AutoscalerResourceLimitsGPULimit) *AutoscalerResourceLimitsGPULimitBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.range_ != nil {
		b.range_ = NewResourceRange().Copy(object.range_)
	} else {
		b.range_ = nil
	}
	b.type_ = object.type_
	return b
}

// Build creates a 'autoscaler_resource_limits_GPU_limit' object using the configuration stored in the builder.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Build() (object *AutoscalerResourceLimitsGPULimit, err error) {
	object = new(AutoscalerResourceLimitsGPULimit)
	object.bitmap_ = b.bitmap_
	if b.range_ != nil {
		object.range_, err = b.range_.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	return
}
