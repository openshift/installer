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

// StatusUpdateBuilder contains the data and logic needed to build 'status_update' objects.
type StatusUpdateBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	metadata    interface{}
	service     *ServiceBuilder
	serviceInfo *ServiceInfoBuilder
	status      string
	updatedAt   time.Time
}

// NewStatusUpdate creates a new builder of 'status_update' objects.
func NewStatusUpdate() *StatusUpdateBuilder {
	return &StatusUpdateBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *StatusUpdateBuilder) Link(value bool) *StatusUpdateBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *StatusUpdateBuilder) ID(value string) *StatusUpdateBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *StatusUpdateBuilder) HREF(value string) *StatusUpdateBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StatusUpdateBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *StatusUpdateBuilder) CreatedAt(value time.Time) *StatusUpdateBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *StatusUpdateBuilder) Metadata(value interface{}) *StatusUpdateBuilder {
	b.metadata = value
	b.bitmap_ |= 16
	return b
}

// Service sets the value of the 'service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *StatusUpdateBuilder) Service(value *ServiceBuilder) *StatusUpdateBuilder {
	b.service = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// ServiceInfo sets the value of the 'service_info' attribute to the given value.
//
// Definition of a Status Board service info.
func (b *StatusUpdateBuilder) ServiceInfo(value *ServiceInfoBuilder) *StatusUpdateBuilder {
	b.serviceInfo = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *StatusUpdateBuilder) Status(value string) *StatusUpdateBuilder {
	b.status = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *StatusUpdateBuilder) UpdatedAt(value time.Time) *StatusUpdateBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StatusUpdateBuilder) Copy(object *StatusUpdate) *StatusUpdateBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.metadata = object.metadata
	if object.service != nil {
		b.service = NewService().Copy(object.service)
	} else {
		b.service = nil
	}
	if object.serviceInfo != nil {
		b.serviceInfo = NewServiceInfo().Copy(object.serviceInfo)
	} else {
		b.serviceInfo = nil
	}
	b.status = object.status
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'status_update' object using the configuration stored in the builder.
func (b *StatusUpdateBuilder) Build() (object *StatusUpdate, err error) {
	object = new(StatusUpdate)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.metadata = b.metadata
	if b.service != nil {
		object.service, err = b.service.Build()
		if err != nil {
			return
		}
	}
	if b.serviceInfo != nil {
		object.serviceInfo, err = b.serviceInfo.Build()
		if err != nil {
			return
		}
	}
	object.status = b.status
	object.updatedAt = b.updatedAt
	return
}
