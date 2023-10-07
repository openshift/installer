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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	time "time"
)

// ServiceBuilder contains the data and logic needed to build 'service' objects.
//
// Definition of a Status Board Service.
type ServiceBuilder struct {
	bitmap_         uint32
	id              string
	href            string
	application     *ApplicationBuilder
	createdAt       time.Time
	currentStatus   string
	fullname        string
	lastPingAt      time.Time
	metadata        interface{}
	name            string
	owners          []*OwnerBuilder
	serviceEndpoint string
	statusType      string
	statusUpdatedAt time.Time
	token           string
	updatedAt       time.Time
	private         bool
}

// NewService creates a new builder of 'service' objects.
func NewService() *ServiceBuilder {
	return &ServiceBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ServiceBuilder) Link(value bool) *ServiceBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ServiceBuilder) ID(value string) *ServiceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ServiceBuilder) HREF(value string) *ServiceBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Application sets the value of the 'application' attribute to the given value.
//
// Definition of a Status Board application.
func (b *ServiceBuilder) Application(value *ApplicationBuilder) *ServiceBuilder {
	b.application = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ServiceBuilder) CreatedAt(value time.Time) *ServiceBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// CurrentStatus sets the value of the 'current_status' attribute to the given value.
func (b *ServiceBuilder) CurrentStatus(value string) *ServiceBuilder {
	b.currentStatus = value
	b.bitmap_ |= 32
	return b
}

// Fullname sets the value of the 'fullname' attribute to the given value.
func (b *ServiceBuilder) Fullname(value string) *ServiceBuilder {
	b.fullname = value
	b.bitmap_ |= 64
	return b
}

// LastPingAt sets the value of the 'last_ping_at' attribute to the given value.
func (b *ServiceBuilder) LastPingAt(value time.Time) *ServiceBuilder {
	b.lastPingAt = value
	b.bitmap_ |= 128
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *ServiceBuilder) Metadata(value interface{}) *ServiceBuilder {
	b.metadata = value
	b.bitmap_ |= 256
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ServiceBuilder) Name(value string) *ServiceBuilder {
	b.name = value
	b.bitmap_ |= 512
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *ServiceBuilder) Owners(values ...*OwnerBuilder) *ServiceBuilder {
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.bitmap_ |= 1024
	return b
}

// Private sets the value of the 'private' attribute to the given value.
func (b *ServiceBuilder) Private(value bool) *ServiceBuilder {
	b.private = value
	b.bitmap_ |= 2048
	return b
}

// ServiceEndpoint sets the value of the 'service_endpoint' attribute to the given value.
func (b *ServiceBuilder) ServiceEndpoint(value string) *ServiceBuilder {
	b.serviceEndpoint = value
	b.bitmap_ |= 4096
	return b
}

// StatusType sets the value of the 'status_type' attribute to the given value.
func (b *ServiceBuilder) StatusType(value string) *ServiceBuilder {
	b.statusType = value
	b.bitmap_ |= 8192
	return b
}

// StatusUpdatedAt sets the value of the 'status_updated_at' attribute to the given value.
func (b *ServiceBuilder) StatusUpdatedAt(value time.Time) *ServiceBuilder {
	b.statusUpdatedAt = value
	b.bitmap_ |= 16384
	return b
}

// Token sets the value of the 'token' attribute to the given value.
func (b *ServiceBuilder) Token(value string) *ServiceBuilder {
	b.token = value
	b.bitmap_ |= 32768
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ServiceBuilder) UpdatedAt(value time.Time) *ServiceBuilder {
	b.updatedAt = value
	b.bitmap_ |= 65536
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceBuilder) Copy(object *Service) *ServiceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.application != nil {
		b.application = NewApplication().Copy(object.application)
	} else {
		b.application = nil
	}
	b.createdAt = object.createdAt
	b.currentStatus = object.currentStatus
	b.fullname = object.fullname
	b.lastPingAt = object.lastPingAt
	b.metadata = object.metadata
	b.name = object.name
	if object.owners != nil {
		b.owners = make([]*OwnerBuilder, len(object.owners))
		for i, v := range object.owners {
			b.owners[i] = NewOwner().Copy(v)
		}
	} else {
		b.owners = nil
	}
	b.private = object.private
	b.serviceEndpoint = object.serviceEndpoint
	b.statusType = object.statusType
	b.statusUpdatedAt = object.statusUpdatedAt
	b.token = object.token
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'service' object using the configuration stored in the builder.
func (b *ServiceBuilder) Build() (object *Service, err error) {
	object = new(Service)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.application != nil {
		object.application, err = b.application.Build()
		if err != nil {
			return
		}
	}
	object.createdAt = b.createdAt
	object.currentStatus = b.currentStatus
	object.fullname = b.fullname
	object.lastPingAt = b.lastPingAt
	object.metadata = b.metadata
	object.name = b.name
	if b.owners != nil {
		object.owners = make([]*Owner, len(b.owners))
		for i, v := range b.owners {
			object.owners[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.private = b.private
	object.serviceEndpoint = b.serviceEndpoint
	object.statusType = b.statusType
	object.statusUpdatedAt = b.statusUpdatedAt
	object.token = b.token
	object.updatedAt = b.updatedAt
	return
}
