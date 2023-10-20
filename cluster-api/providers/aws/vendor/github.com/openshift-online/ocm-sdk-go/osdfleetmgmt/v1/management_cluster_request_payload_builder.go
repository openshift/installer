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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

// ManagementClusterRequestPayloadBuilder contains the data and logic needed to build 'management_cluster_request_payload' objects.
type ManagementClusterRequestPayloadBuilder struct {
	bitmap_                                                                    uint32
	service_cluster_idService_cluster_idService_cluster_idService_cluster_idId string
}

// NewManagementClusterRequestPayload creates a new builder of 'management_cluster_request_payload' objects.
func NewManagementClusterRequestPayload() *ManagementClusterRequestPayloadBuilder {
	return &ManagementClusterRequestPayloadBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManagementClusterRequestPayloadBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Service_cluster_idService_cluster_idService_cluster_idService_cluster_idId sets the value of the 'service_cluster_id_service_cluster_id_service_cluster_id_service_cluster_id_id' attribute to the given value.
func (b *ManagementClusterRequestPayloadBuilder) Service_cluster_idService_cluster_idService_cluster_idService_cluster_idId(value string) *ManagementClusterRequestPayloadBuilder {
	b.service_cluster_idService_cluster_idService_cluster_idService_cluster_idId = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManagementClusterRequestPayloadBuilder) Copy(object *ManagementClusterRequestPayload) *ManagementClusterRequestPayloadBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.service_cluster_idService_cluster_idService_cluster_idService_cluster_idId = object.service_cluster_idService_cluster_idService_cluster_idService_cluster_idId
	return b
}

// Build creates a 'management_cluster_request_payload' object using the configuration stored in the builder.
func (b *ManagementClusterRequestPayloadBuilder) Build() (object *ManagementClusterRequestPayload, err error) {
	object = new(ManagementClusterRequestPayload)
	object.bitmap_ = b.bitmap_
	object.service_cluster_idService_cluster_idService_cluster_idService_cluster_idId = b.service_cluster_idService_cluster_idService_cluster_idService_cluster_idId
	return
}
