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

type RoleBindingBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	account        *AccountBuilder
	accountID      string
	accountGroup   *AccountGroupBuilder
	accountGroupID string
	createdAt      time.Time
	managedBy      string
	organization   *OrganizationBuilder
	organizationID string
	role           *RoleBuilder
	roleID         string
	subscription   *SubscriptionBuilder
	subscriptionID string
	type_          string
	updatedAt      time.Time
	configManaged  bool
}

// NewRoleBinding creates a new builder of 'role_binding' objects.
func NewRoleBinding() *RoleBindingBuilder {
	return &RoleBindingBuilder{
		fieldSet_: make([]bool, 18),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *RoleBindingBuilder) Link(value bool) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *RoleBindingBuilder) ID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *RoleBindingBuilder) HREF(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RoleBindingBuilder) Empty() bool {
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

// Account sets the value of the 'account' attribute to the given value.
func (b *RoleBindingBuilder) Account(value *AccountBuilder) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.account = value
	if value != nil {
		b.fieldSet_[3] = true
	} else {
		b.fieldSet_[3] = false
	}
	return b
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *RoleBindingBuilder) AccountID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.accountID = value
	b.fieldSet_[4] = true
	return b
}

// AccountGroup sets the value of the 'account_group' attribute to the given value.
func (b *RoleBindingBuilder) AccountGroup(value *AccountGroupBuilder) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.accountGroup = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// AccountGroupID sets the value of the 'account_group_ID' attribute to the given value.
func (b *RoleBindingBuilder) AccountGroupID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.accountGroupID = value
	b.fieldSet_[6] = true
	return b
}

// ConfigManaged sets the value of the 'config_managed' attribute to the given value.
func (b *RoleBindingBuilder) ConfigManaged(value bool) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.configManaged = value
	b.fieldSet_[7] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RoleBindingBuilder) CreatedAt(value time.Time) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.createdAt = value
	b.fieldSet_[8] = true
	return b
}

// ManagedBy sets the value of the 'managed_by' attribute to the given value.
func (b *RoleBindingBuilder) ManagedBy(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.managedBy = value
	b.fieldSet_[9] = true
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
func (b *RoleBindingBuilder) Organization(value *OrganizationBuilder) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.organization = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *RoleBindingBuilder) OrganizationID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.organizationID = value
	b.fieldSet_[11] = true
	return b
}

// Role sets the value of the 'role' attribute to the given value.
func (b *RoleBindingBuilder) Role(value *RoleBuilder) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.role = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// RoleID sets the value of the 'role_ID' attribute to the given value.
func (b *RoleBindingBuilder) RoleID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.roleID = value
	b.fieldSet_[13] = true
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
func (b *RoleBindingBuilder) Subscription(value *SubscriptionBuilder) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.subscription = value
	if value != nil {
		b.fieldSet_[14] = true
	} else {
		b.fieldSet_[14] = false
	}
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *RoleBindingBuilder) SubscriptionID(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.subscriptionID = value
	b.fieldSet_[15] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RoleBindingBuilder) Type(value string) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.type_ = value
	b.fieldSet_[16] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RoleBindingBuilder) UpdatedAt(value time.Time) *RoleBindingBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 18)
	}
	b.updatedAt = value
	b.fieldSet_[17] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RoleBindingBuilder) Copy(object *RoleBinding) *RoleBindingBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.account != nil {
		b.account = NewAccount().Copy(object.account)
	} else {
		b.account = nil
	}
	b.accountID = object.accountID
	if object.accountGroup != nil {
		b.accountGroup = NewAccountGroup().Copy(object.accountGroup)
	} else {
		b.accountGroup = nil
	}
	b.accountGroupID = object.accountGroupID
	b.configManaged = object.configManaged
	b.createdAt = object.createdAt
	b.managedBy = object.managedBy
	if object.organization != nil {
		b.organization = NewOrganization().Copy(object.organization)
	} else {
		b.organization = nil
	}
	b.organizationID = object.organizationID
	if object.role != nil {
		b.role = NewRole().Copy(object.role)
	} else {
		b.role = nil
	}
	b.roleID = object.roleID
	if object.subscription != nil {
		b.subscription = NewSubscription().Copy(object.subscription)
	} else {
		b.subscription = nil
	}
	b.subscriptionID = object.subscriptionID
	b.type_ = object.type_
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'role_binding' object using the configuration stored in the builder.
func (b *RoleBindingBuilder) Build() (object *RoleBinding, err error) {
	object = new(RoleBinding)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.account != nil {
		object.account, err = b.account.Build()
		if err != nil {
			return
		}
	}
	object.accountID = b.accountID
	if b.accountGroup != nil {
		object.accountGroup, err = b.accountGroup.Build()
		if err != nil {
			return
		}
	}
	object.accountGroupID = b.accountGroupID
	object.configManaged = b.configManaged
	object.createdAt = b.createdAt
	object.managedBy = b.managedBy
	if b.organization != nil {
		object.organization, err = b.organization.Build()
		if err != nil {
			return
		}
	}
	object.organizationID = b.organizationID
	if b.role != nil {
		object.role, err = b.role.Build()
		if err != nil {
			return
		}
	}
	object.roleID = b.roleID
	if b.subscription != nil {
		object.subscription, err = b.subscription.Build()
		if err != nil {
			return
		}
	}
	object.subscriptionID = b.subscriptionID
	object.type_ = b.type_
	object.updatedAt = b.updatedAt
	return
}
