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

// UpgradePolicyStateBuilder contains the data and logic needed to build 'upgrade_policy_state' objects.
//
// Representation of an upgrade policy state that that is set for a cluster.
type UpgradePolicyStateBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	description string
	value       UpgradePolicyStateValue
}

// NewUpgradePolicyState creates a new builder of 'upgrade_policy_state' objects.
func NewUpgradePolicyState() *UpgradePolicyStateBuilder {
	return &UpgradePolicyStateBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *UpgradePolicyStateBuilder) Link(value bool) *UpgradePolicyStateBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *UpgradePolicyStateBuilder) ID(value string) *UpgradePolicyStateBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *UpgradePolicyStateBuilder) HREF(value string) *UpgradePolicyStateBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *UpgradePolicyStateBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Description sets the value of the 'description' attribute to the given value.
func (b *UpgradePolicyStateBuilder) Description(value string) *UpgradePolicyStateBuilder {
	b.description = value
	b.bitmap_ |= 8
	return b
}

// Value sets the value of the 'value' attribute to the given value.
//
// Overall state of a cluster upgrade policy.
func (b *UpgradePolicyStateBuilder) Value(value UpgradePolicyStateValue) *UpgradePolicyStateBuilder {
	b.value = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *UpgradePolicyStateBuilder) Copy(object *UpgradePolicyState) *UpgradePolicyStateBuilder {
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

// Build creates a 'upgrade_policy_state' object using the configuration stored in the builder.
func (b *UpgradePolicyStateBuilder) Build() (object *UpgradePolicyState, err error) {
	object = new(UpgradePolicyState)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.description = b.description
	object.value = b.value
	return
}
