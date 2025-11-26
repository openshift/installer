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

// MarshalCIDRBlockAccessList writes a list of values of the 'CIDR_block_access' type to
// the given writer.
func MarshalCIDRBlockAccessList(list []*CIDRBlockAccess, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCIDRBlockAccessList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCIDRBlockAccessList writes a list of value of the 'CIDR_block_access' type to
// the given stream.
func WriteCIDRBlockAccessList(list []*CIDRBlockAccess, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteCIDRBlockAccess(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalCIDRBlockAccessList reads a list of values of the 'CIDR_block_access' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalCIDRBlockAccessList(source interface{}) (items []*CIDRBlockAccess, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadCIDRBlockAccessList(iterator)
	err = iterator.Error
	return
}

// ReadCIDRBlockAccessList reads list of values of the ”CIDR_block_access' type from
// the given iterator.
func ReadCIDRBlockAccessList(iterator *jsoniter.Iterator) []*CIDRBlockAccess {
	list := []*CIDRBlockAccess{}
	for iterator.ReadArray() {
		item := ReadCIDRBlockAccess(iterator)
		list = append(list, item)
	}
	return list
}
