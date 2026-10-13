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

type CloudAccountBuilder struct {
	fieldSet_       []bool
	cloudAccountID  string
	cloudProviderID string
	contracts       []*ContractBuilder
}

// NewCloudAccount creates a new builder of 'cloud_account' objects.
func NewCloudAccount() *CloudAccountBuilder {
	return &CloudAccountBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudAccountBuilder) Empty() bool {
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

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *CloudAccountBuilder) CloudAccountID(value string) *CloudAccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.cloudAccountID = value
	b.fieldSet_[0] = true
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *CloudAccountBuilder) CloudProviderID(value string) *CloudAccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.cloudProviderID = value
	b.fieldSet_[1] = true
	return b
}

// Contracts sets the value of the 'contracts' attribute to the given values.
func (b *CloudAccountBuilder) Contracts(values ...*ContractBuilder) *CloudAccountBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.contracts = make([]*ContractBuilder, len(values))
	copy(b.contracts, values)
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudAccountBuilder) Copy(object *CloudAccount) *CloudAccountBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.cloudAccountID = object.cloudAccountID
	b.cloudProviderID = object.cloudProviderID
	if object.contracts != nil {
		b.contracts = make([]*ContractBuilder, len(object.contracts))
		for i, v := range object.contracts {
			b.contracts[i] = NewContract().Copy(v)
		}
	} else {
		b.contracts = nil
	}
	return b
}

// Build creates a 'cloud_account' object using the configuration stored in the builder.
func (b *CloudAccountBuilder) Build() (object *CloudAccount, err error) {
	object = new(CloudAccount)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.cloudAccountID = b.cloudAccountID
	object.cloudProviderID = b.cloudProviderID
	if b.contracts != nil {
		object.contracts = make([]*Contract, len(b.contracts))
		for i, v := range b.contracts {
			object.contracts[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
