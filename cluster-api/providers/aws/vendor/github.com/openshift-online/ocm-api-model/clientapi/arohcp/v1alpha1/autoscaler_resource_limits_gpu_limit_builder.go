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

type AutoscalerResourceLimitsGPULimitBuilder struct {
	fieldSet_ []bool
	range_    *ResourceRangeBuilder
	type_     string
}

// NewAutoscalerResourceLimitsGPULimit creates a new builder of 'autoscaler_resource_limits_GPU_limit' objects.
func NewAutoscalerResourceLimitsGPULimit() *AutoscalerResourceLimitsGPULimitBuilder {
	return &AutoscalerResourceLimitsGPULimitBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Empty() bool {
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

// Range sets the value of the 'range' attribute to the given value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Range(value *ResourceRangeBuilder) *AutoscalerResourceLimitsGPULimitBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.range_ = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Type(value string) *AutoscalerResourceLimitsGPULimitBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.type_ = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AutoscalerResourceLimitsGPULimitBuilder) Copy(object *AutoscalerResourceLimitsGPULimit) *AutoscalerResourceLimitsGPULimitBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.range_ != nil {
		object.range_, err = b.range_.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	return
}
