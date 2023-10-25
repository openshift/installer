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

// SummaryMetricsBuilder contains the data and logic needed to build 'summary_metrics' objects.
type SummaryMetricsBuilder struct {
	bitmap_ uint32
	name    string
	vector  []*SummarySampleBuilder
}

// NewSummaryMetrics creates a new builder of 'summary_metrics' objects.
func NewSummaryMetrics() *SummaryMetricsBuilder {
	return &SummaryMetricsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SummaryMetricsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *SummaryMetricsBuilder) Name(value string) *SummaryMetricsBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Vector sets the value of the 'vector' attribute to the given values.
func (b *SummaryMetricsBuilder) Vector(values ...*SummarySampleBuilder) *SummaryMetricsBuilder {
	b.vector = make([]*SummarySampleBuilder, len(values))
	copy(b.vector, values)
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SummaryMetricsBuilder) Copy(object *SummaryMetrics) *SummaryMetricsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	if object.vector != nil {
		b.vector = make([]*SummarySampleBuilder, len(object.vector))
		for i, v := range object.vector {
			b.vector[i] = NewSummarySample().Copy(v)
		}
	} else {
		b.vector = nil
	}
	return b
}

// Build creates a 'summary_metrics' object using the configuration stored in the builder.
func (b *SummaryMetricsBuilder) Build() (object *SummaryMetrics, err error) {
	object = new(SummaryMetrics)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	if b.vector != nil {
		object.vector = make([]*SummarySample, len(b.vector))
		for i, v := range b.vector {
			object.vector[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
