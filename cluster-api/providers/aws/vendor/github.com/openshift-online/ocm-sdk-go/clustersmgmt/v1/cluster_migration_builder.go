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

// ClusterMigrationBuilder contains the data and logic needed to build 'cluster_migration' objects.
//
// Representation of a cluster migration.
type ClusterMigrationBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	clusterID         string
	creationTimestamp time.Time
	sdnToOvn          *SdnToOvnClusterMigrationBuilder
	state             *ClusterMigrationStateBuilder
	type_             ClusterMigrationType
	updatedTimestamp  time.Time
}

// NewClusterMigration creates a new builder of 'cluster_migration' objects.
func NewClusterMigration() *ClusterMigrationBuilder {
	return &ClusterMigrationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterMigrationBuilder) Link(value bool) *ClusterMigrationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ClusterMigrationBuilder) ID(value string) *ClusterMigrationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ClusterMigrationBuilder) HREF(value string) *ClusterMigrationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMigrationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterMigrationBuilder) ClusterID(value string) *ClusterMigrationBuilder {
	b.clusterID = value
	b.bitmap_ |= 8
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ClusterMigrationBuilder) CreationTimestamp(value time.Time) *ClusterMigrationBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 16
	return b
}

// SdnToOvn sets the value of the 'sdn_to_ovn' attribute to the given value.
//
// Details for `SdnToOvn` cluster migrations.
func (b *ClusterMigrationBuilder) SdnToOvn(value *SdnToOvnClusterMigrationBuilder) *ClusterMigrationBuilder {
	b.sdnToOvn = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of a cluster migration state.
func (b *ClusterMigrationBuilder) State(value *ClusterMigrationStateBuilder) *ClusterMigrationBuilder {
	b.state = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Type of cluster migration.
func (b *ClusterMigrationBuilder) Type(value ClusterMigrationType) *ClusterMigrationBuilder {
	b.type_ = value
	b.bitmap_ |= 128
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ClusterMigrationBuilder) UpdatedTimestamp(value time.Time) *ClusterMigrationBuilder {
	b.updatedTimestamp = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMigrationBuilder) Copy(object *ClusterMigration) *ClusterMigrationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.clusterID = object.clusterID
	b.creationTimestamp = object.creationTimestamp
	if object.sdnToOvn != nil {
		b.sdnToOvn = NewSdnToOvnClusterMigration().Copy(object.sdnToOvn)
	} else {
		b.sdnToOvn = nil
	}
	if object.state != nil {
		b.state = NewClusterMigrationState().Copy(object.state)
	} else {
		b.state = nil
	}
	b.type_ = object.type_
	b.updatedTimestamp = object.updatedTimestamp
	return b
}

// Build creates a 'cluster_migration' object using the configuration stored in the builder.
func (b *ClusterMigrationBuilder) Build() (object *ClusterMigration, err error) {
	object = new(ClusterMigration)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.clusterID = b.clusterID
	object.creationTimestamp = b.creationTimestamp
	if b.sdnToOvn != nil {
		object.sdnToOvn, err = b.sdnToOvn.Build()
		if err != nil {
			return
		}
	}
	if b.state != nil {
		object.state, err = b.state.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	object.updatedTimestamp = b.updatedTimestamp
	return
}
