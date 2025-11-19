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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

import (
	time "time"
)

// Represents data about a running Managed Service.
type ManagedServiceBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	addon        *StatefulObjectBuilder
	cluster      *ClusterBuilder
	createdAt    time.Time
	expiredAt    time.Time
	parameters   []*ServiceParameterBuilder
	resources    []*StatefulObjectBuilder
	service      string
	serviceState string
	updatedAt    time.Time
}

// NewManagedService creates a new builder of 'managed_service' objects.
func NewManagedService() *ManagedServiceBuilder {
	return &ManagedServiceBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ManagedServiceBuilder) Link(value bool) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ManagedServiceBuilder) ID(value string) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ManagedServiceBuilder) HREF(value string) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagedServiceBuilder) Empty() bool {
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

// Addon sets the value of the 'addon' attribute to the given value.
func (b *ManagedServiceBuilder) Addon(value *StatefulObjectBuilder) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.addon = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// Cluster sets the value of the 'cluster' attribute to the given value.
//
// This represents the parameters needed by Managed Service to create a cluster.
func (b *ManagedServiceBuilder) Cluster(value *ClusterBuilder) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.cluster = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ManagedServiceBuilder) CreatedAt(value time.Time) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.createdAt = value
	b.fieldSet_[5] = true
	return b
}

// ExpiredAt sets the value of the 'expired_at' attribute to the given value.
func (b *ManagedServiceBuilder) ExpiredAt(value time.Time) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.expiredAt = value
	b.fieldSet_[6] = true
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *ManagedServiceBuilder) Parameters(values ...*ServiceParameterBuilder) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.parameters = make([]*ServiceParameterBuilder, len(values))
	copy(b.parameters, values)
	b.fieldSet_[7] = true
	return b
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *ManagedServiceBuilder) Resources(values ...*StatefulObjectBuilder) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.resources = make([]*StatefulObjectBuilder, len(values))
	copy(b.resources, values)
	b.fieldSet_[8] = true
	return b
}

// Service sets the value of the 'service' attribute to the given value.
func (b *ManagedServiceBuilder) Service(value string) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.service = value
	b.fieldSet_[9] = true
	return b
}

// ServiceState sets the value of the 'service_state' attribute to the given value.
func (b *ManagedServiceBuilder) ServiceState(value string) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.serviceState = value
	b.fieldSet_[10] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ManagedServiceBuilder) UpdatedAt(value time.Time) *ManagedServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.updatedAt = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagedServiceBuilder) Copy(object *ManagedService) *ManagedServiceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.addon != nil {
		b.addon = NewStatefulObject().Copy(object.addon)
	} else {
		b.addon = nil
	}
	if object.cluster != nil {
		b.cluster = NewCluster().Copy(object.cluster)
	} else {
		b.cluster = nil
	}
	b.createdAt = object.createdAt
	b.expiredAt = object.expiredAt
	if object.parameters != nil {
		b.parameters = make([]*ServiceParameterBuilder, len(object.parameters))
		for i, v := range object.parameters {
			b.parameters[i] = NewServiceParameter().Copy(v)
		}
	} else {
		b.parameters = nil
	}
	if object.resources != nil {
		b.resources = make([]*StatefulObjectBuilder, len(object.resources))
		for i, v := range object.resources {
			b.resources[i] = NewStatefulObject().Copy(v)
		}
	} else {
		b.resources = nil
	}
	b.service = object.service
	b.serviceState = object.serviceState
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'managed_service' object using the configuration stored in the builder.
func (b *ManagedServiceBuilder) Build() (object *ManagedService, err error) {
	object = new(ManagedService)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.addon != nil {
		object.addon, err = b.addon.Build()
		if err != nil {
			return
		}
	}
	if b.cluster != nil {
		object.cluster, err = b.cluster.Build()
		if err != nil {
			return
		}
	}
	object.createdAt = b.createdAt
	object.expiredAt = b.expiredAt
	if b.parameters != nil {
		object.parameters = make([]*ServiceParameter, len(b.parameters))
		for i, v := range b.parameters {
			object.parameters[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.resources != nil {
		object.resources = make([]*StatefulObject, len(b.resources))
		for i, v := range b.resources {
			object.resources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.service = b.service
	object.serviceState = b.serviceState
	object.updatedAt = b.updatedAt
	return
}
