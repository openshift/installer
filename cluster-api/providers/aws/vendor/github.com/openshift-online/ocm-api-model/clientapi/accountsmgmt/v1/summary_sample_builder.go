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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

type SummarySampleBuilder struct {
	fieldSet_ []bool
	time      string
	value     float64
}

// NewSummarySample creates a new builder of 'summary_sample' objects.
func NewSummarySample() *SummarySampleBuilder {
	return &SummarySampleBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SummarySampleBuilder) Empty() bool {
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

// Time sets the value of the 'time' attribute to the given value.
func (b *SummarySampleBuilder) Time(value string) *SummarySampleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.time = value
	b.fieldSet_[0] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *SummarySampleBuilder) Value(value float64) *SummarySampleBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.value = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SummarySampleBuilder) Copy(object *SummarySample) *SummarySampleBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.time = object.time
	b.value = object.value
	return b
}

// Build creates a 'summary_sample' object using the configuration stored in the builder.
func (b *SummarySampleBuilder) Build() (object *SummarySample, err error) {
	object = new(SummarySample)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.time = b.time
	object.value = b.value
	return
}
