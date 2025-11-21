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

// Registration of a new subscription.
type SubscriptionRegistrationBuilder struct {
	fieldSet_   []bool
	clusterUUID string
	consoleURL  string
	displayName string
	planID      PlanID
	status      string
}

// NewSubscriptionRegistration creates a new builder of 'subscription_registration' objects.
func NewSubscriptionRegistration() *SubscriptionRegistrationBuilder {
	return &SubscriptionRegistrationBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubscriptionRegistrationBuilder) Empty() bool {
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

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) ClusterUUID(value string) *SubscriptionRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.clusterUUID = value
	b.fieldSet_[0] = true
	return b
}

// ConsoleURL sets the value of the 'console_URL' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) ConsoleURL(value string) *SubscriptionRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.consoleURL = value
	b.fieldSet_[1] = true
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) DisplayName(value string) *SubscriptionRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.displayName = value
	b.fieldSet_[2] = true
	return b
}

// PlanID sets the value of the 'plan_ID' attribute to the given value.
//
// Plan ID of subscription.
func (b *SubscriptionRegistrationBuilder) PlanID(value PlanID) *SubscriptionRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.planID = value
	b.fieldSet_[3] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *SubscriptionRegistrationBuilder) Status(value string) *SubscriptionRegistrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.status = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubscriptionRegistrationBuilder) Copy(object *SubscriptionRegistration) *SubscriptionRegistrationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterUUID = b.clusterUUID
	object.consoleURL = b.consoleURL
	object.displayName = b.displayName
	object.planID = b.planID
	object.status = b.status
	return
}
