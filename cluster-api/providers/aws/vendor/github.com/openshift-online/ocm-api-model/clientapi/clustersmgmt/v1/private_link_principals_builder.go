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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Contains a list of principals for the Private Link.
type PrivateLinkPrincipalsBuilder struct {
	fieldSet_  []bool
	id         string
	href       string
	principals []*PrivateLinkPrincipalBuilder
}

// NewPrivateLinkPrincipals creates a new builder of 'private_link_principals' objects.
func NewPrivateLinkPrincipals() *PrivateLinkPrincipalsBuilder {
	return &PrivateLinkPrincipalsBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *PrivateLinkPrincipalsBuilder) Link(value bool) *PrivateLinkPrincipalsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *PrivateLinkPrincipalsBuilder) ID(value string) *PrivateLinkPrincipalsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *PrivateLinkPrincipalsBuilder) HREF(value string) *PrivateLinkPrincipalsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PrivateLinkPrincipalsBuilder) Empty() bool {
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

// Principals sets the value of the 'principals' attribute to the given values.
func (b *PrivateLinkPrincipalsBuilder) Principals(values ...*PrivateLinkPrincipalBuilder) *PrivateLinkPrincipalsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.principals = make([]*PrivateLinkPrincipalBuilder, len(values))
	copy(b.principals, values)
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PrivateLinkPrincipalsBuilder) Copy(object *PrivateLinkPrincipals) *PrivateLinkPrincipalsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.principals != nil {
		b.principals = make([]*PrivateLinkPrincipalBuilder, len(object.principals))
		for i, v := range object.principals {
			b.principals[i] = NewPrivateLinkPrincipal().Copy(v)
		}
	} else {
		b.principals = nil
	}
	return b
}

// Build creates a 'private_link_principals' object using the configuration stored in the builder.
func (b *PrivateLinkPrincipalsBuilder) Build() (object *PrivateLinkPrincipals, err error) {
	object = new(PrivateLinkPrincipals)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.principals != nil {
		object.principals = make([]*PrivateLinkPrincipal, len(b.principals))
		for i, v := range b.principals {
			object.principals[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
