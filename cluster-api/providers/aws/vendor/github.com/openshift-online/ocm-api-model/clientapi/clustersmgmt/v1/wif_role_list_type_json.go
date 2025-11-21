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

// MarshalWifRoleList writes a list of values of the 'wif_role' type to
// the given writer.
func MarshalWifRoleList(list []*WifRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifRoleList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifRoleList writes a list of value of the 'wif_role' type to
// the given stream.
func WriteWifRoleList(list []*WifRole, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteWifRole(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalWifRoleList reads a list of values of the 'wif_role' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalWifRoleList(source interface{}) (items []*WifRole, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadWifRoleList(iterator)
	err = iterator.Error
	return
}

// ReadWifRoleList reads list of values of the ”wif_role' type from
// the given iterator.
func ReadWifRoleList(iterator *jsoniter.Iterator) []*WifRole {
	list := []*WifRole{}
	for iterator.ReadArray() {
		item := ReadWifRole(iterator)
		list = append(list, item)
	}
	return list
}
