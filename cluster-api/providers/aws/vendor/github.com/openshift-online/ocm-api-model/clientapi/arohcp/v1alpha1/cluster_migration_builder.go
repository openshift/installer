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

// Representation of a cluster migration.
type ClusterMigrationBuilder struct {
	fieldSet_         []bool
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
	return &ClusterMigrationBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterMigrationBuilder) Link(value bool) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ClusterMigrationBuilder) ID(value string) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ClusterMigrationBuilder) HREF(value string) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterMigrationBuilder) Empty() bool {
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
func (b *ClusterMigrationBuilder) ClusterID(value string) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ClusterMigrationBuilder) CreationTimestamp(value time.Time) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.creationTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// SdnToOvn sets the value of the 'sdn_to_ovn' attribute to the given value.
//
// Details for `SdnToOvn` cluster migrations.
func (b *ClusterMigrationBuilder) SdnToOvn(value *SdnToOvnClusterMigrationBuilder) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.sdnToOvn = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of a cluster migration state.
func (b *ClusterMigrationBuilder) State(value *ClusterMigrationStateBuilder) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.state = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
//
// Type of cluster migration.
func (b *ClusterMigrationBuilder) Type(value ClusterMigrationType) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.type_ = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ClusterMigrationBuilder) UpdatedTimestamp(value time.Time) *ClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedTimestamp = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterMigrationBuilder) Copy(object *ClusterMigration) *ClusterMigrationBuilder {
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
