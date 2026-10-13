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

type ClusterMetricsNodesBuilder struct {
	fieldSet_ []bool
	compute   float64
	infra     float64
	master    float64
	total     float64
}

// NewClusterMetricsNodes creates a new builder of 'cluster_metrics_nodes' objects.
func NewClusterMetricsNodes() *ClusterMetricsNodesBuilder {
	return &ClusterMetricsNodesBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMetricsNodesBuilder) Empty() bool {
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

// Compute sets the value of the 'compute' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Compute(value float64) *ClusterMetricsNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.compute = value
	b.fieldSet_[0] = true
	return b
}

// Infra sets the value of the 'infra' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Infra(value float64) *ClusterMetricsNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.infra = value
	b.fieldSet_[1] = true
	return b
}

// Master sets the value of the 'master' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Master(value float64) *ClusterMetricsNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.master = value
	b.fieldSet_[2] = true
	return b
}

// Total sets the value of the 'total' attribute to the given value.
func (b *ClusterMetricsNodesBuilder) Total(value float64) *ClusterMetricsNodesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.total = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMetricsNodesBuilder) Copy(object *ClusterMetricsNodes) *ClusterMetricsNodesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.compute = object.compute
	b.infra = object.infra
	b.master = object.master
	b.total = object.total
	return b
}

// Build creates a 'cluster_metrics_nodes' object using the configuration stored in the builder.
func (b *ClusterMetricsNodesBuilder) Build() (object *ClusterMetricsNodes, err error) {
	object = new(ClusterMetricsNodes)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.compute = b.compute
	object.infra = b.infra
	object.master = b.master
	object.total = b.total
	return
}
