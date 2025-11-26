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

// MarshalAzureNodePoolEncryptionAtHostList writes a list of values of the 'azure_node_pool_encryption_at_host' type to
// the given writer.
func MarshalAzureNodePoolEncryptionAtHostList(list []*AzureNodePoolEncryptionAtHost, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureNodePoolEncryptionAtHostList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureNodePoolEncryptionAtHostList writes a list of value of the 'azure_node_pool_encryption_at_host' type to
// the given stream.
func WriteAzureNodePoolEncryptionAtHostList(list []*AzureNodePoolEncryptionAtHost, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAzureNodePoolEncryptionAtHost(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAzureNodePoolEncryptionAtHostList reads a list of values of the 'azure_node_pool_encryption_at_host' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAzureNodePoolEncryptionAtHostList(source interface{}) (items []*AzureNodePoolEncryptionAtHost, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAzureNodePoolEncryptionAtHostList(iterator)
	err = iterator.Error
	return
}

// ReadAzureNodePoolEncryptionAtHostList reads list of values of the ”azure_node_pool_encryption_at_host' type from
// the given iterator.
func ReadAzureNodePoolEncryptionAtHostList(iterator *jsoniter.Iterator) []*AzureNodePoolEncryptionAtHost {
	list := []*AzureNodePoolEncryptionAtHost{}
	for iterator.ReadArray() {
		item := ReadAzureNodePoolEncryptionAtHost(iterator)
		list = append(list, item)
	}
	return list
}
