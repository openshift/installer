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

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAzureControlPlaneManagedIdentity writes a value of the 'azure_control_plane_managed_identity' type to the given writer.
func MarshalAzureControlPlaneManagedIdentity(object *AzureControlPlaneManagedIdentity, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAzureControlPlaneManagedIdentity(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAzureControlPlaneManagedIdentity writes a value of the 'azure_control_plane_managed_identity' type to the given stream.
func writeAzureControlPlaneManagedIdentity(object *AzureControlPlaneManagedIdentity, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principal_id")
		stream.WriteString(object.principalID)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_id")
		stream.WriteString(object.resourceID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureControlPlaneManagedIdentity reads a value of the 'azure_control_plane_managed_identity' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureControlPlaneManagedIdentity(source interface{}) (object *AzureControlPlaneManagedIdentity, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAzureControlPlaneManagedIdentity(iterator)
	err = iterator.Error
	return
}

// readAzureControlPlaneManagedIdentity reads a value of the 'azure_control_plane_managed_identity' type from the given iterator.
func readAzureControlPlaneManagedIdentity(iterator *jsoniter.Iterator) *AzureControlPlaneManagedIdentity {
	object := &AzureControlPlaneManagedIdentity{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.bitmap_ |= 1
		case "principal_id":
			value := iterator.ReadString()
			object.principalID = value
			object.bitmap_ |= 2
		case "resource_id":
			value := iterator.ReadString()
			object.resourceID = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
