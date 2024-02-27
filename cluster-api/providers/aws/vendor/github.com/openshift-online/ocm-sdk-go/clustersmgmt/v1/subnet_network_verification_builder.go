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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// SubnetNetworkVerificationBuilder contains the data and logic needed to build 'subnet_network_verification' objects.
type SubnetNetworkVerificationBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	details []string
	state   string
	tags    map[string]string
}

// NewSubnetNetworkVerification creates a new builder of 'subnet_network_verification' objects.
func NewSubnetNetworkVerification() *SubnetNetworkVerificationBuilder {
	return &SubnetNetworkVerificationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SubnetNetworkVerificationBuilder) Link(value bool) *SubnetNetworkVerificationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SubnetNetworkVerificationBuilder) ID(value string) *SubnetNetworkVerificationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SubnetNetworkVerificationBuilder) HREF(value string) *SubnetNetworkVerificationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubnetNetworkVerificationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Details sets the value of the 'details' attribute to the given values.
func (b *SubnetNetworkVerificationBuilder) Details(values ...string) *SubnetNetworkVerificationBuilder {
	b.details = make([]string, len(values))
	copy(b.details, values)
	b.bitmap_ |= 8
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *SubnetNetworkVerificationBuilder) State(value string) *SubnetNetworkVerificationBuilder {
	b.state = value
	b.bitmap_ |= 16
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *SubnetNetworkVerificationBuilder) Tags(value map[string]string) *SubnetNetworkVerificationBuilder {
	b.tags = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubnetNetworkVerificationBuilder) Copy(object *SubnetNetworkVerification) *SubnetNetworkVerificationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.details != nil {
		b.details = make([]string, len(object.details))
		copy(b.details, object.details)
	} else {
		b.details = nil
	}
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
	object.bitmap_ = b.bitmap_
	if b.details != nil {
		object.details = make([]string, len(b.details))
		copy(object.details, b.details)
	}
	object.state = b.state
	if b.tags != nil {
		object.tags = make(map[string]string)
		for k, v := range b.tags {
			object.tags[k] = v
		}
	}
	return
}
