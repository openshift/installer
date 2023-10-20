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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// PrivateLinkPrincipalsBuilder contains the data and logic needed to build 'private_link_principals' objects.
//
// Contains a list of principals for the Private Link.
type PrivateLinkPrincipalsBuilder struct {
	bitmap_    uint32
	id         string
	href       string
	principals []*PrivateLinkPrincipalBuilder
}

// NewPrivateLinkPrincipals creates a new builder of 'private_link_principals' objects.
func NewPrivateLinkPrincipals() *PrivateLinkPrincipalsBuilder {
	return &PrivateLinkPrincipalsBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *PrivateLinkPrincipalsBuilder) Link(value bool) *PrivateLinkPrincipalsBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *PrivateLinkPrincipalsBuilder) ID(value string) *PrivateLinkPrincipalsBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *PrivateLinkPrincipalsBuilder) HREF(value string) *PrivateLinkPrincipalsBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *PrivateLinkPrincipalsBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Principals sets the value of the 'principals' attribute to the given values.
func (b *PrivateLinkPrincipalsBuilder) Principals(values ...*PrivateLinkPrincipalBuilder) *PrivateLinkPrincipalsBuilder {
	b.principals = make([]*PrivateLinkPrincipalBuilder, len(values))
	copy(b.principals, values)
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *PrivateLinkPrincipalsBuilder) Copy(object *PrivateLinkPrincipals) *PrivateLinkPrincipalsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
