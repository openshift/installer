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

// ControlPlaneUpgradePolicyBuilder contains the data and logic needed to build 'control_plane_upgrade_policy' objects.
//
// Representation of an upgrade policy that can be set for a cluster.
type ControlPlaneUpgradePolicyBuilder struct {
	bitmap_                    uint32
	id                         string
	href                       string
	clusterID                  string
	creationTimestamp          time.Time
	lastUpdateTimestamp        time.Time
	nextRun                    time.Time
	schedule                   string
	scheduleType               ScheduleType
	state                      *UpgradePolicyStateBuilder
	upgradeType                UpgradeType
	version                    string
	enableMinorVersionUpgrades bool
}

// NewControlPlaneUpgradePolicy creates a new builder of 'control_plane_upgrade_policy' objects.
func NewControlPlaneUpgradePolicy() *ControlPlaneUpgradePolicyBuilder {
	return &ControlPlaneUpgradePolicyBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ControlPlaneUpgradePolicyBuilder) Link(value bool) *ControlPlaneUpgradePolicyBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ControlPlaneUpgradePolicyBuilder) ID(value string) *ControlPlaneUpgradePolicyBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ControlPlaneUpgradePolicyBuilder) HREF(value string) *ControlPlaneUpgradePolicyBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ControlPlaneUpgradePolicyBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) ClusterID(value string) *ControlPlaneUpgradePolicyBuilder {
	b.clusterID = value
	b.bitmap_ |= 8
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) CreationTimestamp(value time.Time) *ControlPlaneUpgradePolicyBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 16
	return b
}

// EnableMinorVersionUpgrades sets the value of the 'enable_minor_version_upgrades' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) EnableMinorVersionUpgrades(value bool) *ControlPlaneUpgradePolicyBuilder {
	b.enableMinorVersionUpgrades = value
	b.bitmap_ |= 32
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) LastUpdateTimestamp(value time.Time) *ControlPlaneUpgradePolicyBuilder {
	b.lastUpdateTimestamp = value
	b.bitmap_ |= 64
	return b
}

// NextRun sets the value of the 'next_run' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) NextRun(value time.Time) *ControlPlaneUpgradePolicyBuilder {
	b.nextRun = value
	b.bitmap_ |= 128
	return b
}

// Schedule sets the value of the 'schedule' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) Schedule(value string) *ControlPlaneUpgradePolicyBuilder {
	b.schedule = value
	b.bitmap_ |= 256
	return b
}

// ScheduleType sets the value of the 'schedule_type' attribute to the given value.
//
// ScheduleType defines which type of scheduling should be used for the upgrade policy.
func (b *ControlPlaneUpgradePolicyBuilder) ScheduleType(value ScheduleType) *ControlPlaneUpgradePolicyBuilder {
	b.scheduleType = value
	b.bitmap_ |= 512
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of an upgrade policy state that that is set for a cluster.
func (b *ControlPlaneUpgradePolicyBuilder) State(value *UpgradePolicyStateBuilder) *ControlPlaneUpgradePolicyBuilder {
	b.state = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// UpgradeType sets the value of the 'upgrade_type' attribute to the given value.
//
// UpgradeType defines which type of upgrade should be used.
func (b *ControlPlaneUpgradePolicyBuilder) UpgradeType(value UpgradeType) *ControlPlaneUpgradePolicyBuilder {
	b.upgradeType = value
	b.bitmap_ |= 2048
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) Version(value string) *ControlPlaneUpgradePolicyBuilder {
	b.version = value
	b.bitmap_ |= 4096
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ControlPlaneUpgradePolicyBuilder) Copy(object *ControlPlaneUpgradePolicy) *ControlPlaneUpgradePolicyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.creationTimestamp = object.creationTimestamp
	b.enableMinorVersionUpgrades = object.enableMinorVersionUpgrades
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	b.nextRun = object.nextRun
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

// Build creates a 'control_plane_upgrade_policy' object using the configuration stored in the builder.
func (b *ControlPlaneUpgradePolicyBuilder) Build() (object *ControlPlaneUpgradePolicy, err error) {
	object = new(ControlPlaneUpgradePolicy)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.clusterID = b.clusterID
	object.creationTimestamp = b.creationTimestamp
	object.enableMinorVersionUpgrades = b.enableMinorVersionUpgrades
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	object.nextRun = b.nextRun
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
