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

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAzureServiceManagedIdentity writes a value of the 'azure_service_managed_identity' type to the given writer.
func MarshalAzureServiceManagedIdentity(object *AzureServiceManagedIdentity, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureServiceManagedIdentity(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureServiceManagedIdentity writes a value of the 'azure_service_managed_identity' type to the given stream.
func WriteAzureServiceManagedIdentity(object *AzureServiceManagedIdentity, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principal_id")
		stream.WriteString(object.principalID)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_id")
		stream.WriteString(object.resourceID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureServiceManagedIdentity reads a value of the 'azure_service_managed_identity' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureServiceManagedIdentity(source interface{}) (object *AzureServiceManagedIdentity, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureServiceManagedIdentity(iterator)
	err = iterator.Error
	return
}

// ReadAzureServiceManagedIdentity reads a value of the 'azure_service_managed_identity' type from the given iterator.
func ReadAzureServiceManagedIdentity(iterator *jsoniter.Iterator) *AzureServiceManagedIdentity {
	object := &AzureServiceManagedIdentity{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.fieldSet_[0] = true
		case "principal_id":
			value := iterator.ReadString()
			object.principalID = value
			object.fieldSet_[1] = true
		case "resource_id":
			value := iterator.ReadString()
			object.resourceID = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
