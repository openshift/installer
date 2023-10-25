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

// StatusBuilder contains the data and logic needed to build 'status' objects.
//
// Definition of a Status Board status.
type StatusBuilder struct {
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

// NewStatus creates a new builder of 'status' objects.
func NewStatus() *StatusBuilder {
	return &StatusBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *StatusBuilder) Link(value bool) *StatusBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *StatusBuilder) ID(value string) *StatusBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *StatusBuilder) HREF(value string) *StatusBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StatusBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *StatusBuilder) CreatedAt(value time.Time) *StatusBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *StatusBuilder) Metadata(value interface{}) *StatusBuilder {
	b.metadata = value
	b.bitmap_ |= 16
	return b
}

// Service sets the value of the 'service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *StatusBuilder) Service(value *ServiceBuilder) *StatusBuilder {
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
func (b *StatusBuilder) ServiceInfo(value *ServiceInfoBuilder) *StatusBuilder {
	b.serviceInfo = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *StatusBuilder) Status(value string) *StatusBuilder {
	b.status = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *StatusBuilder) UpdatedAt(value time.Time) *StatusBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StatusBuilder) Copy(object *Status) *StatusBuilder {
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

// Build creates a 'status' object using the configuration stored in the builder.
func (b *StatusBuilder) Build() (object *Status, err error) {
	object = new(Status)
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
