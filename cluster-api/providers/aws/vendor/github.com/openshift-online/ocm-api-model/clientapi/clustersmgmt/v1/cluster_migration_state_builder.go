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

// Representation of a cluster migration state.
type ClusterMigrationStateBuilder struct {
	fieldSet_   []bool
	description string
	value       ClusterMigrationStateValue
}

// NewClusterMigrationState creates a new builder of 'cluster_migration_state' objects.
func NewClusterMigrationState() *ClusterMigrationStateBuilder {
	return &ClusterMigrationStateBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMigrationStateBuilder) Empty() bool {
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

// Description sets the value of the 'description' attribute to the given value.
func (b *ClusterMigrationStateBuilder) Description(value string) *ClusterMigrationStateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.description = value
	b.fieldSet_[0] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
//
// The state of the cluster migration.
func (b *ClusterMigrationStateBuilder) Value(value ClusterMigrationStateValue) *ClusterMigrationStateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.value = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMigrationStateBuilder) Copy(object *ClusterMigrationState) *ClusterMigrationStateBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.description = object.description
	b.value = object.value
	return b
}

// Build creates a 'cluster_migration_state' object using the configuration stored in the builder.
func (b *ClusterMigrationStateBuilder) Build() (object *ClusterMigrationState, err error) {
	object = new(ClusterMigrationState)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.description = b.description
	object.value = b.value
	return
}
