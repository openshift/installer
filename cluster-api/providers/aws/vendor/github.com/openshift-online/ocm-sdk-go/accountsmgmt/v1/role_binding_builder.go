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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	time "time"
)

// RoleBindingBuilder contains the data and logic needed to build 'role_binding' objects.
type RoleBindingBuilder struct {
	bitmap_        uint32
	id             string
	href           string
	account        *AccountBuilder
	accountID      string
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
	return &RoleBindingBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *RoleBindingBuilder) Link(value bool) *RoleBindingBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *RoleBindingBuilder) ID(value string) *RoleBindingBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *RoleBindingBuilder) HREF(value string) *RoleBindingBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RoleBindingBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Account sets the value of the 'account' attribute to the given value.
func (b *RoleBindingBuilder) Account(value *AccountBuilder) *RoleBindingBuilder {
	b.account = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *RoleBindingBuilder) AccountID(value string) *RoleBindingBuilder {
	b.accountID = value
	b.bitmap_ |= 16
	return b
}

// ConfigManaged sets the value of the 'config_managed' attribute to the given value.
func (b *RoleBindingBuilder) ConfigManaged(value bool) *RoleBindingBuilder {
	b.configManaged = value
	b.bitmap_ |= 32
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *RoleBindingBuilder) CreatedAt(value time.Time) *RoleBindingBuilder {
	b.createdAt = value
	b.bitmap_ |= 64
	return b
}

// ManagedBy sets the value of the 'managed_by' attribute to the given value.
func (b *RoleBindingBuilder) ManagedBy(value string) *RoleBindingBuilder {
	b.managedBy = value
	b.bitmap_ |= 128
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
func (b *RoleBindingBuilder) Organization(value *OrganizationBuilder) *RoleBindingBuilder {
	b.organization = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *RoleBindingBuilder) OrganizationID(value string) *RoleBindingBuilder {
	b.organizationID = value
	b.bitmap_ |= 512
	return b
}

// Role sets the value of the 'role' attribute to the given value.
func (b *RoleBindingBuilder) Role(value *RoleBuilder) *RoleBindingBuilder {
	b.role = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// RoleID sets the value of the 'role_ID' attribute to the given value.
func (b *RoleBindingBuilder) RoleID(value string) *RoleBindingBuilder {
	b.roleID = value
	b.bitmap_ |= 2048
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
func (b *RoleBindingBuilder) Subscription(value *SubscriptionBuilder) *RoleBindingBuilder {
	b.subscription = value
	if value != nil {
		b.bitmap_ |= 4096
	} else {
		b.bitmap_ &^= 4096
	}
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *RoleBindingBuilder) SubscriptionID(value string) *RoleBindingBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 8192
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *RoleBindingBuilder) Type(value string) *RoleBindingBuilder {
	b.type_ = value
	b.bitmap_ |= 16384
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *RoleBindingBuilder) UpdatedAt(value time.Time) *RoleBindingBuilder {
	b.updatedAt = value
	b.bitmap_ |= 32768
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RoleBindingBuilder) Copy(object *RoleBinding) *RoleBindingBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.account != nil {
		b.account = NewAccount().Copy(object.account)
	} else {
		b.account = nil
	}
	b.accountID = object.accountID
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
	object.bitmap_ = b.bitmap_
	if b.account != nil {
		object.account, err = b.account.Build()
		if err != nil {
			return
		}
	}
	object.accountID = b.accountID
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
