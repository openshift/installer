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

// GenericNotifyDetailsResponseBuilder contains the data and logic needed to build 'generic_notify_details_response' objects.
//
// class that defines notify details response in general.
type GenericNotifyDetailsResponseBuilder struct {
	bitmap_    uint32
	id         string
	href       string
	associates []string
	items      []*NotificationDetailsResponseBuilder
	recipients []string
}

// NewGenericNotifyDetailsResponse creates a new builder of 'generic_notify_details_response' objects.
func NewGenericNotifyDetailsResponse() *GenericNotifyDetailsResponseBuilder {
	return &GenericNotifyDetailsResponseBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *GenericNotifyDetailsResponseBuilder) Link(value bool) *GenericNotifyDetailsResponseBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *GenericNotifyDetailsResponseBuilder) ID(value string) *GenericNotifyDetailsResponseBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *GenericNotifyDetailsResponseBuilder) HREF(value string) *GenericNotifyDetailsResponseBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GenericNotifyDetailsResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Associates sets the value of the 'associates' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Associates(values ...string) *GenericNotifyDetailsResponseBuilder {
	b.associates = make([]string, len(values))
	copy(b.associates, values)
	b.bitmap_ |= 8
	return b
}

// Items sets the value of the 'items' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Items(values ...*NotificationDetailsResponseBuilder) *GenericNotifyDetailsResponseBuilder {
	b.items = make([]*NotificationDetailsResponseBuilder, len(values))
	copy(b.items, values)
	b.bitmap_ |= 16
	return b
}

// Recipients sets the value of the 'recipients' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Recipients(values ...string) *GenericNotifyDetailsResponseBuilder {
	b.recipients = make([]string, len(values))
	copy(b.recipients, values)
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GenericNotifyDetailsResponseBuilder) Copy(object *GenericNotifyDetailsResponse) *GenericNotifyDetailsResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.associates != nil {
		b.associates = make([]string, len(object.associates))
		copy(b.associates, object.associates)
	} else {
		b.associates = nil
	}
	if object.items != nil {
		b.items = make([]*NotificationDetailsResponseBuilder, len(object.items))
		for i, v := range object.items {
			b.items[i] = NewNotificationDetailsResponse().Copy(v)
		}
	} else {
		b.items = nil
	}
	if object.recipients != nil {
		b.recipients = make([]string, len(object.recipients))
		copy(b.recipients, object.recipients)
	} else {
		b.recipients = nil
	}
	return b
}

// Build creates a 'generic_notify_details_response' object using the configuration stored in the builder.
func (b *GenericNotifyDetailsResponseBuilder) Build() (object *GenericNotifyDetailsResponse, err error) {
	object = new(GenericNotifyDetailsResponse)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.associates != nil {
		object.associates = make([]string, len(b.associates))
		copy(object.associates, b.associates)
	}
	if b.items != nil {
		object.items = make([]*NotificationDetailsResponse, len(b.items))
		for i, v := range b.items {
			object.items[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.recipients != nil {
		object.recipients = make([]string, len(b.recipients))
		copy(object.recipients, b.recipients)
	}
	return
}
