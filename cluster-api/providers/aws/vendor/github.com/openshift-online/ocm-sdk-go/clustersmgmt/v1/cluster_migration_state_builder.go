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

// ClusterMigrationStateBuilder contains the data and logic needed to build 'cluster_migration_state' objects.
//
// Representation of a cluster migration state.
type ClusterMigrationStateBuilder struct {
	bitmap_     uint32
	description string
	value       ClusterMigrationStateValue
}

// NewClusterMigrationState creates a new builder of 'cluster_migration_state' objects.
func NewClusterMigrationState() *ClusterMigrationStateBuilder {
	return &ClusterMigrationStateBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMigrationStateBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Description sets the value of the 'description' attribute to the given value.
func (b *ClusterMigrationStateBuilder) Description(value string) *ClusterMigrationStateBuilder {
	b.description = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
//
// The state of the cluster migration.
func (b *ClusterMigrationStateBuilder) Value(value ClusterMigrationStateValue) *ClusterMigrationStateBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMigrationStateBuilder) Copy(object *ClusterMigrationState) *ClusterMigrationStateBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.description = object.description
	b.value = object.value
	return b
}

// Build creates a 'cluster_migration_state' object using the configuration stored in the builder.
func (b *ClusterMigrationStateBuilder) Build() (object *ClusterMigrationState, err error) {
	object = new(ClusterMigrationState)
	object.bitmap_ = b.bitmap_
	object.description = b.description
	object.value = b.value
	return
}
