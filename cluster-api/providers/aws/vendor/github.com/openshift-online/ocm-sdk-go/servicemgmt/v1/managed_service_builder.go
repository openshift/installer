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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	time "time"
)

// ManagedServiceBuilder contains the data and logic needed to build 'managed_service' objects.
//
// Represents data about a running Managed Service.
type ManagedServiceBuilder struct {
	bitmap_      uint32
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
	return &ManagedServiceBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ManagedServiceBuilder) Link(value bool) *ManagedServiceBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ManagedServiceBuilder) ID(value string) *ManagedServiceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ManagedServiceBuilder) HREF(value string) *ManagedServiceBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagedServiceBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Addon sets the value of the 'addon' attribute to the given value.
func (b *ManagedServiceBuilder) Addon(value *StatefulObjectBuilder) *ManagedServiceBuilder {
	b.addon = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Cluster sets the value of the 'cluster' attribute to the given value.
//
// This represents the parameters needed by Managed Service to create a cluster.
func (b *ManagedServiceBuilder) Cluster(value *ClusterBuilder) *ManagedServiceBuilder {
	b.cluster = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ManagedServiceBuilder) CreatedAt(value time.Time) *ManagedServiceBuilder {
	b.createdAt = value
	b.bitmap_ |= 32
	return b
}

// ExpiredAt sets the value of the 'expired_at' attribute to the given value.
func (b *ManagedServiceBuilder) ExpiredAt(value time.Time) *ManagedServiceBuilder {
	b.expiredAt = value
	b.bitmap_ |= 64
	return b
}

// Parameters sets the value of the 'parameters' attribute to the given values.
func (b *ManagedServiceBuilder) Parameters(values ...*ServiceParameterBuilder) *ManagedServiceBuilder {
	b.parameters = make([]*ServiceParameterBuilder, len(values))
	copy(b.parameters, values)
	b.bitmap_ |= 128
	return b
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *ManagedServiceBuilder) Resources(values ...*StatefulObjectBuilder) *ManagedServiceBuilder {
	b.resources = make([]*StatefulObjectBuilder, len(values))
	copy(b.resources, values)
	b.bitmap_ |= 256
	return b
}

// Service sets the value of the 'service' attribute to the given value.
func (b *ManagedServiceBuilder) Service(value string) *ManagedServiceBuilder {
	b.service = value
	b.bitmap_ |= 512
	return b
}

// ServiceState sets the value of the 'service_state' attribute to the given value.
func (b *ManagedServiceBuilder) ServiceState(value string) *ManagedServiceBuilder {
	b.serviceState = value
	b.bitmap_ |= 1024
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ManagedServiceBuilder) UpdatedAt(value time.Time) *ManagedServiceBuilder {
	b.updatedAt = value
	b.bitmap_ |= 2048
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagedServiceBuilder) Copy(object *ManagedService) *ManagedServiceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
