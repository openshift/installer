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

// Definition of a Status Board status.
type StatusBuilder struct {
	fieldSet_   []bool
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
	return &StatusBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *StatusBuilder) Link(value bool) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *StatusBuilder) ID(value string) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *StatusBuilder) HREF(value string) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StatusBuilder) Empty() bool {
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *StatusBuilder) CreatedAt(value time.Time) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *StatusBuilder) Metadata(value interface{}) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.metadata = value
	b.fieldSet_[4] = true
	return b
}

// Service sets the value of the 'service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *StatusBuilder) Service(value *ServiceBuilder) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.service = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// ServiceInfo sets the value of the 'service_info' attribute to the given value.
//
// Definition of a Status Board service info.
func (b *StatusBuilder) ServiceInfo(value *ServiceInfoBuilder) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.serviceInfo = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *StatusBuilder) Status(value string) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.status = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *StatusBuilder) UpdatedAt(value time.Time) *StatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StatusBuilder) Copy(object *Status) *StatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
