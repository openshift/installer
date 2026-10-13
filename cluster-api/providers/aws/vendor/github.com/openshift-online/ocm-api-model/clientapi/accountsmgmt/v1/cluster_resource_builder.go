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

type ClusterResourceBuilder struct {
	fieldSet_        []bool
	total            *ValueUnitBuilder
	updatedTimestamp time.Time
	used             *ValueUnitBuilder
}

// NewClusterResource creates a new builder of 'cluster_resource' objects.
func NewClusterResource() *ClusterResourceBuilder {
	return &ClusterResourceBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterResourceBuilder) Empty() bool {
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

// Total sets the value of the 'total' attribute to the given value.
func (b *ClusterResourceBuilder) Total(value *ValueUnitBuilder) *ClusterResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.total = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ClusterResourceBuilder) UpdatedTimestamp(value time.Time) *ClusterResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.updatedTimestamp = value
	b.fieldSet_[1] = true
	return b
}

// Used sets the value of the 'used' attribute to the given value.
func (b *ClusterResourceBuilder) Used(value *ValueUnitBuilder) *ClusterResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.used = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterResourceBuilder) Copy(object *ClusterResource) *ClusterResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.total != nil {
		b.total = NewValueUnit().Copy(object.total)
	} else {
		b.total = nil
	}
	b.updatedTimestamp = object.updatedTimestamp
	if object.used != nil {
		b.used = NewValueUnit().Copy(object.used)
	} else {
		b.used = nil
	}
	return b
}

// Build creates a 'cluster_resource' object using the configuration stored in the builder.
func (b *ClusterResourceBuilder) Build() (object *ClusterResource, err error) {
	object = new(ClusterResource)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.total != nil {
		object.total, err = b.total.Build()
		if err != nil {
			return
		}
	}
	object.updatedTimestamp = b.updatedTimestamp
	if b.used != nil {
		object.used, err = b.used.Build()
		if err != nil {
			return
		}
	}
	return
}
