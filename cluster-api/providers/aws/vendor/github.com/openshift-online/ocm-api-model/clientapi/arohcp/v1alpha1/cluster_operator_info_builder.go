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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	time "time"
)

type ClusterOperatorInfoBuilder struct {
	fieldSet_ []bool
	condition ClusterOperatorState
	name      string
	reason    string
	time      time.Time
	version   string
}

// NewClusterOperatorInfo creates a new builder of 'cluster_operator_info' objects.
func NewClusterOperatorInfo() *ClusterOperatorInfoBuilder {
	return &ClusterOperatorInfoBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterOperatorInfoBuilder) Empty() bool {
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

// Condition sets the value of the 'condition' attribute to the given value.
//
// Overall state of a cluster operator.
func (b *ClusterOperatorInfoBuilder) Condition(value ClusterOperatorState) *ClusterOperatorInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.condition = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Name(value string) *ClusterOperatorInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Reason(value string) *ClusterOperatorInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.reason = value
	b.fieldSet_[2] = true
	return b
}

// Time sets the value of the 'time' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Time(value time.Time) *ClusterOperatorInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.time = value
	b.fieldSet_[3] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Version(value string) *ClusterOperatorInfoBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.version = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterOperatorInfoBuilder) Copy(object *ClusterOperatorInfo) *ClusterOperatorInfoBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.condition = object.condition
	b.name = object.name
	b.reason = object.reason
	b.time = object.time
	b.version = object.version
	return b
}

// Build creates a 'cluster_operator_info' object using the configuration stored in the builder.
func (b *ClusterOperatorInfoBuilder) Build() (object *ClusterOperatorInfo, err error) {
	object = new(ClusterOperatorInfo)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.condition = b.condition
	object.name = b.name
	object.reason = b.reason
	object.time = b.time
	object.version = b.version
	return
}
