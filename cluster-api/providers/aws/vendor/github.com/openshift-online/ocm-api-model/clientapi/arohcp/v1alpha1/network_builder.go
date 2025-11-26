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

// Network configuration of a cluster.
type NetworkBuilder struct {
	fieldSet_   []bool
	hostPrefix  int
	machineCIDR string
	podCIDR     string
	serviceCIDR string
	type_       string
}

// NewNetwork creates a new builder of 'network' objects.
func NewNetwork() *NetworkBuilder {
	return &NetworkBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NetworkBuilder) Empty() bool {
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

// HostPrefix sets the value of the 'host_prefix' attribute to the given value.
func (b *NetworkBuilder) HostPrefix(value int) *NetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.hostPrefix = value
	b.fieldSet_[0] = true
	return b
}

// MachineCIDR sets the value of the 'machine_CIDR' attribute to the given value.
func (b *NetworkBuilder) MachineCIDR(value string) *NetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.machineCIDR = value
	b.fieldSet_[1] = true
	return b
}

// PodCIDR sets the value of the 'pod_CIDR' attribute to the given value.
func (b *NetworkBuilder) PodCIDR(value string) *NetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.podCIDR = value
	b.fieldSet_[2] = true
	return b
}

// ServiceCIDR sets the value of the 'service_CIDR' attribute to the given value.
func (b *NetworkBuilder) ServiceCIDR(value string) *NetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.serviceCIDR = value
	b.fieldSet_[3] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *NetworkBuilder) Type(value string) *NetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.type_ = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NetworkBuilder) Copy(object *Network) *NetworkBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.hostPrefix = object.hostPrefix
	b.machineCIDR = object.machineCIDR
	b.podCIDR = object.podCIDR
	b.serviceCIDR = object.serviceCIDR
	b.type_ = object.type_
	return b
}

// Build creates a 'network' object using the configuration stored in the builder.
func (b *NetworkBuilder) Build() (object *Network, err error) {
	object = new(Network)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.hostPrefix = b.hostPrefix
	object.machineCIDR = b.machineCIDR
	object.podCIDR = b.podCIDR
	object.serviceCIDR = b.serviceCIDR
	object.type_ = b.type_
	return
}
