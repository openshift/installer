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

type QuotaAuthorizationRequestBuilder struct {
	fieldSet_        []bool
	accountUsername  string
	availabilityZone string
	displayName      string
	productID        string
	productCategory  string
	quotaVersion     string
	resources        []*ReservedResourceBuilder
	reserve          bool
}

// NewQuotaAuthorizationRequest creates a new builder of 'quota_authorization_request' objects.
func NewQuotaAuthorizationRequest() *QuotaAuthorizationRequestBuilder {
	return &QuotaAuthorizationRequestBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *QuotaAuthorizationRequestBuilder) Empty() bool {
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

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) AccountUsername(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.accountUsername = value
	b.fieldSet_[0] = true
	return b
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) AvailabilityZone(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.availabilityZone = value
	b.fieldSet_[1] = true
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) DisplayName(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.displayName = value
	b.fieldSet_[2] = true
	return b
}

// ProductID sets the value of the 'product_ID' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) ProductID(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.productID = value
	b.fieldSet_[3] = true
	return b
}

// ProductCategory sets the value of the 'product_category' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) ProductCategory(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.productCategory = value
	b.fieldSet_[4] = true
	return b
}

// QuotaVersion sets the value of the 'quota_version' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) QuotaVersion(value string) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.quotaVersion = value
	b.fieldSet_[5] = true
	return b
}

// Reserve sets the value of the 'reserve' attribute to the given value.
func (b *QuotaAuthorizationRequestBuilder) Reserve(value bool) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.reserve = value
	b.fieldSet_[6] = true
	return b
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *QuotaAuthorizationRequestBuilder) Resources(values ...*ReservedResourceBuilder) *QuotaAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.resources = make([]*ReservedResourceBuilder, len(values))
	copy(b.resources, values)
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *QuotaAuthorizationRequestBuilder) Copy(object *QuotaAuthorizationRequest) *QuotaAuthorizationRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.accountUsername = object.accountUsername
	b.availabilityZone = object.availabilityZone
	b.displayName = object.displayName
	b.productID = object.productID
	b.productCategory = object.productCategory
	b.quotaVersion = object.quotaVersion
	b.reserve = object.reserve
	if object.resources != nil {
		b.resources = make([]*ReservedResourceBuilder, len(object.resources))
		for i, v := range object.resources {
			b.resources[i] = NewReservedResource().Copy(v)
		}
	} else {
		b.resources = nil
	}
	return b
}

// Build creates a 'quota_authorization_request' object using the configuration stored in the builder.
func (b *QuotaAuthorizationRequestBuilder) Build() (object *QuotaAuthorizationRequest, err error) {
	object = new(QuotaAuthorizationRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountUsername = b.accountUsername
	object.availabilityZone = b.availabilityZone
	object.displayName = b.displayName
	object.productID = b.productID
	object.productCategory = b.productCategory
	object.quotaVersion = b.quotaVersion
	object.reserve = b.reserve
	if b.resources != nil {
		object.resources = make([]*ReservedResource, len(b.resources))
		for i, v := range b.resources {
			object.resources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
