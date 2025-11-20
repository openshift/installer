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

type ClusterAuthorizationRequestBuilder struct {
	fieldSet_         []bool
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
	rhRegionID        string
	scope             string
	byoc              bool
	disconnected      bool
	managed           bool
	reserve           bool
}

// NewClusterAuthorizationRequest creates a new builder of 'cluster_authorization_request' objects.
func NewClusterAuthorizationRequest() *ClusterAuthorizationRequestBuilder {
	return &ClusterAuthorizationRequestBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterAuthorizationRequestBuilder) Empty() bool {
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

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) BYOC(value bool) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.byoc = value
	b.fieldSet_[0] = true
	return b
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) AccountUsername(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.accountUsername = value
	b.fieldSet_[1] = true
	return b
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) AvailabilityZone(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.availabilityZone = value
	b.fieldSet_[2] = true
	return b
}

// CloudAccountID sets the value of the 'cloud_account_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) CloudAccountID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.cloudAccountID = value
	b.fieldSet_[3] = true
	return b
}

// CloudProviderID sets the value of the 'cloud_provider_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) CloudProviderID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.cloudProviderID = value
	b.fieldSet_[4] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ClusterID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.clusterID = value
	b.fieldSet_[5] = true
	return b
}

// Disconnected sets the value of the 'disconnected' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Disconnected(value bool) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.disconnected = value
	b.fieldSet_[6] = true
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) DisplayName(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.displayName = value
	b.fieldSet_[7] = true
	return b
}

// ExternalClusterID sets the value of the 'external_cluster_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ExternalClusterID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.externalClusterID = value
	b.fieldSet_[8] = true
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Managed(value bool) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.managed = value
	b.fieldSet_[9] = true
	return b
}

// ProductID sets the value of the 'product_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ProductID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.productID = value
	b.fieldSet_[10] = true
	return b
}

// ProductCategory sets the value of the 'product_category' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) ProductCategory(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.productCategory = value
	b.fieldSet_[11] = true
	return b
}

// QuotaVersion sets the value of the 'quota_version' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) QuotaVersion(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.quotaVersion = value
	b.fieldSet_[12] = true
	return b
}

// Reserve sets the value of the 'reserve' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Reserve(value bool) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.reserve = value
	b.fieldSet_[13] = true
	return b
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *ClusterAuthorizationRequestBuilder) Resources(values ...*ReservedResourceBuilder) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.resources = make([]*ReservedResourceBuilder, len(values))
	copy(b.resources, values)
	b.fieldSet_[14] = true
	return b
}

// RhRegionID sets the value of the 'rh_region_ID' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) RhRegionID(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.rhRegionID = value
	b.fieldSet_[15] = true
	return b
}

// Scope sets the value of the 'scope' attribute to the given value.
func (b *ClusterAuthorizationRequestBuilder) Scope(value string) *ClusterAuthorizationRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.scope = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterAuthorizationRequestBuilder) Copy(object *ClusterAuthorizationRequest) *ClusterAuthorizationRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	b.rhRegionID = object.rhRegionID
	b.scope = object.scope
	return b
}

// Build creates a 'cluster_authorization_request' object using the configuration stored in the builder.
func (b *ClusterAuthorizationRequestBuilder) Build() (object *ClusterAuthorizationRequest, err error) {
	object = new(ClusterAuthorizationRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	object.rhRegionID = b.rhRegionID
	object.scope = b.scope
	return
}
