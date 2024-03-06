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

// NotificationDetailsRequestBuilder contains the data and logic needed to build 'notification_details_request' objects.
//
// This struct is a request to get a templated email to a user related to this.
// subscription/cluster.
type NotificationDetailsRequestBuilder struct {
	bitmap_                 uint32
	bccAddress              string
	clusterID               string
	clusterUUID             string
	subject                 string
	subscriptionID          string
	includeRedHatAssociates bool
	internalOnly            bool
}

// NewNotificationDetailsRequest creates a new builder of 'notification_details_request' objects.
func NewNotificationDetailsRequest() *NotificationDetailsRequestBuilder {
	return &NotificationDetailsRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NotificationDetailsRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// BccAddress sets the value of the 'bcc_address' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) BccAddress(value string) *NotificationDetailsRequestBuilder {
	b.bccAddress = value
	b.bitmap_ |= 1
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) ClusterID(value string) *NotificationDetailsRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 2
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) ClusterUUID(value string) *NotificationDetailsRequestBuilder {
	b.clusterUUID = value
	b.bitmap_ |= 4
	return b
}

// IncludeRedHatAssociates sets the value of the 'include_red_hat_associates' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) IncludeRedHatAssociates(value bool) *NotificationDetailsRequestBuilder {
	b.includeRedHatAssociates = value
	b.bitmap_ |= 8
	return b
}

// InternalOnly sets the value of the 'internal_only' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) InternalOnly(value bool) *NotificationDetailsRequestBuilder {
	b.internalOnly = value
	b.bitmap_ |= 16
	return b
}

// Subject sets the value of the 'subject' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) Subject(value string) *NotificationDetailsRequestBuilder {
	b.subject = value
	b.bitmap_ |= 32
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) SubscriptionID(value string) *NotificationDetailsRequestBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NotificationDetailsRequestBuilder) Copy(object *NotificationDetailsRequest) *NotificationDetailsRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.bccAddress = object.bccAddress
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.includeRedHatAssociates = object.includeRedHatAssociates
	b.internalOnly = object.internalOnly
	b.subject = object.subject
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'notification_details_request' object using the configuration stored in the builder.
func (b *NotificationDetailsRequestBuilder) Build() (object *NotificationDetailsRequest, err error) {
	object = new(NotificationDetailsRequest)
	object.bitmap_ = b.bitmap_
	object.bccAddress = b.bccAddress
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.includeRedHatAssociates = b.includeRedHatAssociates
	object.internalOnly = b.internalOnly
	object.subject = b.subject
	object.subscriptionID = b.subscriptionID
	return
}
