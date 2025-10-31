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

// MarshalAzureNodePoolOsDiskList writes a list of values of the 'azure_node_pool_os_disk' type to
// the given writer.
func MarshalAzureNodePoolOsDiskList(list []*AzureNodePoolOsDisk, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureNodePoolOsDiskList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureNodePoolOsDiskList writes a list of value of the 'azure_node_pool_os_disk' type to
// the given stream.
func WriteAzureNodePoolOsDiskList(list []*AzureNodePoolOsDisk, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAzureNodePoolOsDisk(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAzureNodePoolOsDiskList reads a list of values of the 'azure_node_pool_os_disk' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAzureNodePoolOsDiskList(source interface{}) (items []*AzureNodePoolOsDisk, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAzureNodePoolOsDiskList(iterator)
	err = iterator.Error
	return
}

// ReadAzureNodePoolOsDiskList reads list of values of the ”azure_node_pool_os_disk' type from
// the given iterator.
func ReadAzureNodePoolOsDiskList(iterator *jsoniter.Iterator) []*AzureNodePoolOsDisk {
	list := []*AzureNodePoolOsDisk{}
	for iterator.ReadArray() {
		item := ReadAzureNodePoolOsDisk(iterator)
		list = append(list, item)
	}
	return list
}
