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

// GCPMachinePoolBuilder contains the data and logic needed to build 'GCP_machine_pool' objects.
//
// Representation of gcp machine pool specific parameters.
type GCPMachinePoolBuilder struct {
	bitmap_    uint32
	secureBoot bool
}

// NewGCPMachinePool creates a new builder of 'GCP_machine_pool' objects.
func NewGCPMachinePool() *GCPMachinePoolBuilder {
	return &GCPMachinePoolBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPMachinePoolBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// SecureBoot sets the value of the 'secure_boot' attribute to the given value.
func (b *GCPMachinePoolBuilder) SecureBoot(value bool) *GCPMachinePoolBuilder {
	b.secureBoot = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPMachinePoolBuilder) Copy(object *GCPMachinePool) *GCPMachinePoolBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.secureBoot = object.secureBoot
	return b
}

// Build creates a 'GCP_machine_pool' object using the configuration stored in the builder.
func (b *GCPMachinePoolBuilder) Build() (object *GCPMachinePool, err error) {
	object = new(GCPMachinePool)
	object.bitmap_ = b.bitmap_
	object.secureBoot = b.secureBoot
	return
}
