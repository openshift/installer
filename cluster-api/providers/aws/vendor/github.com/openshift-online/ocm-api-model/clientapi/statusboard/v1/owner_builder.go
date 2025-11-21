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

// Definition of a Status Board owner.
type OwnerBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	email     string
	username  string
}

// NewOwner creates a new builder of 'owner' objects.
func NewOwner() *OwnerBuilder {
	return &OwnerBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *OwnerBuilder) Link(value bool) *OwnerBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *OwnerBuilder) ID(value string) *OwnerBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *OwnerBuilder) HREF(value string) *OwnerBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OwnerBuilder) Empty() bool {
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

// Email sets the value of the 'email' attribute to the given value.
func (b *OwnerBuilder) Email(value string) *OwnerBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.email = value
	b.fieldSet_[3] = true
	return b
}

// Username sets the value of the 'username' attribute to the given value.
func (b *OwnerBuilder) Username(value string) *OwnerBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.username = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OwnerBuilder) Copy(object *Owner) *OwnerBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.email = object.email
	b.username = object.username
	return b
}

// Build creates a 'owner' object using the configuration stored in the builder.
func (b *OwnerBuilder) Build() (object *Owner, err error) {
	object = new(Owner)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.email = b.email
	object.username = b.username
	return
}
