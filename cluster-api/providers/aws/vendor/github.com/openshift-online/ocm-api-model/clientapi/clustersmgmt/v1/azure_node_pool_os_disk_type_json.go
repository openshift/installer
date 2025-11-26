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

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAzureNodePoolOsDisk writes a value of the 'azure_node_pool_os_disk' type to the given writer.
func MarshalAzureNodePoolOsDisk(object *AzureNodePoolOsDisk, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureNodePoolOsDisk(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureNodePoolOsDisk writes a value of the 'azure_node_pool_os_disk' type to the given stream.
func WriteAzureNodePoolOsDisk(object *AzureNodePoolOsDisk, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("persistence")
		stream.WriteString(object.persistence)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("size_gibibytes")
		stream.WriteInt(object.sizeGibibytes)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sse_encryption_set_resource_id")
		stream.WriteString(object.sseEncryptionSetResourceId)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("storage_account_type")
		stream.WriteString(object.storageAccountType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureNodePoolOsDisk reads a value of the 'azure_node_pool_os_disk' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureNodePoolOsDisk(source interface{}) (object *AzureNodePoolOsDisk, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureNodePoolOsDisk(iterator)
	err = iterator.Error
	return
}

// ReadAzureNodePoolOsDisk reads a value of the 'azure_node_pool_os_disk' type from the given iterator.
func ReadAzureNodePoolOsDisk(iterator *jsoniter.Iterator) *AzureNodePoolOsDisk {
	object := &AzureNodePoolOsDisk{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "persistence":
			value := iterator.ReadString()
			object.persistence = value
			object.fieldSet_[0] = true
		case "size_gibibytes":
			value := iterator.ReadInt()
			object.sizeGibibytes = value
			object.fieldSet_[1] = true
		case "sse_encryption_set_resource_id":
			value := iterator.ReadString()
			object.sseEncryptionSetResourceId = value
			object.fieldSet_[2] = true
		case "storage_account_type":
			value := iterator.ReadString()
			object.storageAccountType = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
