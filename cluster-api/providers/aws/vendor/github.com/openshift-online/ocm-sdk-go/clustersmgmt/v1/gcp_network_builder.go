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

// GCPNetworkBuilder contains the data and logic needed to build 'GCP_network' objects.
//
// GCP Network configuration of a cluster.
type GCPNetworkBuilder struct {
	bitmap_            uint32
	vpcName            string
	vpcProjectID       string
	computeSubnet      string
	controlPlaneSubnet string
}

// NewGCPNetwork creates a new builder of 'GCP_network' objects.
func NewGCPNetwork() *GCPNetworkBuilder {
	return &GCPNetworkBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPNetworkBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// VPCName sets the value of the 'VPC_name' attribute to the given value.
func (b *GCPNetworkBuilder) VPCName(value string) *GCPNetworkBuilder {
	b.vpcName = value
	b.bitmap_ |= 1
	return b
}

// VPCProjectID sets the value of the 'VPC_project_ID' attribute to the given value.
func (b *GCPNetworkBuilder) VPCProjectID(value string) *GCPNetworkBuilder {
	b.vpcProjectID = value
	b.bitmap_ |= 2
	return b
}

// ComputeSubnet sets the value of the 'compute_subnet' attribute to the given value.
func (b *GCPNetworkBuilder) ComputeSubnet(value string) *GCPNetworkBuilder {
	b.computeSubnet = value
	b.bitmap_ |= 4
	return b
}

// ControlPlaneSubnet sets the value of the 'control_plane_subnet' attribute to the given value.
func (b *GCPNetworkBuilder) ControlPlaneSubnet(value string) *GCPNetworkBuilder {
	b.controlPlaneSubnet = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPNetworkBuilder) Copy(object *GCPNetwork) *GCPNetworkBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.vpcName = object.vpcName
	b.vpcProjectID = object.vpcProjectID
	b.computeSubnet = object.computeSubnet
	b.controlPlaneSubnet = object.controlPlaneSubnet
	return b
}

// Build creates a 'GCP_network' object using the configuration stored in the builder.
func (b *GCPNetworkBuilder) Build() (object *GCPNetwork, err error) {
	object = new(GCPNetwork)
	object.bitmap_ = b.bitmap_
	object.vpcName = b.vpcName
	object.vpcProjectID = b.vpcProjectID
	object.computeSubnet = b.computeSubnet
	object.controlPlaneSubnet = b.controlPlaneSubnet
	return
}
