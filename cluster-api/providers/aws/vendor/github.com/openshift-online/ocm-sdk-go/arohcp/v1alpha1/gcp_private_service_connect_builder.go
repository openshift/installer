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

// GcpPrivateServiceConnectBuilder contains the data and logic needed to build 'gcp_private_service_connect' objects.
//
// Google cloud platform private service connect configuration of a cluster.
type GcpPrivateServiceConnectBuilder struct {
	bitmap_                 uint32
	serviceAttachmentSubnet string
}

// NewGcpPrivateServiceConnect creates a new builder of 'gcp_private_service_connect' objects.
func NewGcpPrivateServiceConnect() *GcpPrivateServiceConnectBuilder {
	return &GcpPrivateServiceConnectBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GcpPrivateServiceConnectBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ServiceAttachmentSubnet sets the value of the 'service_attachment_subnet' attribute to the given value.
func (b *GcpPrivateServiceConnectBuilder) ServiceAttachmentSubnet(value string) *GcpPrivateServiceConnectBuilder {
	b.serviceAttachmentSubnet = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GcpPrivateServiceConnectBuilder) Copy(object *GcpPrivateServiceConnect) *GcpPrivateServiceConnectBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.serviceAttachmentSubnet = object.serviceAttachmentSubnet
	return b
}

// Build creates a 'gcp_private_service_connect' object using the configuration stored in the builder.
func (b *GcpPrivateServiceConnectBuilder) Build() (object *GcpPrivateServiceConnect, err error) {
	object = new(GcpPrivateServiceConnect)
	object.bitmap_ = b.bitmap_
	object.serviceAttachmentSubnet = b.serviceAttachmentSubnet
	return
}
