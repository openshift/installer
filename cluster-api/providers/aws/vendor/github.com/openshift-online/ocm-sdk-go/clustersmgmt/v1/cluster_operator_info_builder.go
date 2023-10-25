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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	time "time"
)

// ClusterOperatorInfoBuilder contains the data and logic needed to build 'cluster_operator_info' objects.
type ClusterOperatorInfoBuilder struct {
	bitmap_   uint32
	condition ClusterOperatorState
	name      string
	reason    string
	time      time.Time
	version   string
}

// NewClusterOperatorInfo creates a new builder of 'cluster_operator_info' objects.
func NewClusterOperatorInfo() *ClusterOperatorInfoBuilder {
	return &ClusterOperatorInfoBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterOperatorInfoBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Condition sets the value of the 'condition' attribute to the given value.
//
// Overall state of a cluster operator.
func (b *ClusterOperatorInfoBuilder) Condition(value ClusterOperatorState) *ClusterOperatorInfoBuilder {
	b.condition = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Name(value string) *ClusterOperatorInfoBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Reason(value string) *ClusterOperatorInfoBuilder {
	b.reason = value
	b.bitmap_ |= 4
	return b
}

// Time sets the value of the 'time' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Time(value time.Time) *ClusterOperatorInfoBuilder {
	b.time = value
	b.bitmap_ |= 8
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *ClusterOperatorInfoBuilder) Version(value string) *ClusterOperatorInfoBuilder {
	b.version = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterOperatorInfoBuilder) Copy(object *ClusterOperatorInfo) *ClusterOperatorInfoBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.condition = b.condition
	object.name = b.name
	object.reason = b.reason
	object.time = b.time
	object.version = b.version
	return
}
