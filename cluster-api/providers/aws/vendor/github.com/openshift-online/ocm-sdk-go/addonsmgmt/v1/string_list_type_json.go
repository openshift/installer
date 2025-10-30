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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalStringList writes a list of values of the 'string' type to
// the given writer.
func MarshalStringList(list []string, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteStringList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteStringList writes a list of value of the 'string' type to
// the given stream.
func WriteStringList(list []string, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		stream.WriteString(value)
	}
	stream.WriteArrayEnd()
}

// UnmarshalStringList reads a list of values of the 'string' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalStringList(source interface{}) (items []string, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadStringList(iterator)
	err = iterator.Error
	return
}

// ReadStringList reads list of values of the ”string' type from
// the given iterator.
func ReadStringList(iterator *jsoniter.Iterator) []string {
	list := []string{}
	for iterator.ReadArray() {
		item := iterator.ReadString()
		list = append(list, item)
	}
	return list
}
