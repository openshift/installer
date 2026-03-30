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

// Specification for different classes of nodes inside a flavour.
type GCPFlavourBuilder struct {
	fieldSet_           []bool
	computeInstanceType string
	infraInstanceType   string
	infraVolume         *GCPVolumeBuilder
	masterInstanceType  string
	masterVolume        *GCPVolumeBuilder
	workerVolume        *GCPVolumeBuilder
}

// NewGCPFlavour creates a new builder of 'GCP_flavour' objects.
func NewGCPFlavour() *GCPFlavourBuilder {
	return &GCPFlavourBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPFlavourBuilder) Empty() bool {
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

// ComputeInstanceType sets the value of the 'compute_instance_type' attribute to the given value.
func (b *GCPFlavourBuilder) ComputeInstanceType(value string) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.computeInstanceType = value
	b.fieldSet_[0] = true
	return b
}

// InfraInstanceType sets the value of the 'infra_instance_type' attribute to the given value.
func (b *GCPFlavourBuilder) InfraInstanceType(value string) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.infraInstanceType = value
	b.fieldSet_[1] = true
	return b
}

// InfraVolume sets the value of the 'infra_volume' attribute to the given value.
//
// Holds settings for an GCP storage volume.
func (b *GCPFlavourBuilder) InfraVolume(value *GCPVolumeBuilder) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.infraVolume = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// MasterInstanceType sets the value of the 'master_instance_type' attribute to the given value.
func (b *GCPFlavourBuilder) MasterInstanceType(value string) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.masterInstanceType = value
	b.fieldSet_[3] = true
	return b
}

// MasterVolume sets the value of the 'master_volume' attribute to the given value.
//
// Holds settings for an GCP storage volume.
func (b *GCPFlavourBuilder) MasterVolume(value *GCPVolumeBuilder) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.masterVolume = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// WorkerVolume sets the value of the 'worker_volume' attribute to the given value.
//
// Holds settings for an GCP storage volume.
func (b *GCPFlavourBuilder) WorkerVolume(value *GCPVolumeBuilder) *GCPFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.workerVolume = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPFlavourBuilder) Copy(object *GCPFlavour) *GCPFlavourBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.computeInstanceType = object.computeInstanceType
	b.infraInstanceType = object.infraInstanceType
	if object.infraVolume != nil {
		b.infraVolume = NewGCPVolume().Copy(object.infraVolume)
	} else {
		b.infraVolume = nil
	}
	b.masterInstanceType = object.masterInstanceType
	if object.masterVolume != nil {
		b.masterVolume = NewGCPVolume().Copy(object.masterVolume)
	} else {
		b.masterVolume = nil
	}
	if object.workerVolume != nil {
		b.workerVolume = NewGCPVolume().Copy(object.workerVolume)
	} else {
		b.workerVolume = nil
	}
	return b
}

// Build creates a 'GCP_flavour' object using the configuration stored in the builder.
func (b *GCPFlavourBuilder) Build() (object *GCPFlavour, err error) {
	object = new(GCPFlavour)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.computeInstanceType = b.computeInstanceType
	object.infraInstanceType = b.infraInstanceType
	if b.infraVolume != nil {
		object.infraVolume, err = b.infraVolume.Build()
		if err != nil {
			return
		}
	}
	object.masterInstanceType = b.masterInstanceType
	if b.masterVolume != nil {
		object.masterVolume, err = b.masterVolume.Build()
		if err != nil {
			return
		}
	}
	if b.workerVolume != nil {
		object.workerVolume, err = b.workerVolume.Build()
		if err != nil {
			return
		}
	}
	return
}
