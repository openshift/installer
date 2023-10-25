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

// AWSFlavourBuilder contains the data and logic needed to build 'AWS_flavour' objects.
//
// Specification for different classes of nodes inside a flavour.
type AWSFlavourBuilder struct {
	bitmap_             uint32
	computeInstanceType string
	infraInstanceType   string
	infraVolume         *AWSVolumeBuilder
	masterInstanceType  string
	masterVolume        *AWSVolumeBuilder
	workerVolume        *AWSVolumeBuilder
}

// NewAWSFlavour creates a new builder of 'AWS_flavour' objects.
func NewAWSFlavour() *AWSFlavourBuilder {
	return &AWSFlavourBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSFlavourBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ComputeInstanceType sets the value of the 'compute_instance_type' attribute to the given value.
func (b *AWSFlavourBuilder) ComputeInstanceType(value string) *AWSFlavourBuilder {
	b.computeInstanceType = value
	b.bitmap_ |= 1
	return b
}

// InfraInstanceType sets the value of the 'infra_instance_type' attribute to the given value.
func (b *AWSFlavourBuilder) InfraInstanceType(value string) *AWSFlavourBuilder {
	b.infraInstanceType = value
	b.bitmap_ |= 2
	return b
}

// InfraVolume sets the value of the 'infra_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) InfraVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
	b.infraVolume = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// MasterInstanceType sets the value of the 'master_instance_type' attribute to the given value.
func (b *AWSFlavourBuilder) MasterInstanceType(value string) *AWSFlavourBuilder {
	b.masterInstanceType = value
	b.bitmap_ |= 8
	return b
}

// MasterVolume sets the value of the 'master_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) MasterVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
	b.masterVolume = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// WorkerVolume sets the value of the 'worker_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) WorkerVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
	b.workerVolume = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSFlavourBuilder) Copy(object *AWSFlavour) *AWSFlavourBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.computeInstanceType = object.computeInstanceType
	b.infraInstanceType = object.infraInstanceType
	if object.infraVolume != nil {
		b.infraVolume = NewAWSVolume().Copy(object.infraVolume)
	} else {
		b.infraVolume = nil
	}
	b.masterInstanceType = object.masterInstanceType
	if object.masterVolume != nil {
		b.masterVolume = NewAWSVolume().Copy(object.masterVolume)
	} else {
		b.masterVolume = nil
	}
	if object.workerVolume != nil {
		b.workerVolume = NewAWSVolume().Copy(object.workerVolume)
	} else {
		b.workerVolume = nil
	}
	return b
}

// Build creates a 'AWS_flavour' object using the configuration stored in the builder.
func (b *AWSFlavourBuilder) Build() (object *AWSFlavour, err error) {
	object = new(AWSFlavour)
	object.bitmap_ = b.bitmap_
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
