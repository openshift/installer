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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalProvisionShardReferenceList writes a list of values of the 'provision_shard_reference' type to
// the given writer.
func MarshalProvisionShardReferenceList(list []*ProvisionShardReference, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteProvisionShardReferenceList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteProvisionShardReferenceList writes a list of value of the 'provision_shard_reference' type to
// the given stream.
func WriteProvisionShardReferenceList(list []*ProvisionShardReference, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteProvisionShardReference(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalProvisionShardReferenceList reads a list of values of the 'provision_shard_reference' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalProvisionShardReferenceList(source interface{}) (items []*ProvisionShardReference, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadProvisionShardReferenceList(iterator)
	err = iterator.Error
	return
}

// ReadProvisionShardReferenceList reads list of values of the ”provision_shard_reference' type from
// the given iterator.
func ReadProvisionShardReferenceList(iterator *jsoniter.Iterator) []*ProvisionShardReference {
	list := []*ProvisionShardReference{}
	for iterator.ReadArray() {
		item := ReadProvisionShardReference(iterator)
		list = append(list, item)
	}
	return list
}
