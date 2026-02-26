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
type AWSFlavourBuilder struct {
	fieldSet_           []bool
	computeInstanceType string
	infraInstanceType   string
	infraVolume         *AWSVolumeBuilder
	masterInstanceType  string
	masterVolume        *AWSVolumeBuilder
	workerVolume        *AWSVolumeBuilder
}

// NewAWSFlavour creates a new builder of 'AWS_flavour' objects.
func NewAWSFlavour() *AWSFlavourBuilder {
	return &AWSFlavourBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSFlavourBuilder) Empty() bool {
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
func (b *AWSFlavourBuilder) ComputeInstanceType(value string) *AWSFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.computeInstanceType = value
	b.fieldSet_[0] = true
	return b
}

// InfraInstanceType sets the value of the 'infra_instance_type' attribute to the given value.
func (b *AWSFlavourBuilder) InfraInstanceType(value string) *AWSFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.infraInstanceType = value
	b.fieldSet_[1] = true
	return b
}

// InfraVolume sets the value of the 'infra_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) InfraVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
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
func (b *AWSFlavourBuilder) MasterInstanceType(value string) *AWSFlavourBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.masterInstanceType = value
	b.fieldSet_[3] = true
	return b
}

// MasterVolume sets the value of the 'master_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) MasterVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
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
// Holds settings for an AWS storage volume.
func (b *AWSFlavourBuilder) WorkerVolume(value *AWSVolumeBuilder) *AWSFlavourBuilder {
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
func (b *AWSFlavourBuilder) Copy(object *AWSFlavour) *AWSFlavourBuilder {
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
