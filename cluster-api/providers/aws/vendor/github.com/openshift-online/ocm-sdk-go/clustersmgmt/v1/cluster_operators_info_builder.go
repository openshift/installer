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

// ClusterOperatorsInfoBuilder contains the data and logic needed to build 'cluster_operators_info' objects.
//
// Provides detailed information about the operators installed on the cluster.
type ClusterOperatorsInfoBuilder struct {
	bitmap_   uint32
	operators []*ClusterOperatorInfoBuilder
}

// NewClusterOperatorsInfo creates a new builder of 'cluster_operators_info' objects.
func NewClusterOperatorsInfo() *ClusterOperatorsInfoBuilder {
	return &ClusterOperatorsInfoBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterOperatorsInfoBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Operators sets the value of the 'operators' attribute to the given values.
func (b *ClusterOperatorsInfoBuilder) Operators(values ...*ClusterOperatorInfoBuilder) *ClusterOperatorsInfoBuilder {
	b.operators = make([]*ClusterOperatorInfoBuilder, len(values))
	copy(b.operators, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterOperatorsInfoBuilder) Copy(object *ClusterOperatorsInfo) *ClusterOperatorsInfoBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.operators != nil {
		b.operators = make([]*ClusterOperatorInfoBuilder, len(object.operators))
		for i, v := range object.operators {
			b.operators[i] = NewClusterOperatorInfo().Copy(v)
		}
	} else {
		b.operators = nil
	}
	return b
}

// Build creates a 'cluster_operators_info' object using the configuration stored in the builder.
func (b *ClusterOperatorsInfoBuilder) Build() (object *ClusterOperatorsInfo, err error) {
	object = new(ClusterOperatorsInfo)
	object.bitmap_ = b.bitmap_
	if b.operators != nil {
		object.operators = make([]*ClusterOperatorInfo, len(b.operators))
		for i, v := range b.operators {
			object.operators[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
