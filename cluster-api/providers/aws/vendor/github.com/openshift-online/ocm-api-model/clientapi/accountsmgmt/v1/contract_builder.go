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

type ContractBuilder struct {
	fieldSet_  []bool
	dimensions []*ContractDimensionBuilder
	endDate    time.Time
	startDate  time.Time
}

// NewContract creates a new builder of 'contract' objects.
func NewContract() *ContractBuilder {
	return &ContractBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ContractBuilder) Empty() bool {
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

// Dimensions sets the value of the 'dimensions' attribute to the given values.
func (b *ContractBuilder) Dimensions(values ...*ContractDimensionBuilder) *ContractBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.dimensions = make([]*ContractDimensionBuilder, len(values))
	copy(b.dimensions, values)
	b.fieldSet_[0] = true
	return b
}

// EndDate sets the value of the 'end_date' attribute to the given value.
func (b *ContractBuilder) EndDate(value time.Time) *ContractBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.endDate = value
	b.fieldSet_[1] = true
	return b
}

// StartDate sets the value of the 'start_date' attribute to the given value.
func (b *ContractBuilder) StartDate(value time.Time) *ContractBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.startDate = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ContractBuilder) Copy(object *Contract) *ContractBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.dimensions != nil {
		b.dimensions = make([]*ContractDimensionBuilder, len(object.dimensions))
		for i, v := range object.dimensions {
			b.dimensions[i] = NewContractDimension().Copy(v)
		}
	} else {
		b.dimensions = nil
	}
	b.endDate = object.endDate
	b.startDate = object.startDate
	return b
}

// Build creates a 'contract' object using the configuration stored in the builder.
func (b *ContractBuilder) Build() (object *Contract, err error) {
	object = new(Contract)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.dimensions != nil {
		object.dimensions = make([]*ContractDimension, len(b.dimensions))
		for i, v := range b.dimensions {
			object.dimensions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.endDate = b.endDate
	object.startDate = b.startDate
	return
}
