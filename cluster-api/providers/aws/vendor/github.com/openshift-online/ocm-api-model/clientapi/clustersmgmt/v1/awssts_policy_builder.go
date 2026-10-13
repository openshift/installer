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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Representation of an sts policies for rosa cluster
type AWSSTSPolicyBuilder struct {
	fieldSet_ []bool
	arn       string
	id        string
	details   string
	type_     string
}

// NewAWSSTSPolicy creates a new builder of 'AWSSTS_policy' objects.
func NewAWSSTSPolicy() *AWSSTSPolicyBuilder {
	return &AWSSTSPolicyBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSPolicyBuilder) Empty() bool {
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

// ARN sets the value of the 'ARN' attribute to the given value.
func (b *AWSSTSPolicyBuilder) ARN(value string) *AWSSTSPolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.arn = value
	b.fieldSet_[0] = true
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AWSSTSPolicyBuilder) ID(value string) *AWSSTSPolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// Details sets the value of the 'details' attribute to the given value.
func (b *AWSSTSPolicyBuilder) Details(value string) *AWSSTSPolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.details = value
	b.fieldSet_[2] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *AWSSTSPolicyBuilder) Type(value string) *AWSSTSPolicyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.type_ = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSPolicyBuilder) Copy(object *AWSSTSPolicy) *AWSSTSPolicyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.arn = object.arn
	b.id = object.id
	b.details = object.details
	b.type_ = object.type_
	return b
}

// Build creates a 'AWSSTS_policy' object using the configuration stored in the builder.
func (b *AWSSTSPolicyBuilder) Build() (object *AWSSTSPolicy, err error) {
	object = new(AWSSTSPolicy)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.arn = b.arn
	object.id = b.id
	object.details = b.details
	object.type_ = b.type_
	return
}
