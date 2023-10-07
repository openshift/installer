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

// EnvironmentBuilder contains the data and logic needed to build 'environment' objects.
//
// Description of an environment
type EnvironmentBuilder struct {
	bitmap_                   uint32
	lastLimitedSupportCheck   time.Time
	lastUpgradeAvailableCheck time.Time
	name                      string
}

// NewEnvironment creates a new builder of 'environment' objects.
func NewEnvironment() *EnvironmentBuilder {
	return &EnvironmentBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EnvironmentBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// LastLimitedSupportCheck sets the value of the 'last_limited_support_check' attribute to the given value.
func (b *EnvironmentBuilder) LastLimitedSupportCheck(value time.Time) *EnvironmentBuilder {
	b.lastLimitedSupportCheck = value
	b.bitmap_ |= 1
	return b
}

// LastUpgradeAvailableCheck sets the value of the 'last_upgrade_available_check' attribute to the given value.
func (b *EnvironmentBuilder) LastUpgradeAvailableCheck(value time.Time) *EnvironmentBuilder {
	b.lastUpgradeAvailableCheck = value
	b.bitmap_ |= 2
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *EnvironmentBuilder) Name(value string) *EnvironmentBuilder {
	b.name = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EnvironmentBuilder) Copy(object *Environment) *EnvironmentBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.lastLimitedSupportCheck = object.lastLimitedSupportCheck
	b.lastUpgradeAvailableCheck = object.lastUpgradeAvailableCheck
	b.name = object.name
	return b
}

// Build creates a 'environment' object using the configuration stored in the builder.
func (b *EnvironmentBuilder) Build() (object *Environment, err error) {
	object = new(Environment)
	object.bitmap_ = b.bitmap_
	object.lastLimitedSupportCheck = b.lastLimitedSupportCheck
	object.lastUpgradeAvailableCheck = b.lastUpgradeAvailableCheck
	object.name = b.name
	return
}
