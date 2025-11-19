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

type AutoscalerResourceLimitsBuilder struct {
	fieldSet_     []bool
	gpus          []*AutoscalerResourceLimitsGPULimitBuilder
	cores         *ResourceRangeBuilder
	maxNodesTotal int
	memory        *ResourceRangeBuilder
}

// NewAutoscalerResourceLimits creates a new builder of 'autoscaler_resource_limits' objects.
func NewAutoscalerResourceLimits() *AutoscalerResourceLimitsBuilder {
	return &AutoscalerResourceLimitsBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AutoscalerResourceLimitsBuilder) Empty() bool {
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

// GPUS sets the value of the 'GPUS' attribute to the given values.
func (b *AutoscalerResourceLimitsBuilder) GPUS(values ...*AutoscalerResourceLimitsGPULimitBuilder) *AutoscalerResourceLimitsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.gpus = make([]*AutoscalerResourceLimitsGPULimitBuilder, len(values))
	copy(b.gpus, values)
	b.fieldSet_[0] = true
	return b
}

// Cores sets the value of the 'cores' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) Cores(value *ResourceRangeBuilder) *AutoscalerResourceLimitsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.cores = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// MaxNodesTotal sets the value of the 'max_nodes_total' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) MaxNodesTotal(value int) *AutoscalerResourceLimitsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.maxNodesTotal = value
	b.fieldSet_[2] = true
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) Memory(value *ResourceRangeBuilder) *AutoscalerResourceLimitsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.memory = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AutoscalerResourceLimitsBuilder) Copy(object *AutoscalerResourceLimits) *AutoscalerResourceLimitsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.gpus != nil {
		b.gpus = make([]*AutoscalerResourceLimitsGPULimitBuilder, len(object.gpus))
		for i, v := range object.gpus {
			b.gpus[i] = NewAutoscalerResourceLimitsGPULimit().Copy(v)
		}
	} else {
		b.gpus = nil
	}
	if object.cores != nil {
		b.cores = NewResourceRange().Copy(object.cores)
	} else {
		b.cores = nil
	}
	b.maxNodesTotal = object.maxNodesTotal
	if object.memory != nil {
		b.memory = NewResourceRange().Copy(object.memory)
	} else {
		b.memory = nil
	}
	return b
}

// Build creates a 'autoscaler_resource_limits' object using the configuration stored in the builder.
func (b *AutoscalerResourceLimitsBuilder) Build() (object *AutoscalerResourceLimits, err error) {
	object = new(AutoscalerResourceLimits)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.gpus != nil {
		object.gpus = make([]*AutoscalerResourceLimitsGPULimit, len(b.gpus))
		for i, v := range b.gpus {
			object.gpus[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.cores != nil {
		object.cores, err = b.cores.Build()
		if err != nil {
			return
		}
	}
	object.maxNodesTotal = b.maxNodesTotal
	if b.memory != nil {
		object.memory, err = b.memory.Build()
		if err != nil {
			return
		}
	}
	return
}
