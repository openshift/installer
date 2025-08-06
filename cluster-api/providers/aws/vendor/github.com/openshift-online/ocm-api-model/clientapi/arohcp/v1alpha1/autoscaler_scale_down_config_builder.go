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

type AutoscalerScaleDownConfigBuilder struct {
	fieldSet_            []bool
	delayAfterAdd        string
	delayAfterDelete     string
	delayAfterFailure    string
	unneededTime         string
	utilizationThreshold string
	enabled              bool
}

// NewAutoscalerScaleDownConfig creates a new builder of 'autoscaler_scale_down_config' objects.
func NewAutoscalerScaleDownConfig() *AutoscalerScaleDownConfigBuilder {
	return &AutoscalerScaleDownConfigBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AutoscalerScaleDownConfigBuilder) Empty() bool {
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

// DelayAfterAdd sets the value of the 'delay_after_add' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) DelayAfterAdd(value string) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.delayAfterAdd = value
	b.fieldSet_[0] = true
	return b
}

// DelayAfterDelete sets the value of the 'delay_after_delete' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) DelayAfterDelete(value string) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.delayAfterDelete = value
	b.fieldSet_[1] = true
	return b
}

// DelayAfterFailure sets the value of the 'delay_after_failure' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) DelayAfterFailure(value string) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.delayAfterFailure = value
	b.fieldSet_[2] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) Enabled(value bool) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.enabled = value
	b.fieldSet_[3] = true
	return b
}

// UnneededTime sets the value of the 'unneeded_time' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) UnneededTime(value string) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.unneededTime = value
	b.fieldSet_[4] = true
	return b
}

// UtilizationThreshold sets the value of the 'utilization_threshold' attribute to the given value.
func (b *AutoscalerScaleDownConfigBuilder) UtilizationThreshold(value string) *AutoscalerScaleDownConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.utilizationThreshold = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AutoscalerScaleDownConfigBuilder) Copy(object *AutoscalerScaleDownConfig) *AutoscalerScaleDownConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.delayAfterAdd = object.delayAfterAdd
	b.delayAfterDelete = object.delayAfterDelete
	b.delayAfterFailure = object.delayAfterFailure
	b.enabled = object.enabled
	b.unneededTime = object.unneededTime
	b.utilizationThreshold = object.utilizationThreshold
	return b
}

// Build creates a 'autoscaler_scale_down_config' object using the configuration stored in the builder.
func (b *AutoscalerScaleDownConfigBuilder) Build() (object *AutoscalerScaleDownConfig, err error) {
	object = new(AutoscalerScaleDownConfig)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.delayAfterAdd = b.delayAfterAdd
	object.delayAfterDelete = b.delayAfterDelete
	object.delayAfterFailure = b.delayAfterFailure
	object.enabled = b.enabled
	object.unneededTime = b.unneededTime
	object.utilizationThreshold = b.utilizationThreshold
	return
}
