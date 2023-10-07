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

// AddonUpgradePolicyStateBuilder contains the data and logic needed to build 'addon_upgrade_policy_state' objects.
//
// Representation of an addon upgrade policy state that that is set for a cluster.
type AddonUpgradePolicyStateBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	description string
	value       UpgradePolicyStateValue
}

// NewAddonUpgradePolicyState creates a new builder of 'addon_upgrade_policy_state' objects.
func NewAddonUpgradePolicyState() *AddonUpgradePolicyStateBuilder {
	return &AddonUpgradePolicyStateBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonUpgradePolicyStateBuilder) Link(value bool) *AddonUpgradePolicyStateBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddonUpgradePolicyStateBuilder) ID(value string) *AddonUpgradePolicyStateBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddonUpgradePolicyStateBuilder) HREF(value string) *AddonUpgradePolicyStateBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonUpgradePolicyStateBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Description sets the value of the 'description' attribute to the given value.
func (b *AddonUpgradePolicyStateBuilder) Description(value string) *AddonUpgradePolicyStateBuilder {
	b.description = value
	b.bitmap_ |= 8
	return b
}

// Value sets the value of the 'value' attribute to the given value.
//
// Overall state of a cluster upgrade policy.
func (b *AddonUpgradePolicyStateBuilder) Value(value UpgradePolicyStateValue) *AddonUpgradePolicyStateBuilder {
	b.value = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonUpgradePolicyStateBuilder) Copy(object *AddonUpgradePolicyState) *AddonUpgradePolicyStateBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.description = object.description
	b.value = object.value
	return b
}

// Build creates a 'addon_upgrade_policy_state' object using the configuration stored in the builder.
func (b *AddonUpgradePolicyStateBuilder) Build() (object *AddonUpgradePolicyState, err error) {
	object = new(AddonUpgradePolicyState)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.description = b.description
	object.value = b.value
	return
}
