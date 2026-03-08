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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

type SubnetNetworkVerificationBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	details   []string
	platform  Platform
	state     string
	tags      map[string]string
}

// NewSubnetNetworkVerification creates a new builder of 'subnet_network_verification' objects.
func NewSubnetNetworkVerification() *SubnetNetworkVerificationBuilder {
	return &SubnetNetworkVerificationBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *SubnetNetworkVerificationBuilder) Link(value bool) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *SubnetNetworkVerificationBuilder) ID(value string) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *SubnetNetworkVerificationBuilder) HREF(value string) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubnetNetworkVerificationBuilder) Empty() bool {
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

// Details sets the value of the 'details' attribute to the given values.
func (b *SubnetNetworkVerificationBuilder) Details(values ...string) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.details = make([]string, len(values))
	copy(b.details, values)
	b.fieldSet_[3] = true
	return b
}

// Platform sets the value of the 'platform' attribute to the given value.
//
// Representation of an platform type field.
func (b *SubnetNetworkVerificationBuilder) Platform(value Platform) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.platform = value
	b.fieldSet_[4] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *SubnetNetworkVerificationBuilder) State(value string) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.state = value
	b.fieldSet_[5] = true
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *SubnetNetworkVerificationBuilder) Tags(value map[string]string) *SubnetNetworkVerificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.tags = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubnetNetworkVerificationBuilder) Copy(object *SubnetNetworkVerification) *SubnetNetworkVerificationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.details != nil {
		b.details = make([]string, len(object.details))
		copy(b.details, object.details)
	} else {
		b.details = nil
	}
	b.platform = object.platform
	b.state = object.state
	if len(object.tags) > 0 {
		b.tags = map[string]string{}
		for k, v := range object.tags {
			b.tags[k] = v
		}
	} else {
		b.tags = nil
	}
	return b
}

// Build creates a 'subnet_network_verification' object using the configuration stored in the builder.
func (b *SubnetNetworkVerificationBuilder) Build() (object *SubnetNetworkVerification, err error) {
	object = new(SubnetNetworkVerification)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.details != nil {
		object.details = make([]string, len(b.details))
		copy(object.details, b.details)
	}
	object.platform = b.platform
	object.state = b.state
	if b.tags != nil {
		object.tags = make(map[string]string)
		for k, v := range b.tags {
			object.tags[k] = v
		}
	}
	return
}
