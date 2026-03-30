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

import (
	time "time"
)

// Representation of information from telemetry about a the CPU capacity by node role and OS.
type CPUTotalNodeRoleOSMetricNodeBuilder struct {
	fieldSet_       []bool
	cpuTotal        float64
	nodeRoles       []string
	operatingSystem string
	time            time.Time
}

// NewCPUTotalNodeRoleOSMetricNode creates a new builder of 'CPU_total_node_role_OS_metric_node' objects.
func NewCPUTotalNodeRoleOSMetricNode() *CPUTotalNodeRoleOSMetricNodeBuilder {
	return &CPUTotalNodeRoleOSMetricNodeBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Empty() bool {
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

// CPUTotal sets the value of the 'CPU_total' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) CPUTotal(value float64) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.cpuTotal = value
	b.fieldSet_[0] = true
	return b
}

// NodeRoles sets the value of the 'node_roles' attribute to the given values.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) NodeRoles(values ...string) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.nodeRoles = make([]string, len(values))
	copy(b.nodeRoles, values)
	b.fieldSet_[1] = true
	return b
}

// OperatingSystem sets the value of the 'operating_system' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) OperatingSystem(value string) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.operatingSystem = value
	b.fieldSet_[2] = true
	return b
}

// Time sets the value of the 'time' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Time(value time.Time) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.time = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Copy(object *CPUTotalNodeRoleOSMetricNode) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.cpuTotal = object.cpuTotal
	if object.nodeRoles != nil {
		b.nodeRoles = make([]string, len(object.nodeRoles))
		copy(b.nodeRoles, object.nodeRoles)
	} else {
		b.nodeRoles = nil
	}
	b.operatingSystem = object.operatingSystem
	b.time = object.time
	return b
}

// Build creates a 'CPU_total_node_role_OS_metric_node' object using the configuration stored in the builder.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Build() (object *CPUTotalNodeRoleOSMetricNode, err error) {
	object = new(CPUTotalNodeRoleOSMetricNode)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.cpuTotal = b.cpuTotal
	if b.nodeRoles != nil {
		object.nodeRoles = make([]string, len(b.nodeRoles))
		copy(object.nodeRoles, b.nodeRoles)
	}
	object.operatingSystem = b.operatingSystem
	object.time = b.time
	return
}
