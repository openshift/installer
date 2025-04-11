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

// NodePoolManagementUpgradeBuilder contains the data and logic needed to build 'node_pool_management_upgrade' objects.
//
// Representation of node pool management.
type NodePoolManagementUpgradeBuilder struct {
	bitmap_        uint32
	id             string
	href           string
	maxSurge       string
	maxUnavailable string
	type_          string
}

// NewNodePoolManagementUpgrade creates a new builder of 'node_pool_management_upgrade' objects.
func NewNodePoolManagementUpgrade() *NodePoolManagementUpgradeBuilder {
	return &NodePoolManagementUpgradeBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolManagementUpgradeBuilder) Link(value bool) *NodePoolManagementUpgradeBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolManagementUpgradeBuilder) ID(value string) *NodePoolManagementUpgradeBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *NodePoolManagementUpgradeBuilder) HREF(value string) *NodePoolManagementUpgradeBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolManagementUpgradeBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// MaxSurge sets the value of the 'max_surge' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) MaxSurge(value string) *NodePoolManagementUpgradeBuilder {
	b.maxSurge = value
	b.bitmap_ |= 8
	return b
}

// MaxUnavailable sets the value of the 'max_unavailable' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) MaxUnavailable(value string) *NodePoolManagementUpgradeBuilder {
	b.maxUnavailable = value
	b.bitmap_ |= 16
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) Type(value string) *NodePoolManagementUpgradeBuilder {
	b.type_ = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolManagementUpgradeBuilder) Copy(object *NodePoolManagementUpgrade) *NodePoolManagementUpgradeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.maxSurge = object.maxSurge
	b.maxUnavailable = object.maxUnavailable
	b.type_ = object.type_
	return b
}

// Build creates a 'node_pool_management_upgrade' object using the configuration stored in the builder.
func (b *NodePoolManagementUpgradeBuilder) Build() (object *NodePoolManagementUpgrade, err error) {
	object = new(NodePoolManagementUpgrade)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.maxSurge = b.maxSurge
	object.maxUnavailable = b.maxUnavailable
	object.type_ = b.type_
	return
}
