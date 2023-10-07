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

// ServiceDependencyBuilder contains the data and logic needed to build 'service_dependency' objects.
//
// Definition of a Status Board service dependency.
type ServiceDependencyBuilder struct {
	bitmap_       uint32
	id            string
	href          string
	childService  *ServiceBuilder
	createdAt     time.Time
	metadata      interface{}
	name          string
	owners        []*OwnerBuilder
	parentService *ServiceBuilder
	type_         string
	updatedAt     time.Time
}

// NewServiceDependency creates a new builder of 'service_dependency' objects.
func NewServiceDependency() *ServiceDependencyBuilder {
	return &ServiceDependencyBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ServiceDependencyBuilder) Link(value bool) *ServiceDependencyBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ServiceDependencyBuilder) ID(value string) *ServiceDependencyBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ServiceDependencyBuilder) HREF(value string) *ServiceDependencyBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceDependencyBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ChildService sets the value of the 'child_service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *ServiceDependencyBuilder) ChildService(value *ServiceBuilder) *ServiceDependencyBuilder {
	b.childService = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ServiceDependencyBuilder) CreatedAt(value time.Time) *ServiceDependencyBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *ServiceDependencyBuilder) Metadata(value interface{}) *ServiceDependencyBuilder {
	b.metadata = value
	b.bitmap_ |= 32
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ServiceDependencyBuilder) Name(value string) *ServiceDependencyBuilder {
	b.name = value
	b.bitmap_ |= 64
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *ServiceDependencyBuilder) Owners(values ...*OwnerBuilder) *ServiceDependencyBuilder {
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.bitmap_ |= 128
	return b
}

// ParentService sets the value of the 'parent_service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *ServiceDependencyBuilder) ParentService(value *ServiceBuilder) *ServiceDependencyBuilder {
	b.parentService = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *ServiceDependencyBuilder) Type(value string) *ServiceDependencyBuilder {
	b.type_ = value
	b.bitmap_ |= 512
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ServiceDependencyBuilder) UpdatedAt(value time.Time) *ServiceDependencyBuilder {
	b.updatedAt = value
	b.bitmap_ |= 1024
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceDependencyBuilder) Copy(object *ServiceDependency) *ServiceDependencyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.childService != nil {
		b.childService = NewService().Copy(object.childService)
	} else {
		b.childService = nil
	}
	b.createdAt = object.createdAt
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
	if object.parentService != nil {
		b.parentService = NewService().Copy(object.parentService)
	} else {
		b.parentService = nil
	}
	b.type_ = object.type_
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'service_dependency' object using the configuration stored in the builder.
func (b *ServiceDependencyBuilder) Build() (object *ServiceDependency, err error) {
	object = new(ServiceDependency)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.childService != nil {
		object.childService, err = b.childService.Build()
		if err != nil {
			return
		}
	}
	object.createdAt = b.createdAt
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
	if b.parentService != nil {
		object.parentService, err = b.parentService.Build()
		if err != nil {
			return
		}
	}
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	return
}
