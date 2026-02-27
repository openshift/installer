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

import (
	time "time"
)

// Description of an environment
type EnvironmentBuilder struct {
	fieldSet_                       []bool
	backplaneURL                    string
	lastClusterImagesetSync         time.Time
	lastHibernationCheck            time.Time
	lastLimitedSupportCheck         time.Time
	lastLimitedSupportOverrideCheck time.Time
	lastUpgradeAvailableCheck       time.Time
	name                            string
}

// NewEnvironment creates a new builder of 'environment' objects.
func NewEnvironment() *EnvironmentBuilder {
	return &EnvironmentBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EnvironmentBuilder) Empty() bool {
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

// BackplaneURL sets the value of the 'backplane_URL' attribute to the given value.
func (b *EnvironmentBuilder) BackplaneURL(value string) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.backplaneURL = value
	b.fieldSet_[0] = true
	return b
}

// LastClusterImagesetSync sets the value of the 'last_cluster_imageset_sync' attribute to the given value.
func (b *EnvironmentBuilder) LastClusterImagesetSync(value time.Time) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastClusterImagesetSync = value
	b.fieldSet_[1] = true
	return b
}

// LastHibernationCheck sets the value of the 'last_hibernation_check' attribute to the given value.
func (b *EnvironmentBuilder) LastHibernationCheck(value time.Time) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastHibernationCheck = value
	b.fieldSet_[2] = true
	return b
}

// LastLimitedSupportCheck sets the value of the 'last_limited_support_check' attribute to the given value.
func (b *EnvironmentBuilder) LastLimitedSupportCheck(value time.Time) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastLimitedSupportCheck = value
	b.fieldSet_[3] = true
	return b
}

// LastLimitedSupportOverrideCheck sets the value of the 'last_limited_support_override_check' attribute to the given value.
func (b *EnvironmentBuilder) LastLimitedSupportOverrideCheck(value time.Time) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastLimitedSupportOverrideCheck = value
	b.fieldSet_[4] = true
	return b
}

// LastUpgradeAvailableCheck sets the value of the 'last_upgrade_available_check' attribute to the given value.
func (b *EnvironmentBuilder) LastUpgradeAvailableCheck(value time.Time) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.lastUpgradeAvailableCheck = value
	b.fieldSet_[5] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *EnvironmentBuilder) Name(value string) *EnvironmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.name = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EnvironmentBuilder) Copy(object *Environment) *EnvironmentBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.backplaneURL = object.backplaneURL
	b.lastClusterImagesetSync = object.lastClusterImagesetSync
	b.lastHibernationCheck = object.lastHibernationCheck
	b.lastLimitedSupportCheck = object.lastLimitedSupportCheck
	b.lastLimitedSupportOverrideCheck = object.lastLimitedSupportOverrideCheck
	b.lastUpgradeAvailableCheck = object.lastUpgradeAvailableCheck
	b.name = object.name
	return b
}

// Build creates a 'environment' object using the configuration stored in the builder.
func (b *EnvironmentBuilder) Build() (object *Environment, err error) {
	object = new(Environment)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.backplaneURL = b.backplaneURL
	object.lastClusterImagesetSync = b.lastClusterImagesetSync
	object.lastHibernationCheck = b.lastHibernationCheck
	object.lastLimitedSupportCheck = b.lastLimitedSupportCheck
	object.lastLimitedSupportOverrideCheck = b.lastLimitedSupportOverrideCheck
	object.lastUpgradeAvailableCheck = b.lastUpgradeAvailableCheck
	object.name = b.name
	return
}
