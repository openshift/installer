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

// AWSSTSPolicyBuilder contains the data and logic needed to build 'AWSSTS_policy' objects.
//
// Representation of an sts policies for rosa cluster
type AWSSTSPolicyBuilder struct {
	bitmap_ uint32
	arn     string
	id      string
	details string
	type_   string
}

// NewAWSSTSPolicy creates a new builder of 'AWSSTS_policy' objects.
func NewAWSSTSPolicy() *AWSSTSPolicyBuilder {
	return &AWSSTSPolicyBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSSTSPolicyBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ARN sets the value of the 'ARN' attribute to the given value.
func (b *AWSSTSPolicyBuilder) ARN(value string) *AWSSTSPolicyBuilder {
	b.arn = value
	b.bitmap_ |= 1
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AWSSTSPolicyBuilder) ID(value string) *AWSSTSPolicyBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// Details sets the value of the 'details' attribute to the given value.
func (b *AWSSTSPolicyBuilder) Details(value string) *AWSSTSPolicyBuilder {
	b.details = value
	b.bitmap_ |= 4
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *AWSSTSPolicyBuilder) Type(value string) *AWSSTSPolicyBuilder {
	b.type_ = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSSTSPolicyBuilder) Copy(object *AWSSTSPolicy) *AWSSTSPolicyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.arn = object.arn
	b.id = object.id
	b.details = object.details
	b.type_ = object.type_
	return b
}

// Build creates a 'AWSSTS_policy' object using the configuration stored in the builder.
func (b *AWSSTSPolicyBuilder) Build() (object *AWSSTSPolicy, err error) {
	object = new(AWSSTSPolicy)
	object.bitmap_ = b.bitmap_
	object.arn = b.arn
	object.id = b.id
	object.details = b.details
	object.type_ = b.type_
	return
}
