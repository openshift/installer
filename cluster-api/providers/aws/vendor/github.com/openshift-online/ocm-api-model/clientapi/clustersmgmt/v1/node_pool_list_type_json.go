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

// MarshalNodePoolList writes a list of values of the 'node_pool' type to
// the given writer.
func MarshalNodePoolList(list []*NodePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNodePoolList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNodePoolList writes a list of value of the 'node_pool' type to
// the given stream.
func WriteNodePoolList(list []*NodePool, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteNodePool(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalNodePoolList reads a list of values of the 'node_pool' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalNodePoolList(source interface{}) (items []*NodePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadNodePoolList(iterator)
	err = iterator.Error
	return
}

// ReadNodePoolList reads list of values of the ”node_pool' type from
// the given iterator.
func ReadNodePoolList(iterator *jsoniter.Iterator) []*NodePool {
	list := []*NodePool{}
	for iterator.ReadArray() {
		item := ReadNodePool(iterator)
		list = append(list, item)
	}
	return list
}
