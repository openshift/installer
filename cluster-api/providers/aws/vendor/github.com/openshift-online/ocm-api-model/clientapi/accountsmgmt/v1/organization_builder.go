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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	time "time"
)

type OrganizationBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	capabilities []*CapabilityBuilder
	createdAt    time.Time
	ebsAccountID string
	externalID   string
	labels       []*LabelBuilder
	name         string
	updatedAt    time.Time
}

// NewOrganization creates a new builder of 'organization' objects.
func NewOrganization() *OrganizationBuilder {
	return &OrganizationBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *OrganizationBuilder) Link(value bool) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *OrganizationBuilder) ID(value string) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *OrganizationBuilder) HREF(value string) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OrganizationBuilder) Empty() bool {
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

// Capabilities sets the value of the 'capabilities' attribute to the given values.
func (b *OrganizationBuilder) Capabilities(values ...*CapabilityBuilder) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.capabilities = make([]*CapabilityBuilder, len(values))
	copy(b.capabilities, values)
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *OrganizationBuilder) CreatedAt(value time.Time) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// EbsAccountID sets the value of the 'ebs_account_ID' attribute to the given value.
func (b *OrganizationBuilder) EbsAccountID(value string) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.ebsAccountID = value
	b.fieldSet_[5] = true
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *OrganizationBuilder) ExternalID(value string) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.externalID = value
	b.fieldSet_[6] = true
	return b
}

// Labels sets the value of the 'labels' attribute to the given values.
func (b *OrganizationBuilder) Labels(values ...*LabelBuilder) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.labels = make([]*LabelBuilder, len(values))
	copy(b.labels, values)
	b.fieldSet_[7] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *OrganizationBuilder) Name(value string) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.name = value
	b.fieldSet_[8] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *OrganizationBuilder) UpdatedAt(value time.Time) *OrganizationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OrganizationBuilder) Copy(object *Organization) *OrganizationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.capabilities != nil {
		b.capabilities = make([]*CapabilityBuilder, len(object.capabilities))
		for i, v := range object.capabilities {
			b.capabilities[i] = NewCapability().Copy(v)
		}
	} else {
		b.capabilities = nil
	}
	b.createdAt = object.createdAt
	b.ebsAccountID = object.ebsAccountID
	b.externalID = object.externalID
	if object.labels != nil {
		b.labels = make([]*LabelBuilder, len(object.labels))
		for i, v := range object.labels {
			b.labels[i] = NewLabel().Copy(v)
		}
	} else {
		b.labels = nil
	}
	b.name = object.name
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'organization' object using the configuration stored in the builder.
func (b *OrganizationBuilder) Build() (object *Organization, err error) {
	object = new(Organization)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.capabilities != nil {
		object.capabilities = make([]*Capability, len(b.capabilities))
		for i, v := range b.capabilities {
			object.capabilities[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.createdAt = b.createdAt
	object.ebsAccountID = b.ebsAccountID
	object.externalID = b.externalID
	if b.labels != nil {
		object.labels = make([]*Label, len(b.labels))
		for i, v := range b.labels {
			object.labels[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.name = b.name
	object.updatedAt = b.updatedAt
	return
}
