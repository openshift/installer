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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalIntegerList writes a list of values of the 'integer' type to
// the given writer.
func MarshalIntegerList(list []int, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteIntegerList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteIntegerList writes a list of value of the 'integer' type to
// the given stream.
func WriteIntegerList(list []int, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		stream.WriteInt(value)
	}
	stream.WriteArrayEnd()
}

// UnmarshalIntegerList reads a list of values of the 'integer' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalIntegerList(source interface{}) (items []int, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadIntegerList(iterator)
	err = iterator.Error
	return
}

// ReadIntegerList reads list of values of the ”integer' type from
// the given iterator.
func ReadIntegerList(iterator *jsoniter.Iterator) []int {
	list := []int{}
	for iterator.ReadArray() {
		item := iterator.ReadInt()
		list = append(list, item)
	}
	return list
}
