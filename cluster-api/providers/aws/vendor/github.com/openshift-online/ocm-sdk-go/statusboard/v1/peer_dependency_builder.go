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

// PeerDependencyBuilder contains the data and logic needed to build 'peer_dependency' objects.
//
// Definition of a Status Board peer dependency.
type PeerDependencyBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	metadata  interface{}
	name      string
	owners    []*OwnerBuilder
	services  []*ServiceBuilder
	updatedAt time.Time
}

// NewPeerDependency creates a new builder of 'peer_dependency' objects.
func NewPeerDependency() *PeerDependencyBuilder {
	return &PeerDependencyBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *PeerDependencyBuilder) Link(value bool) *PeerDependencyBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *PeerDependencyBuilder) ID(value string) *PeerDependencyBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *PeerDependencyBuilder) HREF(value string) *PeerDependencyBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PeerDependencyBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *PeerDependencyBuilder) CreatedAt(value time.Time) *PeerDependencyBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *PeerDependencyBuilder) Metadata(value interface{}) *PeerDependencyBuilder {
	b.metadata = value
	b.bitmap_ |= 16
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *PeerDependencyBuilder) Name(value string) *PeerDependencyBuilder {
	b.name = value
	b.bitmap_ |= 32
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *PeerDependencyBuilder) Owners(values ...*OwnerBuilder) *PeerDependencyBuilder {
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.bitmap_ |= 64
	return b
}

// Services sets the value of the 'services' attribute to the given values.
func (b *PeerDependencyBuilder) Services(values ...*ServiceBuilder) *PeerDependencyBuilder {
	b.services = make([]*ServiceBuilder, len(values))
	copy(b.services, values)
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *PeerDependencyBuilder) UpdatedAt(value time.Time) *PeerDependencyBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PeerDependencyBuilder) Copy(object *PeerDependency) *PeerDependencyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
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
	if object.services != nil {
		b.services = make([]*ServiceBuilder, len(object.services))
		for i, v := range object.services {
			b.services[i] = NewService().Copy(v)
		}
	} else {
		b.services = nil
	}
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'peer_dependency' object using the configuration stored in the builder.
func (b *PeerDependencyBuilder) Build() (object *PeerDependency, err error) {
	object = new(PeerDependency)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
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
	if b.services != nil {
		object.services = make([]*Service, len(b.services))
		for i, v := range b.services {
			object.services[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.updatedAt = b.updatedAt
	return
}
