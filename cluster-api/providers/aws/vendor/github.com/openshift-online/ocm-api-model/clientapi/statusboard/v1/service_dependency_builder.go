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

// Definition of a Status Board service dependency.
type ServiceDependencyBuilder struct {
	fieldSet_     []bool
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
	return &ServiceDependencyBuilder{
		fieldSet_: make([]bool, 11),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ServiceDependencyBuilder) Link(value bool) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ServiceDependencyBuilder) ID(value string) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ServiceDependencyBuilder) HREF(value string) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ServiceDependencyBuilder) Empty() bool {
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

// ChildService sets the value of the 'child_service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *ServiceDependencyBuilder) ChildService(value *ServiceBuilder) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.childService = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ServiceDependencyBuilder) CreatedAt(value time.Time) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *ServiceDependencyBuilder) Metadata(value interface{}) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.metadata = value
	b.fieldSet_[5] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ServiceDependencyBuilder) Name(value string) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.name = value
	b.fieldSet_[6] = true
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *ServiceDependencyBuilder) Owners(values ...*OwnerBuilder) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.fieldSet_[7] = true
	return b
}

// ParentService sets the value of the 'parent_service' attribute to the given value.
//
// Definition of a Status Board Service.
func (b *ServiceDependencyBuilder) ParentService(value *ServiceBuilder) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.parentService = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *ServiceDependencyBuilder) Type(value string) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.type_ = value
	b.fieldSet_[9] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ServiceDependencyBuilder) UpdatedAt(value time.Time) *ServiceDependencyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 11)
	}
	b.updatedAt = value
	b.fieldSet_[10] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ServiceDependencyBuilder) Copy(object *ServiceDependency) *ServiceDependencyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
