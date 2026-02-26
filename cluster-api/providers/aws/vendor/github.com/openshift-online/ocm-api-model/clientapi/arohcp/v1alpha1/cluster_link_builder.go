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

// Definition of a cluster link.
type ClusterLinkBuilder struct {
	fieldSet_ []bool
	href      string
	id        string
}

// NewClusterLink creates a new builder of 'cluster_link' objects.
func NewClusterLink() *ClusterLinkBuilder {
	return &ClusterLinkBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterLinkBuilder) Empty() bool {
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

// HREF sets the value of the 'HREF' attribute to the given value.
func (b *ClusterLinkBuilder) HREF(value string) *ClusterLinkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.href = value
	b.fieldSet_[0] = true
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *ClusterLinkBuilder) ID(value string) *ClusterLinkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterLinkBuilder) Copy(object *ClusterLink) *ClusterLinkBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.href = object.href
	b.id = object.id
	return b
}

// Build creates a 'cluster_link' object using the configuration stored in the builder.
func (b *ClusterLinkBuilder) Build() (object *ClusterLink, err error) {
	object = new(ClusterLink)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.href = b.href
	object.id = b.id
	return
}
