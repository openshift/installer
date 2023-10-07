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

// ClusterAuthorizationRequestBuilder contains the data and logic needed to build 'cluster_authorization_request' objects.
type ClusterAuthorizationRequestBuilder struct {
	bitmap_           uint32
	accountUsername   string
	availabilityZone  string
	cloudAccountID    string
	cloudProviderID   string
	clusterID         string
	displayName       string
	externalClusterID string
	productID         string
	productCategory   string
	quotaVersion      string
	resources         []*ReservedResourceBuilder
	byoc              bool
	disconnected      bool
	managed           bool
	reserve           bool
}

// NewClusterAuthorizationRequest creates a new builder of 'cluster_authorization_request' objects.
func NewClusterAuthorizationRequest() *ClusterAuthorizationRequestBuilder {
	return &ClusterAuthorizationRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAuthorizationRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) BYOC(value bool) *ClusterAuthorizationRequestBuilder {
	b.byoc = value
	b.bitmap_ |= 1
	return b
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) AccountUsername(value string) *ClusterAuthorizationRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 2
	return b
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) AvailabilityZone(value string) *ClusterAuthorizationRequestBuilder {
	b.availabilityZone = value
	b.bitmap_ |= 4
	return b
}

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) CloudAccountID(value string) *ClusterAuthorizationRequestBuilder {
	b.cloudAccountID = value
	b.bitmap_ |= 8
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) CloudProviderID(value string) *ClusterAuthorizationRequestBuilder {
	b.cloudProviderID = value
	b.bitmap_ |= 16
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ClusterID(value string) *ClusterAuthorizationRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 32
	return b
}

// Disconnected sets the value of the 'disconnected' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Disconnected(value bool) *ClusterAuthorizationRequestBuilder {
	b.disconnected = value
	b.bitmap_ |= 64
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) DisplayName(value string) *ClusterAuthorizationRequestBuilder {
	b.displayName = value
	b.bitmap_ |= 128
	return b
}

// ExternalClusterID sets the value of the 'external_cluster_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ExternalClusterID(value string) *ClusterAuthorizationRequestBuilder {
	b.externalClusterID = value
	b.bitmap_ |= 256
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Managed(value bool) *ClusterAuthorizationRequestBuilder {
	b.managed = value
	b.bitmap_ |= 512
	return b
}

// ProductID sets the value of the 'product_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ProductID(value string) *ClusterAuthorizationRequestBuilder {
	b.productID = value
	b.bitmap_ |= 1024
	return b
}

// ProductCategory sets the value of the 'product_category' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ProductCategory(value string) *ClusterAuthorizationRequestBuilder {
	b.productCategory = value
	b.bitmap_ |= 2048
	return b
}

// QuotaVersion sets the value of the 'quota_version' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) QuotaVersion(value string) *ClusterAuthorizationRequestBuilder {
	b.quotaVersion = value
	b.bitmap_ |= 4096
	return b
}

// Reserve sets the value of the 'reserve' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Reserve(value bool) *ClusterAuthorizationRequestBuilder {
	b.reserve = value
	b.bitmap_ |= 8192
	return b
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *ClusterAuthorizationRequestBuilder) Resources(values ...*ReservedResourceBuilder) *ClusterAuthorizationRequestBuilder {
	b.resources = make([]*ReservedResourceBuilder, len(values))
	copy(b.resources, values)
	b.bitmap_ |= 16384
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAuthorizationRequestBuilder) Copy(object *ClusterAuthorizationRequest) *ClusterAuthorizationRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.byoc = object.byoc
	b.accountUsername = object.accountUsername
	b.availabilityZone = object.availabilityZone
	b.cloudAccountID = object.cloudAccountID
	b.cloudProviderID = object.cloudProviderID
	b.clusterID = object.clusterID
	b.disconnected = object.disconnected
	b.displayName = object.displayName
	b.externalClusterID = object.externalClusterID
	b.managed = object.managed
	b.productID = object.productID
	b.productCategory = object.productCategory
	b.quotaVersion = object.quotaVersion
	b.reserve = object.reserve
	if object.resources != nil {
		b.resources = make([]*ReservedResourceBuilder, len(object.resources))
		for i, v := range object.resources {
			b.resources[i] = NewReservedResource().Copy(v)
		}
	} else {
		b.resources = nil
	}
	return b
}

// Build creates a 'cluster_authorization_request' object using the configuration stored in the builder.
func (b *ClusterAuthorizationRequestBuilder) Build() (object *ClusterAuthorizationRequest, err error) {
	object = new(ClusterAuthorizationRequest)
	object.bitmap_ = b.bitmap_
	object.byoc = b.byoc
	object.accountUsername = b.accountUsername
	object.availabilityZone = b.availabilityZone
	object.cloudAccountID = b.cloudAccountID
	object.cloudProviderID = b.cloudProviderID
	object.clusterID = b.clusterID
	object.disconnected = b.disconnected
	object.displayName = b.displayName
	object.externalClusterID = b.externalClusterID
	object.managed = b.managed
	object.productID = b.productID
	object.productCategory = b.productCategory
	object.quotaVersion = b.quotaVersion
	object.reserve = b.reserve
	if b.resources != nil {
		object.resources = make([]*ReservedResource, len(b.resources))
		for i, v := range b.resources {
			object.resources[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
