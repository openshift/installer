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

// Representation of information from telemetry about the socket capacity by node
// role and OS of a cluster.
type SocketTotalsNodeRoleOSMetricNodeBuilder struct {
	fieldSet_    []bool
	socketTotals []*SocketTotalNodeRoleOSMetricNodeBuilder
}

// NewSocketTotalsNodeRoleOSMetricNode creates a new builder of 'socket_totals_node_role_OS_metric_node' objects.
func NewSocketTotalsNodeRoleOSMetricNode() *SocketTotalsNodeRoleOSMetricNodeBuilder {
	return &SocketTotalsNodeRoleOSMetricNodeBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) Empty() bool {
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

// SocketTotals sets the value of the 'socket_totals' attribute to the given values.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) SocketTotals(values ...*SocketTotalNodeRoleOSMetricNodeBuilder) *SocketTotalsNodeRoleOSMetricNodeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.socketTotals = make([]*SocketTotalNodeRoleOSMetricNodeBuilder, len(values))
	copy(b.socketTotals, values)
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) Copy(object *SocketTotalsNodeRoleOSMetricNode) *SocketTotalsNodeRoleOSMetricNodeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.socketTotals != nil {
		b.socketTotals = make([]*SocketTotalNodeRoleOSMetricNodeBuilder, len(object.socketTotals))
		for i, v := range object.socketTotals {
			b.socketTotals[i] = NewSocketTotalNodeRoleOSMetricNode().Copy(v)
		}
	} else {
		b.socketTotals = nil
	}
	return b
}

// Build creates a 'socket_totals_node_role_OS_metric_node' object using the configuration stored in the builder.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) Build() (object *SocketTotalsNodeRoleOSMetricNode, err error) {
	object = new(SocketTotalsNodeRoleOSMetricNode)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.socketTotals != nil {
		object.socketTotals = make([]*SocketTotalNodeRoleOSMetricNode, len(b.socketTotals))
		for i, v := range b.socketTotals {
			object.socketTotals[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
