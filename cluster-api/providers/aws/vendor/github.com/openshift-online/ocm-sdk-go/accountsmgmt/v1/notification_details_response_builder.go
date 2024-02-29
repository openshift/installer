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

// NotificationDetailsResponseBuilder contains the data and logic needed to build 'notification_details_response' objects.
//
// This class is a single response item for the notify details list.
type NotificationDetailsResponseBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	key     string
	value   string
}

// NewNotificationDetailsResponse creates a new builder of 'notification_details_response' objects.
func NewNotificationDetailsResponse() *NotificationDetailsResponseBuilder {
	return &NotificationDetailsResponseBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *NotificationDetailsResponseBuilder) Link(value bool) *NotificationDetailsResponseBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *NotificationDetailsResponseBuilder) ID(value string) *NotificationDetailsResponseBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *NotificationDetailsResponseBuilder) HREF(value string) *NotificationDetailsResponseBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NotificationDetailsResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Key sets the value of the 'key' attribute to the given value.
func (b *NotificationDetailsResponseBuilder) Key(value string) *NotificationDetailsResponseBuilder {
	b.key = value
	b.bitmap_ |= 8
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *NotificationDetailsResponseBuilder) Value(value string) *NotificationDetailsResponseBuilder {
	b.value = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NotificationDetailsResponseBuilder) Copy(object *NotificationDetailsResponse) *NotificationDetailsResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.key = object.key
	b.value = object.value
	return b
}

// Build creates a 'notification_details_response' object using the configuration stored in the builder.
func (b *NotificationDetailsResponseBuilder) Build() (object *NotificationDetailsResponse, err error) {
	object = new(NotificationDetailsResponse)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.key = b.key
	object.value = b.value
	return
}
