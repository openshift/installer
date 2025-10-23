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

// Representation of node pool management.
type NodePoolManagementUpgradeBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	maxSurge       string
	maxUnavailable string
	type_          string
}

// NewNodePoolManagementUpgrade creates a new builder of 'node_pool_management_upgrade' objects.
func NewNodePoolManagementUpgrade() *NodePoolManagementUpgradeBuilder {
	return &NodePoolManagementUpgradeBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolManagementUpgradeBuilder) Link(value bool) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolManagementUpgradeBuilder) ID(value string) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *NodePoolManagementUpgradeBuilder) HREF(value string) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolManagementUpgradeBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// MaxSurge sets the value of the 'max_surge' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) MaxSurge(value string) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.maxSurge = value
	b.fieldSet_[3] = true
	return b
}

// MaxUnavailable sets the value of the 'max_unavailable' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) MaxUnavailable(value string) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.maxUnavailable = value
	b.fieldSet_[4] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *NodePoolManagementUpgradeBuilder) Type(value string) *NodePoolManagementUpgradeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.type_ = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolManagementUpgradeBuilder) Copy(object *NodePoolManagementUpgrade) *NodePoolManagementUpgradeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.maxSurge = b.maxSurge
	object.maxUnavailable = b.maxUnavailable
	object.type_ = b.type_
	return
}
