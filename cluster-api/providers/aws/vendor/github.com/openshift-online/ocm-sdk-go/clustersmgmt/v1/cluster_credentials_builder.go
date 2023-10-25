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

// ClusterCredentialsBuilder contains the data and logic needed to build 'cluster_credentials' objects.
//
// Credentials of the a cluster.
type ClusterCredentialsBuilder struct {
	bitmap_    uint32
	id         string
	href       string
	kubeconfig string
}

// NewClusterCredentials creates a new builder of 'cluster_credentials' objects.
func NewClusterCredentials() *ClusterCredentialsBuilder {
	return &ClusterCredentialsBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterCredentialsBuilder) Link(value bool) *ClusterCredentialsBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ClusterCredentialsBuilder) ID(value string) *ClusterCredentialsBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ClusterCredentialsBuilder) HREF(value string) *ClusterCredentialsBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterCredentialsBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Kubeconfig sets the value of the 'kubeconfig' attribute to the given value.
func (b *ClusterCredentialsBuilder) Kubeconfig(value string) *ClusterCredentialsBuilder {
	b.kubeconfig = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterCredentialsBuilder) Copy(object *ClusterCredentials) *ClusterCredentialsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.kubeconfig = object.kubeconfig
	return b
}

// Build creates a 'cluster_credentials' object using the configuration stored in the builder.
func (b *ClusterCredentialsBuilder) Build() (object *ClusterCredentials, err error) {
	object = new(ClusterCredentials)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.kubeconfig = b.kubeconfig
	return
}
