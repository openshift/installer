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

// class that defines notify details response in general.
type GenericNotifyDetailsResponseBuilder struct {
	fieldSet_  []bool
	id         string
	href       string
	associates []string
	items      []*NotificationDetailsResponseBuilder
	recipients []string
}

// NewGenericNotifyDetailsResponse creates a new builder of 'generic_notify_details_response' objects.
func NewGenericNotifyDetailsResponse() *GenericNotifyDetailsResponseBuilder {
	return &GenericNotifyDetailsResponseBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *GenericNotifyDetailsResponseBuilder) Link(value bool) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *GenericNotifyDetailsResponseBuilder) ID(value string) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *GenericNotifyDetailsResponseBuilder) HREF(value string) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GenericNotifyDetailsResponseBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// Associates sets the value of the 'associates' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Associates(values ...string) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.associates = make([]string, len(values))
	copy(b.associates, values)
	b.fieldSet_[3] = true
	return b
}

// Items sets the value of the 'items' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Items(values ...*NotificationDetailsResponseBuilder) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.items = make([]*NotificationDetailsResponseBuilder, len(values))
	copy(b.items, values)
	b.fieldSet_[4] = true
	return b
}

// Recipients sets the value of the 'recipients' attribute to the given values.
func (b *GenericNotifyDetailsResponseBuilder) Recipients(values ...string) *GenericNotifyDetailsResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.recipients = make([]string, len(values))
	copy(b.recipients, values)
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GenericNotifyDetailsResponseBuilder) Copy(object *GenericNotifyDetailsResponse) *GenericNotifyDetailsResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
