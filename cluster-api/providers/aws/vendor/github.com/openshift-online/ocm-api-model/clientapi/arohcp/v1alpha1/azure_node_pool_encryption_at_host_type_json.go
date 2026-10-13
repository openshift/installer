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

// MarshalAzureNodePoolEncryptionAtHost writes a value of the 'azure_node_pool_encryption_at_host' type to the given writer.
func MarshalAzureNodePoolEncryptionAtHost(object *AzureNodePoolEncryptionAtHost, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureNodePoolEncryptionAtHost(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureNodePoolEncryptionAtHost writes a value of the 'azure_node_pool_encryption_at_host' type to the given stream.
func WriteAzureNodePoolEncryptionAtHost(object *AzureNodePoolEncryptionAtHost, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(object.state)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureNodePoolEncryptionAtHost reads a value of the 'azure_node_pool_encryption_at_host' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureNodePoolEncryptionAtHost(source interface{}) (object *AzureNodePoolEncryptionAtHost, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureNodePoolEncryptionAtHost(iterator)
	err = iterator.Error
	return
}

// ReadAzureNodePoolEncryptionAtHost reads a value of the 'azure_node_pool_encryption_at_host' type from the given iterator.
func ReadAzureNodePoolEncryptionAtHost(iterator *jsoniter.Iterator) *AzureNodePoolEncryptionAtHost {
	object := &AzureNodePoolEncryptionAtHost{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "state":
			value := iterator.ReadString()
			object.state = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
