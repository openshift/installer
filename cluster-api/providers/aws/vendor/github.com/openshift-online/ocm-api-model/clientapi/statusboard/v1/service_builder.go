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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

import (
	time "time"
)

// Definition of a Status Board Service.
type ServiceBuilder struct {
	fieldSet_       []bool
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
	return &ServiceBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ServiceBuilder) Link(value bool) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ServiceBuilder) ID(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ServiceBuilder) HREF(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceBuilder) Empty() bool {
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

// Application sets the value of the 'application' attribute to the given value.
//
// Definition of a Status Board application.
func (b *ServiceBuilder) Application(value *ApplicationBuilder) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.application = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ServiceBuilder) CreatedAt(value time.Time) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// CurrentStatus sets the value of the 'current_status' attribute to the given value.
func (b *ServiceBuilder) CurrentStatus(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.currentStatus = value
	b.fieldSet_[5] = true
	return b
}

// Fullname sets the value of the 'fullname' attribute to the given value.
func (b *ServiceBuilder) Fullname(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fullname = value
	b.fieldSet_[6] = true
	return b
}

// LastPingAt sets the value of the 'last_ping_at' attribute to the given value.
func (b *ServiceBuilder) LastPingAt(value time.Time) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.lastPingAt = value
	b.fieldSet_[7] = true
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *ServiceBuilder) Metadata(value interface{}) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.metadata = value
	b.fieldSet_[8] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ServiceBuilder) Name(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.name = value
	b.fieldSet_[9] = true
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *ServiceBuilder) Owners(values ...*OwnerBuilder) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.fieldSet_[10] = true
	return b
}

// Private sets the value of the 'private' attribute to the given value.
func (b *ServiceBuilder) Private(value bool) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.private = value
	b.fieldSet_[11] = true
	return b
}

// ServiceEndpoint sets the value of the 'service_endpoint' attribute to the given value.
func (b *ServiceBuilder) ServiceEndpoint(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.serviceEndpoint = value
	b.fieldSet_[12] = true
	return b
}

// StatusType sets the value of the 'status_type' attribute to the given value.
func (b *ServiceBuilder) StatusType(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.statusType = value
	b.fieldSet_[13] = true
	return b
}

// StatusUpdatedAt sets the value of the 'status_updated_at' attribute to the given value.
func (b *ServiceBuilder) StatusUpdatedAt(value time.Time) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.statusUpdatedAt = value
	b.fieldSet_[14] = true
	return b
}

// Token sets the value of the 'token' attribute to the given value.
func (b *ServiceBuilder) Token(value string) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.token = value
	b.fieldSet_[15] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ServiceBuilder) UpdatedAt(value time.Time) *ServiceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.updatedAt = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceBuilder) Copy(object *Service) *ServiceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
