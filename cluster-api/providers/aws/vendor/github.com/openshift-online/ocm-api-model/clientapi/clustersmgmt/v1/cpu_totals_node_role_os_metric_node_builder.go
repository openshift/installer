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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of information from telemetry about the CPU capacity by node
// role and OS of a cluster.
type CPUTotalsNodeRoleOSMetricNodeBuilder struct {
	fieldSet_ []bool
	cpuTotals []*CPUTotalNodeRoleOSMetricNodeBuilder
}

// NewCPUTotalsNodeRoleOSMetricNode creates a new builder of 'CPU_totals_node_role_OS_metric_node' objects.
func NewCPUTotalsNodeRoleOSMetricNode() *CPUTotalsNodeRoleOSMetricNodeBuilder {
	return &CPUTotalsNodeRoleOSMetricNodeBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CPUTotalsNodeRoleOSMetricNodeBuilder) Empty() bool {
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

// CPUTotals sets the value of the 'CPU_totals' attribute to the given values.
func (b *CPUTotalsNodeRoleOSMetricNodeBuilder) CPUTotals(values ...*CPUTotalNodeRoleOSMetricNodeBuilder) *CPUTotalsNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.cpuTotals = make([]*CPUTotalNodeRoleOSMetricNodeBuilder, len(values))
	copy(b.cpuTotals, values)
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CPUTotalsNodeRoleOSMetricNodeBuilder) Copy(object *CPUTotalsNodeRoleOSMetricNode) *CPUTotalsNodeRoleOSMetricNodeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.cpuTotals != nil {
		b.cpuTotals = make([]*CPUTotalNodeRoleOSMetricNodeBuilder, len(object.cpuTotals))
		for i, v := range object.cpuTotals {
			b.cpuTotals[i] = NewCPUTotalNodeRoleOSMetricNode().Copy(v)
		}
	} else {
		b.cpuTotals = nil
	}
	return b
}

// Build creates a 'CPU_totals_node_role_OS_metric_node' object using the configuration stored in the builder.
func (b *CPUTotalsNodeRoleOSMetricNodeBuilder) Build() (object *CPUTotalsNodeRoleOSMetricNode, err error) {
	object = new(CPUTotalsNodeRoleOSMetricNode)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.cpuTotals != nil {
		object.cpuTotals = make([]*CPUTotalNodeRoleOSMetricNode, len(b.cpuTotals))
		for i, v := range b.cpuTotals {
			object.cpuTotals[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
