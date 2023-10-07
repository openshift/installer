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

import (
	time "time"
)

// CPUTotalNodeRoleOSMetricNodeBuilder contains the data and logic needed to build 'CPU_total_node_role_OS_metric_node' objects.
//
// Representation of information from telemetry about a the CPU capacity by node role and OS.
type CPUTotalNodeRoleOSMetricNodeBuilder struct {
	bitmap_         uint32
	cpuTotal        float64
	nodeRoles       []string
	operatingSystem string
	time            time.Time
}

// NewCPUTotalNodeRoleOSMetricNode creates a new builder of 'CPU_total_node_role_OS_metric_node' objects.
func NewCPUTotalNodeRoleOSMetricNode() *CPUTotalNodeRoleOSMetricNodeBuilder {
	return &CPUTotalNodeRoleOSMetricNodeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CPUTotal sets the value of the 'CPU_total' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) CPUTotal(value float64) *CPUTotalNodeRoleOSMetricNodeBuilder {
	b.cpuTotal = value
	b.bitmap_ |= 1
	return b
}

// NodeRoles sets the value of the 'node_roles' attribute to the given values.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) NodeRoles(values ...string) *CPUTotalNodeRoleOSMetricNodeBuilder {
	b.nodeRoles = make([]string, len(values))
	copy(b.nodeRoles, values)
	b.bitmap_ |= 2
	return b
}

// OperatingSystem sets the value of the 'operating_system' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) OperatingSystem(value string) *CPUTotalNodeRoleOSMetricNodeBuilder {
	b.operatingSystem = value
	b.bitmap_ |= 4
	return b
}

// Time sets the value of the 'time' attribute to the given value.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Time(value time.Time) *CPUTotalNodeRoleOSMetricNodeBuilder {
	b.time = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CPUTotalNodeRoleOSMetricNodeBuilder) Copy(object *CPUTotalNodeRoleOSMetricNode) *CPUTotalNodeRoleOSMetricNodeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.cpuTotal = b.cpuTotal
	if b.nodeRoles != nil {
		object.nodeRoles = make([]string, len(b.nodeRoles))
		copy(object.nodeRoles, b.nodeRoles)
	}
	object.operatingSystem = b.operatingSystem
	object.time = b.time
	return
}
