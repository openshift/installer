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

// AwsEtcdEncryptionBuilder contains the data and logic needed to build 'aws_etcd_encryption' objects.
//
// Contains the necessary attributes to support etcd encryption for AWS based clusters.
type AwsEtcdEncryptionBuilder struct {
	bitmap_   uint32
	kmsKeyARN string
}

// NewAwsEtcdEncryption creates a new builder of 'aws_etcd_encryption' objects.
func NewAwsEtcdEncryption() *AwsEtcdEncryptionBuilder {
	return &AwsEtcdEncryptionBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AwsEtcdEncryptionBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// KMSKeyARN sets the value of the 'KMS_key_ARN' attribute to the given value.
func (b *AwsEtcdEncryptionBuilder) KMSKeyARN(value string) *AwsEtcdEncryptionBuilder {
	b.kmsKeyARN = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AwsEtcdEncryptionBuilder) Copy(object *AwsEtcdEncryption) *AwsEtcdEncryptionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.kmsKeyARN = object.kmsKeyARN
	return b
}

// Build creates a 'aws_etcd_encryption' object using the configuration stored in the builder.
func (b *AwsEtcdEncryptionBuilder) Build() (object *AwsEtcdEncryption, err error) {
	object = new(AwsEtcdEncryption)
	object.bitmap_ = b.bitmap_
	object.kmsKeyARN = b.kmsKeyARN
	return
}
