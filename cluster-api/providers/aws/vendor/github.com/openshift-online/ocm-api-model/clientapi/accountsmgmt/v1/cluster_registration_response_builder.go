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

type ClusterRegistrationResponseBuilder struct {
	fieldSet_          []bool
	accountID          string
	authorizationToken string
	clusterID          string
	expiresAt          string
}

// NewClusterRegistrationResponse creates a new builder of 'cluster_registration_response' objects.
func NewClusterRegistrationResponse() *ClusterRegistrationResponseBuilder {
	return &ClusterRegistrationResponseBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterRegistrationResponseBuilder) Empty() bool {
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

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) AccountID(value string) *ClusterRegistrationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.accountID = value
	b.fieldSet_[0] = true
	return b
}

// AuthorizationToken sets the value of the 'authorization_token' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) AuthorizationToken(value string) *ClusterRegistrationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.authorizationToken = value
	b.fieldSet_[1] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) ClusterID(value string) *ClusterRegistrationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.clusterID = value
	b.fieldSet_[2] = true
	return b
}

// ExpiresAt sets the value of the 'expires_at' attribute to the given value.
func (b *ClusterRegistrationResponseBuilder) ExpiresAt(value string) *ClusterRegistrationResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.expiresAt = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterRegistrationResponseBuilder) Copy(object *ClusterRegistrationResponse) *ClusterRegistrationResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.accountID = object.accountID
	b.authorizationToken = object.authorizationToken
	b.clusterID = object.clusterID
	b.expiresAt = object.expiresAt
	return b
}

// Build creates a 'cluster_registration_response' object using the configuration stored in the builder.
func (b *ClusterRegistrationResponseBuilder) Build() (object *ClusterRegistrationResponse, err error) {
	object = new(ClusterRegistrationResponse)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountID = b.accountID
	object.authorizationToken = b.authorizationToken
	object.clusterID = b.clusterID
	object.expiresAt = b.expiresAt
	return
}
