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

// NetworkVerificationBuilder contains the data and logic needed to build 'network_verification' objects.
type NetworkVerificationBuilder struct {
	bitmap_           uint32
	cloudProviderData *CloudProviderDataBuilder
	clusterId         string
	items             []*SubnetNetworkVerificationBuilder
	platform          Platform
	total             int
}

// NewNetworkVerification creates a new builder of 'network_verification' objects.
func NewNetworkVerification() *NetworkVerificationBuilder {
	return &NetworkVerificationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NetworkVerificationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// CloudProviderData sets the value of the 'cloud_provider_data' attribute to the given value.
//
// Description of a cloud provider data used for cloud provider inquiries.
func (b *NetworkVerificationBuilder) CloudProviderData(value *CloudProviderDataBuilder) *NetworkVerificationBuilder {
	b.cloudProviderData = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *NetworkVerificationBuilder) ClusterId(value string) *NetworkVerificationBuilder {
	b.clusterId = value
	b.bitmap_ |= 2
	return b
}

// Items sets the value of the 'items' attribute to the given values.
func (b *NetworkVerificationBuilder) Items(values ...*SubnetNetworkVerificationBuilder) *NetworkVerificationBuilder {
	b.items = make([]*SubnetNetworkVerificationBuilder, len(values))
	copy(b.items, values)
	b.bitmap_ |= 4
	return b
}

// Platform sets the value of the 'platform' attribute to the given value.
//
// Representation of an platform type field.
func (b *NetworkVerificationBuilder) Platform(value Platform) *NetworkVerificationBuilder {
	b.platform = value
	b.bitmap_ |= 8
	return b
}

// Total sets the value of the 'total' attribute to the given value.
func (b *NetworkVerificationBuilder) Total(value int) *NetworkVerificationBuilder {
	b.total = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NetworkVerificationBuilder) Copy(object *NetworkVerification) *NetworkVerificationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.cloudProviderData != nil {
		b.cloudProviderData = NewCloudProviderData().Copy(object.cloudProviderData)
	} else {
		b.cloudProviderData = nil
	}
	b.clusterId = object.clusterId
	if object.items != nil {
		b.items = make([]*SubnetNetworkVerificationBuilder, len(object.items))
		for i, v := range object.items {
			b.items[i] = NewSubnetNetworkVerification().Copy(v)
		}
	} else {
		b.items = nil
	}
	b.platform = object.platform
	b.total = object.total
	return b
}

// Build creates a 'network_verification' object using the configuration stored in the builder.
func (b *NetworkVerificationBuilder) Build() (object *NetworkVerification, err error) {
	object = new(NetworkVerification)
	object.bitmap_ = b.bitmap_
	if b.cloudProviderData != nil {
		object.cloudProviderData, err = b.cloudProviderData.Build()
		if err != nil {
			return
		}
	}
	object.clusterId = b.clusterId
	if b.items != nil {
		object.items = make([]*SubnetNetworkVerification, len(b.items))
		for i, v := range b.items {
			object.items[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.platform = b.platform
	object.total = b.total
	return
}
