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

// ClusterMetricsNodesBuilder contains the data and logic needed to build 'cluster_metrics_nodes' objects.
type ClusterMetricsNodesBuilder struct {
	bitmap_ uint32
	compute float64
	infra   float64
	master  float64
	total   float64
}

// NewClusterMetricsNodes creates a new builder of 'cluster_metrics_nodes' objects.
func NewClusterMetricsNodes() *ClusterMetricsNodesBuilder {
	return &ClusterMetricsNodesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMetricsNodesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Compute sets the value of the 'compute' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Compute(value float64) *ClusterMetricsNodesBuilder {
	b.compute = value
	b.bitmap_ |= 1
	return b
}

// Infra sets the value of the 'infra' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Infra(value float64) *ClusterMetricsNodesBuilder {
	b.infra = value
	b.bitmap_ |= 2
	return b
}

// Master sets the value of the 'master' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Master(value float64) *ClusterMetricsNodesBuilder {
	b.master = value
	b.bitmap_ |= 4
	return b
}

// Total sets the value of the 'total' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Total(value float64) *ClusterMetricsNodesBuilder {
	b.total = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMetricsNodesBuilder) Copy(object *ClusterMetricsNodes) *ClusterMetricsNodesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.compute = object.compute
	b.infra = object.infra
	b.master = object.master
	b.total = object.total
	return b
}

// Build creates a 'cluster_metrics_nodes' object using the configuration stored in the builder.
func (b *ClusterMetricsNodesBuilder) Build() (object *ClusterMetricsNodes, err error) {
	object = new(ClusterMetricsNodes)
	object.bitmap_ = b.bitmap_
	object.compute = b.compute
	object.infra = b.infra
	object.master = b.master
	object.total = b.total
	return
}
