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

// Registration of a new cluster to the service.
//
// For example, to register a cluster that has been provisioned outside
// of this service, send a a request like this:
//
// ```http
// POST /api/clusters_mgmt/v1/register_cluster HTTP/1.1
// ```
//
// With a request body like this:
//
// ```json
//
//	{
//	  "external_id": "d656aecf-11a6-4782-ad86-8f72638449ba",
//	  "subscription_id": "...",
//	  "organization_id": "..."
//	}
//
// ```
type ClusterRegistrationBuilder struct {
	fieldSet_      []bool
	consoleUrl     string
	externalID     string
	organizationID string
	subscriptionID string
}

// NewClusterRegistration creates a new builder of 'cluster_registration' objects.
func NewClusterRegistration() *ClusterRegistrationBuilder {
	return &ClusterRegistrationBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterRegistrationBuilder) Empty() bool {
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

// ConsoleUrl sets the value of the 'console_url' attribute to the given value.
func (b *ClusterRegistrationBuilder) ConsoleUrl(value string) *ClusterRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.consoleUrl = value
	b.fieldSet_[0] = true
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *ClusterRegistrationBuilder) ExternalID(value string) *ClusterRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.externalID = value
	b.fieldSet_[1] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *ClusterRegistrationBuilder) OrganizationID(value string) *ClusterRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.organizationID = value
	b.fieldSet_[2] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *ClusterRegistrationBuilder) SubscriptionID(value string) *ClusterRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.subscriptionID = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterRegistrationBuilder) Copy(object *ClusterRegistration) *ClusterRegistrationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.consoleUrl = object.consoleUrl
	b.externalID = object.externalID
	b.organizationID = object.organizationID
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'cluster_registration' object using the configuration stored in the builder.
func (b *ClusterRegistrationBuilder) Build() (object *ClusterRegistration, err error) {
	object = new(ClusterRegistration)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.consoleUrl = b.consoleUrl
	object.externalID = b.externalID
	object.organizationID = b.organizationID
	object.subscriptionID = b.subscriptionID
	return
}
