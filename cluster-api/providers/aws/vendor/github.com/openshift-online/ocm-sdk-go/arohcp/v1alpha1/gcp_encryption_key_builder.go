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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// GCPEncryptionKeyBuilder contains the data and logic needed to build 'GCP_encryption_key' objects.
//
// GCP Encryption Key for CCS clusters.
type GCPEncryptionKeyBuilder struct {
	bitmap_              uint32
	kmsKeyServiceAccount string
	keyLocation          string
	keyName              string
	keyRing              string
}

// NewGCPEncryptionKey creates a new builder of 'GCP_encryption_key' objects.
func NewGCPEncryptionKey() *GCPEncryptionKeyBuilder {
	return &GCPEncryptionKeyBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPEncryptionKeyBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// KMSKeyServiceAccount sets the value of the 'KMS_key_service_account' attribute to the given value.
func (b *GCPEncryptionKeyBuilder) KMSKeyServiceAccount(value string) *GCPEncryptionKeyBuilder {
	b.kmsKeyServiceAccount = value
	b.bitmap_ |= 1
	return b
}

// KeyLocation sets the value of the 'key_location' attribute to the given value.
func (b *GCPEncryptionKeyBuilder) KeyLocation(value string) *GCPEncryptionKeyBuilder {
	b.keyLocation = value
	b.bitmap_ |= 2
	return b
}

// KeyName sets the value of the 'key_name' attribute to the given value.
func (b *GCPEncryptionKeyBuilder) KeyName(value string) *GCPEncryptionKeyBuilder {
	b.keyName = value
	b.bitmap_ |= 4
	return b
}

// KeyRing sets the value of the 'key_ring' attribute to the given value.
func (b *GCPEncryptionKeyBuilder) KeyRing(value string) *GCPEncryptionKeyBuilder {
	b.keyRing = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPEncryptionKeyBuilder) Copy(object *GCPEncryptionKey) *GCPEncryptionKeyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.kmsKeyServiceAccount = object.kmsKeyServiceAccount
	b.keyLocation = object.keyLocation
	b.keyName = object.keyName
	b.keyRing = object.keyRing
	return b
}

// Build creates a 'GCP_encryption_key' object using the configuration stored in the builder.
func (b *GCPEncryptionKeyBuilder) Build() (object *GCPEncryptionKey, err error) {
	object = new(GCPEncryptionKey)
	object.bitmap_ = b.bitmap_
	object.kmsKeyServiceAccount = b.kmsKeyServiceAccount
	object.keyLocation = b.keyLocation
	object.keyName = b.keyName
	object.keyRing = b.keyRing
	return
}
