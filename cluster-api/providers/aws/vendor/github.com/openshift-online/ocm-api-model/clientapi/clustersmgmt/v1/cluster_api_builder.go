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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Information about the API of a cluster.
type ClusterAPIBuilder struct {
	fieldSet_       []bool
	cidrBlockAccess *CIDRBlockAccessBuilder
	url             string
	listening       ListeningMethod
}

// NewClusterAPI creates a new builder of 'cluster_API' objects.
func NewClusterAPI() *ClusterAPIBuilder {
	return &ClusterAPIBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAPIBuilder) Empty() bool {
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

// CIDRBlockAccess sets the value of the 'CIDR_block_access' attribute to the given value.
//
// Describes the CIDR Block access policy to the Kubernetes API server.
// Currently, only supported for ARO-HCP based clusters.
// The default policy mode is "allow_all" that is, all access is allowed.
func (b *ClusterAPIBuilder) CIDRBlockAccess(value *CIDRBlockAccessBuilder) *ClusterAPIBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.cidrBlockAccess = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// URL sets the value of the 'URL' attribute to the given value.
func (b *ClusterAPIBuilder) URL(value string) *ClusterAPIBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.url = value
	b.fieldSet_[1] = true
	return b
}

// Listening sets the value of the 'listening' attribute to the given value.
//
// Cluster components listening method.
func (b *ClusterAPIBuilder) Listening(value ListeningMethod) *ClusterAPIBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.listening = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAPIBuilder) Copy(object *ClusterAPI) *ClusterAPIBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.cidrBlockAccess != nil {
		b.cidrBlockAccess = NewCIDRBlockAccess().Copy(object.cidrBlockAccess)
	} else {
		b.cidrBlockAccess = nil
	}
	b.url = object.url
	b.listening = object.listening
	return b
}

// Build creates a 'cluster_API' object using the configuration stored in the builder.
func (b *ClusterAPIBuilder) Build() (object *ClusterAPI, err error) {
	object = new(ClusterAPI)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.cidrBlockAccess != nil {
		object.cidrBlockAccess, err = b.cidrBlockAccess.Build()
		if err != nil {
			return
		}
	}
	object.url = b.url
	object.listening = b.listening
	return
}
