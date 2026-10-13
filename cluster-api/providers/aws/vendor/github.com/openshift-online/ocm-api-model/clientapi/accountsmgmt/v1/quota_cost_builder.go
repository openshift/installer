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

type QuotaCostBuilder struct {
	fieldSet_        []bool
	allowed          int
	cloudAccounts    []*CloudAccountBuilder
	consumed         int
	organizationID   string
	quotaID          string
	relatedResources []*RelatedResourceBuilder
	version          string
}

// NewQuotaCost creates a new builder of 'quota_cost' objects.
func NewQuotaCost() *QuotaCostBuilder {
	return &QuotaCostBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *QuotaCostBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *QuotaCostBuilder) Allowed(value int) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.allowed = value
	b.fieldSet_[0] = true
	return b
}

// CloudAccounts sets the value of the 'cloud_accounts' attribute to the given values.
func (b *QuotaCostBuilder) CloudAccounts(values ...*CloudAccountBuilder) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.cloudAccounts = make([]*CloudAccountBuilder, len(values))
	copy(b.cloudAccounts, values)
	b.fieldSet_[1] = true
	return b
}

// Consumed sets the value of the 'consumed' attribute to the given value.
func (b *QuotaCostBuilder) Consumed(value int) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.consumed = value
	b.fieldSet_[2] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *QuotaCostBuilder) OrganizationID(value string) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.organizationID = value
	b.fieldSet_[3] = true
	return b
}

// QuotaID sets the value of the 'quota_ID' attribute to the given value.
func (b *QuotaCostBuilder) QuotaID(value string) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.quotaID = value
	b.fieldSet_[4] = true
	return b
}

// RelatedResources sets the value of the 'related_resources' attribute to the given values.
func (b *QuotaCostBuilder) RelatedResources(values ...*RelatedResourceBuilder) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.relatedResources = make([]*RelatedResourceBuilder, len(values))
	copy(b.relatedResources, values)
	b.fieldSet_[5] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *QuotaCostBuilder) Version(value string) *QuotaCostBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.version = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *QuotaCostBuilder) Copy(object *QuotaCost) *QuotaCostBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.allowed = object.allowed
	if object.cloudAccounts != nil {
		b.cloudAccounts = make([]*CloudAccountBuilder, len(object.cloudAccounts))
		for i, v := range object.cloudAccounts {
			b.cloudAccounts[i] = NewCloudAccount().Copy(v)
		}
	} else {
		b.cloudAccounts = nil
	}
	b.consumed = object.consumed
	b.organizationID = object.organizationID
	b.quotaID = object.quotaID
	if object.relatedResources != nil {
		b.relatedResources = make([]*RelatedResourceBuilder, len(object.relatedResources))
		for i, v := range object.relatedResources {
			b.relatedResources[i] = NewRelatedResource().Copy(v)
		}
	} else {
		b.relatedResources = nil
	}
	b.version = object.version
	return b
}

// Build creates a 'quota_cost' object using the configuration stored in the builder.
func (b *QuotaCostBuilder) Build() (object *QuotaCost, err error) {
	object = new(QuotaCost)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.allowed = b.allowed
	if b.cloudAccounts != nil {
		object.cloudAccounts = make([]*CloudAccount, len(b.cloudAccounts))
		for i, v := range b.cloudAccounts {
			object.cloudAccounts[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.consumed = b.consumed
	object.organizationID = b.organizationID
	object.quotaID = b.quotaID
	if b.relatedResources != nil {
		object.relatedResources = make([]*RelatedResource, len(b.relatedResources))
		for i, v := range b.relatedResources {
			object.relatedResources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.version = b.version
	return
}
