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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	time "time"
)

// ClusterUpgradeBuilder contains the data and logic needed to build 'cluster_upgrade' objects.
type ClusterUpgradeBuilder struct {
	bitmap_          uint32
	state            string
	updatedTimestamp time.Time
	version          string
	available        bool
}

// NewClusterUpgrade creates a new builder of 'cluster_upgrade' objects.
func NewClusterUpgrade() *ClusterUpgradeBuilder {
	return &ClusterUpgradeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterUpgradeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Available sets the value of the 'available' attribute to the given value.
func (b *ClusterUpgradeBuilder) Available(value bool) *ClusterUpgradeBuilder {
	b.available = value
	b.bitmap_ |= 1
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *ClusterUpgradeBuilder) State(value string) *ClusterUpgradeBuilder {
	b.state = value
	b.bitmap_ |= 2
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ClusterUpgradeBuilder) UpdatedTimestamp(value time.Time) *ClusterUpgradeBuilder {
	b.updatedTimestamp = value
	b.bitmap_ |= 4
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *ClusterUpgradeBuilder) Version(value string) *ClusterUpgradeBuilder {
	b.version = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterUpgradeBuilder) Copy(object *ClusterUpgrade) *ClusterUpgradeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.available = object.available
	b.state = object.state
	b.updatedTimestamp = object.updatedTimestamp
	b.version = object.version
	return b
}

// Build creates a 'cluster_upgrade' object using the configuration stored in the builder.
func (b *ClusterUpgradeBuilder) Build() (object *ClusterUpgrade, err error) {
	object = new(ClusterUpgrade)
	object.bitmap_ = b.bitmap_
	object.available = b.available
	object.state = b.state
	object.updatedTimestamp = b.updatedTimestamp
	object.version = b.version
	return
}
