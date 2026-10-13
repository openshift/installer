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

// MarshalAzureNodePool writes a value of the 'azure_node_pool' type to the given writer.
func MarshalAzureNodePool(object *AzureNodePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureNodePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureNodePool writes a value of the 'azure_node_pool' type to the given stream.
func WriteAzureNodePool(object *AzureNodePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vm_size")
		stream.WriteString(object.vmSize)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.encryptionAtHost != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("encryption_at_host")
		WriteAzureNodePoolEncryptionAtHost(object.encryptionAtHost, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.osDisk != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("os_disk")
		WriteAzureNodePoolOsDisk(object.osDisk, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureNodePool reads a value of the 'azure_node_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureNodePool(source interface{}) (object *AzureNodePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureNodePool(iterator)
	err = iterator.Error
	return
}

// ReadAzureNodePool reads a value of the 'azure_node_pool' type from the given iterator.
func ReadAzureNodePool(iterator *jsoniter.Iterator) *AzureNodePool {
	object := &AzureNodePool{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "vm_size":
			value := iterator.ReadString()
			object.vmSize = value
			object.fieldSet_[0] = true
		case "encryption_at_host":
			value := ReadAzureNodePoolEncryptionAtHost(iterator)
			object.encryptionAtHost = value
			object.fieldSet_[1] = true
		case "os_disk":
			value := ReadAzureNodePoolOsDisk(iterator)
			object.osDisk = value
			object.fieldSet_[2] = true
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
