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

type AccountGroupAssignmentBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	accountID      string
	accountGroup   *AccountGroupBuilder
	accountGroupID string
	createdAt      time.Time
	managedBy      AccountGroupAssignmentManagedBy
	updatedAt      time.Time
}

// NewAccountGroupAssignment creates a new builder of 'account_group_assignment' objects.
func NewAccountGroupAssignment() *AccountGroupAssignmentBuilder {
	return &AccountGroupAssignmentBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AccountGroupAssignmentBuilder) Link(value bool) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AccountGroupAssignmentBuilder) ID(value string) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AccountGroupAssignmentBuilder) HREF(value string) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccountGroupAssignmentBuilder) Empty() bool {
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
func (b *AccountGroupAssignmentBuilder) AccountID(value string) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.accountID = value
	b.fieldSet_[3] = true
	return b
}

// AccountGroup sets the value of the 'account_group' attribute to the given value.
func (b *AccountGroupAssignmentBuilder) AccountGroup(value *AccountGroupBuilder) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.accountGroup = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// AccountGroupID sets the value of the 'account_group_ID' attribute to the given value.
func (b *AccountGroupAssignmentBuilder) AccountGroupID(value string) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.accountGroupID = value
	b.fieldSet_[5] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AccountGroupAssignmentBuilder) CreatedAt(value time.Time) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.createdAt = value
	b.fieldSet_[6] = true
	return b
}

// ManagedBy sets the value of the 'managed_by' attribute to the given value.
func (b *AccountGroupAssignmentBuilder) ManagedBy(value AccountGroupAssignmentManagedBy) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.managedBy = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AccountGroupAssignmentBuilder) UpdatedAt(value time.Time) *AccountGroupAssignmentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccountGroupAssignmentBuilder) Copy(object *AccountGroupAssignment) *AccountGroupAssignmentBuilder {
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
	if object.accountGroup != nil {
		b.accountGroup = NewAccountGroup().Copy(object.accountGroup)
	} else {
		b.accountGroup = nil
	}
	b.accountGroupID = object.accountGroupID
	b.createdAt = object.createdAt
	b.managedBy = object.managedBy
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'account_group_assignment' object using the configuration stored in the builder.
func (b *AccountGroupAssignmentBuilder) Build() (object *AccountGroupAssignment, err error) {
	object = new(AccountGroupAssignment)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountID = b.accountID
	if b.accountGroup != nil {
		object.accountGroup, err = b.accountGroup.Build()
		if err != nil {
			return
		}
	}
	object.accountGroupID = b.accountGroupID
	object.createdAt = b.createdAt
	object.managedBy = b.managedBy
	object.updatedAt = b.updatedAt
	return
}
