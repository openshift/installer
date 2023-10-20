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

// SummaryDashboardBuilder contains the data and logic needed to build 'summary_dashboard' objects.
type SummaryDashboardBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	metrics []*SummaryMetricsBuilder
}

// NewSummaryDashboard creates a new builder of 'summary_dashboard' objects.
func NewSummaryDashboard() *SummaryDashboardBuilder {
	return &SummaryDashboardBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SummaryDashboardBuilder) Link(value bool) *SummaryDashboardBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SummaryDashboardBuilder) ID(value string) *SummaryDashboardBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SummaryDashboardBuilder) HREF(value string) *SummaryDashboardBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SummaryDashboardBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Metrics sets the value of the 'metrics' attribute to the given values.
func (b *SummaryDashboardBuilder) Metrics(values ...*SummaryMetricsBuilder) *SummaryDashboardBuilder {
	b.metrics = make([]*SummaryMetricsBuilder, len(values))
	copy(b.metrics, values)
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SummaryDashboardBuilder) Copy(object *SummaryDashboard) *SummaryDashboardBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.metrics != nil {
		b.metrics = make([]*SummaryMetricsBuilder, len(object.metrics))
		for i, v := range object.metrics {
			b.metrics[i] = NewSummaryMetrics().Copy(v)
		}
	} else {
		b.metrics = nil
	}
	return b
}

// Build creates a 'summary_dashboard' object using the configuration stored in the builder.
func (b *SummaryDashboardBuilder) Build() (object *SummaryDashboard, err error) {
	object = new(SummaryDashboard)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.metrics != nil {
		object.metrics = make([]*SummaryMetrics, len(b.metrics))
		for i, v := range b.metrics {
			object.metrics[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
