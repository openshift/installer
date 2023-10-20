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

// ClusterRegistrationRequestBuilder contains the data and logic needed to build 'cluster_registration_request' objects.
type ClusterRegistrationRequestBuilder struct {
	bitmap_            uint32
	authorizationToken string
	clusterID          string
}

// NewClusterRegistrationRequest creates a new builder of 'cluster_registration_request' objects.
func NewClusterRegistrationRequest() *ClusterRegistrationRequestBuilder {
	return &ClusterRegistrationRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterRegistrationRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AuthorizationToken sets the value of the 'authorization_token' attribute to the given value.
func (b *ClusterRegistrationRequestBuilder) AuthorizationToken(value string) *ClusterRegistrationRequestBuilder {
	b.authorizationToken = value
	b.bitmap_ |= 1
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterRegistrationRequestBuilder) ClusterID(value string) *ClusterRegistrationRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterRegistrationRequestBuilder) Copy(object *ClusterRegistrationRequest) *ClusterRegistrationRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.authorizationToken = object.authorizationToken
	b.clusterID = object.clusterID
	return b
}

// Build creates a 'cluster_registration_request' object using the configuration stored in the builder.
func (b *ClusterRegistrationRequestBuilder) Build() (object *ClusterRegistrationRequest, err error) {
	object = new(ClusterRegistrationRequest)
	object.bitmap_ = b.bitmap_
	object.authorizationToken = b.authorizationToken
	object.clusterID = b.clusterID
	return
}
