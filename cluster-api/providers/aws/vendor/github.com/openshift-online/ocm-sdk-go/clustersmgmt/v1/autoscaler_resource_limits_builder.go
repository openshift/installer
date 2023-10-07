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

// AutoscalerResourceLimitsBuilder contains the data and logic needed to build 'autoscaler_resource_limits' objects.
type AutoscalerResourceLimitsBuilder struct {
	bitmap_       uint32
	cores         *ResourceRangeBuilder
	maxNodesTotal int
	memory        *ResourceRangeBuilder
}

// NewAutoscalerResourceLimits creates a new builder of 'autoscaler_resource_limits' objects.
func NewAutoscalerResourceLimits() *AutoscalerResourceLimitsBuilder {
	return &AutoscalerResourceLimitsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AutoscalerResourceLimitsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Cores sets the value of the 'cores' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) Cores(value *ResourceRangeBuilder) *AutoscalerResourceLimitsBuilder {
	b.cores = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// MaxNodesTotal sets the value of the 'max_nodes_total' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) MaxNodesTotal(value int) *AutoscalerResourceLimitsBuilder {
	b.maxNodesTotal = value
	b.bitmap_ |= 2
	return b
}

// Memory sets the value of the 'memory' attribute to the given value.
func (b *AutoscalerResourceLimitsBuilder) Memory(value *ResourceRangeBuilder) *AutoscalerResourceLimitsBuilder {
	b.memory = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AutoscalerResourceLimitsBuilder) Copy(object *AutoscalerResourceLimits) *AutoscalerResourceLimitsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
