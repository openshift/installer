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

type LabelBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	accountID      string
	createdAt      time.Time
	key            string
	managedBy      string
	organizationID string
	subscriptionID string
	type_          string
	updatedAt      time.Time
	value          string
	internal       bool
}

// NewLabel creates a new builder of 'label' objects.
func NewLabel() *LabelBuilder {
	return &LabelBuilder{
		fieldSet_: make([]bool, 13),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *LabelBuilder) Link(value bool) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *LabelBuilder) ID(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *LabelBuilder) HREF(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LabelBuilder) Empty() bool {
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

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *LabelBuilder) AccountID(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.accountID = value
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *LabelBuilder) CreatedAt(value time.Time) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// Internal sets the value of the 'internal' attribute to the given value.
func (b *LabelBuilder) Internal(value bool) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.internal = value
	b.fieldSet_[5] = true
	return b
}

// Key sets the value of the 'key' attribute to the given value.
func (b *LabelBuilder) Key(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.key = value
	b.fieldSet_[6] = true
	return b
}

// ManagedBy sets the value of the 'managed_by' attribute to the given value.
func (b *LabelBuilder) ManagedBy(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.managedBy = value
	b.fieldSet_[7] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *LabelBuilder) OrganizationID(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.organizationID = value
	b.fieldSet_[8] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *LabelBuilder) SubscriptionID(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.subscriptionID = value
	b.fieldSet_[9] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *LabelBuilder) Type(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.type_ = value
	b.fieldSet_[10] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *LabelBuilder) UpdatedAt(value time.Time) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.updatedAt = value
	b.fieldSet_[11] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *LabelBuilder) Value(value string) *LabelBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.value = value
	b.fieldSet_[12] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LabelBuilder) Copy(object *Label) *LabelBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.accountID = object.accountID
	b.createdAt = object.createdAt
	b.internal = object.internal
	b.key = object.key
	b.managedBy = object.managedBy
	b.organizationID = object.organizationID
	b.subscriptionID = object.subscriptionID
	b.type_ = object.type_
	b.updatedAt = object.updatedAt
	b.value = object.value
	return b
}

// Build creates a 'label' object using the configuration stored in the builder.
func (b *LabelBuilder) Build() (object *Label, err error) {
	object = new(Label)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountID = b.accountID
	object.createdAt = b.createdAt
	object.internal = b.internal
	object.key = b.key
	object.managedBy = b.managedBy
	object.organizationID = b.organizationID
	object.subscriptionID = b.subscriptionID
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	object.value = b.value
	return
}
