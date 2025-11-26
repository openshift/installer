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

// GCP Network configuration of a cluster.
type GCPNetworkBuilder struct {
	fieldSet_          []bool
	vpcName            string
	vpcProjectID       string
	computeSubnet      string
	controlPlaneSubnet string
}

// NewGCPNetwork creates a new builder of 'GCP_network' objects.
func NewGCPNetwork() *GCPNetworkBuilder {
	return &GCPNetworkBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPNetworkBuilder) Empty() bool {
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

// VPCName sets the value of the 'VPC_name' attribute to the given value.
func (b *GCPNetworkBuilder) VPCName(value string) *GCPNetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.vpcName = value
	b.fieldSet_[0] = true
	return b
}

// VPCProjectID sets the value of the 'VPC_project_ID' attribute to the given value.
func (b *GCPNetworkBuilder) VPCProjectID(value string) *GCPNetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.vpcProjectID = value
	b.fieldSet_[1] = true
	return b
}

// ComputeSubnet sets the value of the 'compute_subnet' attribute to the given value.
func (b *GCPNetworkBuilder) ComputeSubnet(value string) *GCPNetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.computeSubnet = value
	b.fieldSet_[2] = true
	return b
}

// ControlPlaneSubnet sets the value of the 'control_plane_subnet' attribute to the given value.
func (b *GCPNetworkBuilder) ControlPlaneSubnet(value string) *GCPNetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.controlPlaneSubnet = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPNetworkBuilder) Copy(object *GCPNetwork) *GCPNetworkBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.vpcName = object.vpcName
	b.vpcProjectID = object.vpcProjectID
	b.computeSubnet = object.computeSubnet
	b.controlPlaneSubnet = object.controlPlaneSubnet
	return b
}

// Build creates a 'GCP_network' object using the configuration stored in the builder.
func (b *GCPNetworkBuilder) Build() (object *GCPNetwork, err error) {
	object = new(GCPNetwork)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.vpcName = b.vpcName
	object.vpcProjectID = b.vpcProjectID
	object.computeSubnet = b.computeSubnet
	object.controlPlaneSubnet = b.controlPlaneSubnet
	return
}
