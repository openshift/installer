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

// ContractBuilder contains the data and logic needed to build 'contract' objects.
type ContractBuilder struct {
	bitmap_    uint32
	dimensions []*ContractDimensionBuilder
	endDate    time.Time
	startDate  time.Time
}

// NewContract creates a new builder of 'contract' objects.
func NewContract() *ContractBuilder {
	return &ContractBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ContractBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Dimensions sets the value of the 'dimensions' attribute to the given values.
func (b *ContractBuilder) Dimensions(values ...*ContractDimensionBuilder) *ContractBuilder {
	b.dimensions = make([]*ContractDimensionBuilder, len(values))
	copy(b.dimensions, values)
	b.bitmap_ |= 1
	return b
}

// EndDate sets the value of the 'end_date' attribute to the given value.
func (b *ContractBuilder) EndDate(value time.Time) *ContractBuilder {
	b.endDate = value
	b.bitmap_ |= 2
	return b
}

// StartDate sets the value of the 'start_date' attribute to the given value.
func (b *ContractBuilder) StartDate(value time.Time) *ContractBuilder {
	b.startDate = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ContractBuilder) Copy(object *Contract) *ContractBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
