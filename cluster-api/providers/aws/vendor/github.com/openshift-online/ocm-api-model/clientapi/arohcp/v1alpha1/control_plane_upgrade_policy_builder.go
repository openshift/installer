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
type ControlPlaneUpgradePolicyBuilder struct {
	fieldSet_           []bool
	id                  string
	href                string
	clusterID           string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	state               *UpgradePolicyStateBuilder
	version             string
}

// NewControlPlaneUpgradePolicy creates a new builder of 'control_plane_upgrade_policy' objects.
func NewControlPlaneUpgradePolicy() *ControlPlaneUpgradePolicyBuilder {
	return &ControlPlaneUpgradePolicyBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ControlPlaneUpgradePolicyBuilder) Link(value bool) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ControlPlaneUpgradePolicyBuilder) ID(value string) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ControlPlaneUpgradePolicyBuilder) HREF(value string) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ControlPlaneUpgradePolicyBuilder) Empty() bool {
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
func (b *ControlPlaneUpgradePolicyBuilder) ClusterID(value string) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) CreationTimestamp(value time.Time) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.creationTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) LastUpdateTimestamp(value time.Time) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[5] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Representation of an upgrade policy state that that is set for a cluster.
func (b *ControlPlaneUpgradePolicyBuilder) State(value *UpgradePolicyStateBuilder) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.state = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *ControlPlaneUpgradePolicyBuilder) Version(value string) *ControlPlaneUpgradePolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.version = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ControlPlaneUpgradePolicyBuilder) Copy(object *ControlPlaneUpgradePolicy) *ControlPlaneUpgradePolicyBuilder {
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
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	if object.state != nil {
		b.state = NewUpgradePolicyState().Copy(object.state)
	} else {
		b.state = nil
	}
	b.version = object.version
	return b
}

// Build creates a 'control_plane_upgrade_policy' object using the configuration stored in the builder.
func (b *ControlPlaneUpgradePolicyBuilder) Build() (object *ControlPlaneUpgradePolicy, err error) {
	object = new(ControlPlaneUpgradePolicy)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterID = b.clusterID
	object.creationTimestamp = b.creationTimestamp
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	if b.state != nil {
		object.state, err = b.state.Build()
		if err != nil {
			return
		}
	}
	object.version = b.version
	return
}
