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

// UpgradePolicyBuilder contains the data and logic needed to build 'upgrade_policy' objects.
//
// Representation of an upgrade policy that can be set for a cluster.
type UpgradePolicyBuilder struct {
	bitmap_                    uint32
	id                         string
	href                       string
	clusterID                  string
	nextRun                    time.Time
	schedule                   string
	scheduleType               ScheduleType
	upgradeType                UpgradeType
	version                    string
	enableMinorVersionUpgrades bool
}

// NewUpgradePolicy creates a new builder of 'upgrade_policy' objects.
func NewUpgradePolicy() *UpgradePolicyBuilder {
	return &UpgradePolicyBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *UpgradePolicyBuilder) Link(value bool) *UpgradePolicyBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *UpgradePolicyBuilder) ID(value string) *UpgradePolicyBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *UpgradePolicyBuilder) HREF(value string) *UpgradePolicyBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *UpgradePolicyBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *UpgradePolicyBuilder) ClusterID(value string) *UpgradePolicyBuilder {
	b.clusterID = value
	b.bitmap_ |= 8
	return b
}

// EnableMinorVersionUpgrades sets the value of the 'enable_minor_version_upgrades' attribute to the given value.
func (b *UpgradePolicyBuilder) EnableMinorVersionUpgrades(value bool) *UpgradePolicyBuilder {
	b.enableMinorVersionUpgrades = value
	b.bitmap_ |= 16
	return b
}

// NextRun sets the value of the 'next_run' attribute to the given value.
func (b *UpgradePolicyBuilder) NextRun(value time.Time) *UpgradePolicyBuilder {
	b.nextRun = value
	b.bitmap_ |= 32
	return b
}

// Schedule sets the value of the 'schedule' attribute to the given value.
func (b *UpgradePolicyBuilder) Schedule(value string) *UpgradePolicyBuilder {
	b.schedule = value
	b.bitmap_ |= 64
	return b
}

// ScheduleType sets the value of the 'schedule_type' attribute to the given value.
//
// ScheduleType defines which type of scheduling should be used for the upgrade policy.
func (b *UpgradePolicyBuilder) ScheduleType(value ScheduleType) *UpgradePolicyBuilder {
	b.scheduleType = value
	b.bitmap_ |= 128
	return b
}

// UpgradeType sets the value of the 'upgrade_type' attribute to the given value.
//
// UpgradeType defines which type of upgrade should be used.
func (b *UpgradePolicyBuilder) UpgradeType(value UpgradeType) *UpgradePolicyBuilder {
	b.upgradeType = value
	b.bitmap_ |= 256
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *UpgradePolicyBuilder) Version(value string) *UpgradePolicyBuilder {
	b.version = value
	b.bitmap_ |= 512
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UpgradePolicyBuilder) Copy(object *UpgradePolicy) *UpgradePolicyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.enableMinorVersionUpgrades = object.enableMinorVersionUpgrades
	b.nextRun = object.nextRun
	b.schedule = object.schedule
	b.scheduleType = object.scheduleType
	b.upgradeType = object.upgradeType
	b.version = object.version
	return b
}

// Build creates a 'upgrade_policy' object using the configuration stored in the builder.
func (b *UpgradePolicyBuilder) Build() (object *UpgradePolicy, err error) {
	object = new(UpgradePolicy)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.clusterID = b.clusterID
	object.enableMinorVersionUpgrades = b.enableMinorVersionUpgrades
	object.nextRun = b.nextRun
	object.schedule = b.schedule
	object.scheduleType = b.scheduleType
	object.upgradeType = b.upgradeType
	object.version = b.version
	return
}
