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

// This struct is a request to get a templated email to a user related to this.
// subscription/cluster.
type NotificationDetailsRequestBuilder struct {
	fieldSet_               []bool
	bccAddress              string
	clusterID               string
	clusterUUID             string
	logType                 string
	subject                 string
	subscriptionID          string
	includeRedHatAssociates bool
	internalOnly            bool
}

// NewNotificationDetailsRequest creates a new builder of 'notification_details_request' objects.
func NewNotificationDetailsRequest() *NotificationDetailsRequestBuilder {
	return &NotificationDetailsRequestBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NotificationDetailsRequestBuilder) Empty() bool {
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

// BccAddress sets the value of the 'bcc_address' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) BccAddress(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.bccAddress = value
	b.fieldSet_[0] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) ClusterID(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.clusterID = value
	b.fieldSet_[1] = true
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) ClusterUUID(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.clusterUUID = value
	b.fieldSet_[2] = true
	return b
}

// IncludeRedHatAssociates sets the value of the 'include_red_hat_associates' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) IncludeRedHatAssociates(value bool) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.includeRedHatAssociates = value
	b.fieldSet_[3] = true
	return b
}

// InternalOnly sets the value of the 'internal_only' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) InternalOnly(value bool) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.internalOnly = value
	b.fieldSet_[4] = true
	return b
}

// LogType sets the value of the 'log_type' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) LogType(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.logType = value
	b.fieldSet_[5] = true
	return b
}

// Subject sets the value of the 'subject' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) Subject(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.subject = value
	b.fieldSet_[6] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *NotificationDetailsRequestBuilder) SubscriptionID(value string) *NotificationDetailsRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.subscriptionID = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NotificationDetailsRequestBuilder) Copy(object *NotificationDetailsRequest) *NotificationDetailsRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.bccAddress = object.bccAddress
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.includeRedHatAssociates = object.includeRedHatAssociates
	b.internalOnly = object.internalOnly
	b.logType = object.logType
	b.subject = object.subject
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'notification_details_request' object using the configuration stored in the builder.
func (b *NotificationDetailsRequestBuilder) Build() (object *NotificationDetailsRequest, err error) {
	object = new(NotificationDetailsRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.bccAddress = b.bccAddress
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.includeRedHatAssociates = b.includeRedHatAssociates
	object.internalOnly = b.internalOnly
	object.logType = b.logType
	object.subject = b.subject
	object.subscriptionID = b.subscriptionID
	return
}
