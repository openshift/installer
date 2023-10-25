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

// SocketTotalsNodeRoleOSMetricNodeBuilder contains the data and logic needed to build 'socket_totals_node_role_OS_metric_node' objects.
//
// Representation of information from telemetry about the socket capacity by node
// role and OS of a cluster.
type SocketTotalsNodeRoleOSMetricNodeBuilder struct {
	bitmap_      uint32
	socketTotals []*SocketTotalNodeRoleOSMetricNodeBuilder
}

// NewSocketTotalsNodeRoleOSMetricNode creates a new builder of 'socket_totals_node_role_OS_metric_node' objects.
func NewSocketTotalsNodeRoleOSMetricNode() *SocketTotalsNodeRoleOSMetricNodeBuilder {
	return &SocketTotalsNodeRoleOSMetricNodeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// SocketTotals sets the value of the 'socket_totals' attribute to the given values.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) SocketTotals(values ...*SocketTotalNodeRoleOSMetricNodeBuilder) *SocketTotalsNodeRoleOSMetricNodeBuilder {
	b.socketTotals = make([]*SocketTotalNodeRoleOSMetricNodeBuilder, len(values))
	copy(b.socketTotals, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SocketTotalsNodeRoleOSMetricNodeBuilder) Copy(object *SocketTotalsNodeRoleOSMetricNode) *SocketTotalsNodeRoleOSMetricNodeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
