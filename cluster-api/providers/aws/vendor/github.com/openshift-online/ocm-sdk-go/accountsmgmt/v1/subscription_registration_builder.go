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

// SubscriptionRegistrationBuilder contains the data and logic needed to build 'subscription_registration' objects.
//
// Registration of a new subscription.
type SubscriptionRegistrationBuilder struct {
	bitmap_     uint32
	clusterUUID string
	consoleURL  string
	displayName string
	planID      PlanID
	status      string
}

// NewSubscriptionRegistration creates a new builder of 'subscription_registration' objects.
func NewSubscriptionRegistration() *SubscriptionRegistrationBuilder {
	return &SubscriptionRegistrationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubscriptionRegistrationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) ClusterUUID(value string) *SubscriptionRegistrationBuilder {
	b.clusterUUID = value
	b.bitmap_ |= 1
	return b
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) ConsoleURL(value string) *SubscriptionRegistrationBuilder {
	b.consoleURL = value
	b.bitmap_ |= 2
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) DisplayName(value string) *SubscriptionRegistrationBuilder {
	b.displayName = value
	b.bitmap_ |= 4
	return b
}

// PlanID sets the value of the 'plan_ID' attribute to the given value.
//
// Plan ID of subscription.
func (b *SubscriptionRegistrationBuilder) PlanID(value PlanID) *SubscriptionRegistrationBuilder {
	b.planID = value
	b.bitmap_ |= 8
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) Status(value string) *SubscriptionRegistrationBuilder {
	b.status = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubscriptionRegistrationBuilder) Copy(object *SubscriptionRegistration) *SubscriptionRegistrationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clusterUUID = object.clusterUUID
	b.consoleURL = object.consoleURL
	b.displayName = object.displayName
	b.planID = object.planID
	b.status = object.status
	return b
}

// Build creates a 'subscription_registration' object using the configuration stored in the builder.
func (b *SubscriptionRegistrationBuilder) Build() (object *SubscriptionRegistration, err error) {
	object = new(SubscriptionRegistration)
	object.bitmap_ = b.bitmap_
	object.clusterUUID = b.clusterUUID
	object.consoleURL = b.consoleURL
	object.displayName = b.displayName
	object.planID = b.planID
	object.status = b.status
	return
}
