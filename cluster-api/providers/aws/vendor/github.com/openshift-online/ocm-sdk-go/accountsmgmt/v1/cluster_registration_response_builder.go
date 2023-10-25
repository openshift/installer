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

// ClusterRegistrationResponseBuilder contains the data and logic needed to build 'cluster_registration_response' objects.
type ClusterRegistrationResponseBuilder struct {
	bitmap_            uint32
	accountID          string
	authorizationToken string
	clusterID          string
	expiresAt          string
}

// NewClusterRegistrationResponse creates a new builder of 'cluster_registration_response' objects.
func NewClusterRegistrationResponse() *ClusterRegistrationResponseBuilder {
	return &ClusterRegistrationResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterRegistrationResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) AccountID(value string) *ClusterRegistrationResponseBuilder {
	b.accountID = value
	b.bitmap_ |= 1
	return b
}

// AuthorizationToken sets the value of the 'authorization_token' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) AuthorizationToken(value string) *ClusterRegistrationResponseBuilder {
	b.authorizationToken = value
	b.bitmap_ |= 2
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) ClusterID(value string) *ClusterRegistrationResponseBuilder {
	b.clusterID = value
	b.bitmap_ |= 4
	return b
}

// ExpiresAt sets the value of the 'expires_at' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) ExpiresAt(value string) *ClusterRegistrationResponseBuilder {
	b.expiresAt = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterRegistrationResponseBuilder) Copy(object *ClusterRegistrationResponse) *ClusterRegistrationResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountID = object.accountID
	b.authorizationToken = object.authorizationToken
	b.clusterID = object.clusterID
	b.expiresAt = object.expiresAt
	return b
}

// Build creates a 'cluster_registration_response' object using the configuration stored in the builder.
func (b *ClusterRegistrationResponseBuilder) Build() (object *ClusterRegistrationResponse, err error) {
	object = new(ClusterRegistrationResponse)
	object.bitmap_ = b.bitmap_
	object.accountID = b.accountID
	object.authorizationToken = b.authorizationToken
	object.clusterID = b.clusterID
	object.expiresAt = b.expiresAt
	return
}
