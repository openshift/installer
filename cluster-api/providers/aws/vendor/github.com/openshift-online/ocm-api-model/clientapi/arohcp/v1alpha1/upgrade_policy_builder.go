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

// Representation of an upgrade policy that can be set for a cluster.
type UpgradePolicyBuilder struct {
	fieldSet_                  []bool
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
	return &UpgradePolicyBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *UpgradePolicyBuilder) Link(value bool) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *UpgradePolicyBuilder) ID(value string) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *UpgradePolicyBuilder) HREF(value string) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *UpgradePolicyBuilder) Empty() bool {
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
func (b *UpgradePolicyBuilder) ClusterID(value string) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// EnableMinorVersionUpgrades sets the value of the 'enable_minor_version_upgrades' attribute to the given value.
func (b *UpgradePolicyBuilder) EnableMinorVersionUpgrades(value bool) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.enableMinorVersionUpgrades = value
	b.fieldSet_[4] = true
	return b
}

// NextRun sets the value of the 'next_run' attribute to the given value.
func (b *UpgradePolicyBuilder) NextRun(value time.Time) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.nextRun = value
	b.fieldSet_[5] = true
	return b
}

// Schedule sets the value of the 'schedule' attribute to the given value.
func (b *UpgradePolicyBuilder) Schedule(value string) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.schedule = value
	b.fieldSet_[6] = true
	return b
}

// ScheduleType sets the value of the 'schedule_type' attribute to the given value.
//
// ScheduleType defines which type of scheduling should be used for the upgrade policy.
func (b *UpgradePolicyBuilder) ScheduleType(value ScheduleType) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.scheduleType = value
	b.fieldSet_[7] = true
	return b
}

// UpgradeType sets the value of the 'upgrade_type' attribute to the given value.
//
// UpgradeType defines which type of upgrade should be used.
func (b *UpgradePolicyBuilder) UpgradeType(value UpgradeType) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.upgradeType = value
	b.fieldSet_[8] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *UpgradePolicyBuilder) Version(value string) *UpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.version = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UpgradePolicyBuilder) Copy(object *UpgradePolicy) *UpgradePolicyBuilder {
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterID = b.clusterID
	object.enableMinorVersionUpgrades = b.enableMinorVersionUpgrades
	object.nextRun = b.nextRun
	object.schedule = b.schedule
	object.scheduleType = b.scheduleType
	object.upgradeType = b.upgradeType
	object.version = b.version
	return
}
