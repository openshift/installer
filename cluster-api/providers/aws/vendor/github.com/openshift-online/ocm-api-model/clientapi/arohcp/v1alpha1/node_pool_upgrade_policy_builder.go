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

// Representation of an upgrade policy that can be set for a node pool.
type NodePoolUpgradePolicyBuilder struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	clusterID                  string
	creationTimestamp          time.Time
	lastUpdateTimestamp        time.Time
	nextRun                    time.Time
	nodePoolID                 string
	schedule                   string
	scheduleType               ScheduleType
	state                      *UpgradePolicyStateBuilder
	upgradeType                UpgradeType
	version                    string
	enableMinorVersionUpgrades bool
}

// NewNodePoolUpgradePolicy creates a new builder of 'node_pool_upgrade_policy' objects.
func NewNodePoolUpgradePolicy() *NodePoolUpgradePolicyBuilder {
	return &NodePoolUpgradePolicyBuilder{
		fieldSet_: make([]bool, 14),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *NodePoolUpgradePolicyBuilder) Link(value bool) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *NodePoolUpgradePolicyBuilder) ID(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *NodePoolUpgradePolicyBuilder) HREF(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NodePoolUpgradePolicyBuilder) Empty() bool {
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

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) ClusterID(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) CreationTimestamp(value time.Time) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.creationTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// EnableMinorVersionUpgrades sets the value of the 'enable_minor_version_upgrades' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) EnableMinorVersionUpgrades(value bool) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.enableMinorVersionUpgrades = value
	b.fieldSet_[5] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) LastUpdateTimestamp(value time.Time) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// NextRun sets the value of the 'next_run' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) NextRun(value time.Time) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.nextRun = value
	b.fieldSet_[7] = true
	return b
}

// NodePoolID sets the value of the 'node_pool_ID' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) NodePoolID(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.nodePoolID = value
	b.fieldSet_[8] = true
	return b
}

// Schedule sets the value of the 'schedule' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) Schedule(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.schedule = value
	b.fieldSet_[9] = true
	return b
}

// ScheduleType sets the value of the 'schedule_type' attribute to the given value.
//
// ScheduleType defines which type of scheduling should be used for the upgrade policy.
func (b *NodePoolUpgradePolicyBuilder) ScheduleType(value ScheduleType) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.scheduleType = value
	b.fieldSet_[10] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of an upgrade policy state that that is set for a cluster.
func (b *NodePoolUpgradePolicyBuilder) State(value *UpgradePolicyStateBuilder) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.state = value
	if value != nil {
		b.fieldSet_[11] = true
	} else {
		b.fieldSet_[11] = false
	}
	return b
}

// UpgradeType sets the value of the 'upgrade_type' attribute to the given value.
//
// UpgradeType defines which type of upgrade should be used.
func (b *NodePoolUpgradePolicyBuilder) UpgradeType(value UpgradeType) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.upgradeType = value
	b.fieldSet_[12] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *NodePoolUpgradePolicyBuilder) Version(value string) *NodePoolUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.version = value
	b.fieldSet_[13] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NodePoolUpgradePolicyBuilder) Copy(object *NodePoolUpgradePolicy) *NodePoolUpgradePolicyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.creationTimestamp = object.creationTimestamp
	b.enableMinorVersionUpgrades = object.enableMinorVersionUpgrades
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	b.nextRun = object.nextRun
	b.nodePoolID = object.nodePoolID
	b.schedule = object.schedule
	b.scheduleType = object.scheduleType
	if object.state != nil {
		b.state = NewUpgradePolicyState().Copy(object.state)
	} else {
		b.state = nil
	}
	b.upgradeType = object.upgradeType
	b.version = object.version
	return b
}

// Build creates a 'node_pool_upgrade_policy' object using the configuration stored in the builder.
func (b *NodePoolUpgradePolicyBuilder) Build() (object *NodePoolUpgradePolicy, err error) {
	object = new(NodePoolUpgradePolicy)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterID = b.clusterID
	object.creationTimestamp = b.creationTimestamp
	object.enableMinorVersionUpgrades = b.enableMinorVersionUpgrades
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	object.nextRun = b.nextRun
	object.nodePoolID = b.nodePoolID
	object.schedule = b.schedule
	object.scheduleType = b.scheduleType
	if b.state != nil {
		object.state, err = b.state.Build()
		if err != nil {
			return
		}
	}
	object.upgradeType = b.upgradeType
	object.version = b.version
	return
}
