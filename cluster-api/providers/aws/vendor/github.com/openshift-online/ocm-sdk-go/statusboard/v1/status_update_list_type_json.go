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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalStatusUpdateList writes a list of values of the 'status_update' type to
// the given writer.
func MarshalStatusUpdateList(list []*StatusUpdate, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteStatusUpdateList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteStatusUpdateList writes a list of value of the 'status_update' type to
// the given stream.
func WriteStatusUpdateList(list []*StatusUpdate, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteStatusUpdate(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalStatusUpdateList reads a list of values of the 'status_update' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalStatusUpdateList(source interface{}) (items []*StatusUpdate, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadStatusUpdateList(iterator)
	err = iterator.Error
	return
}

// ReadStatusUpdateList reads list of values of the ”status_update' type from
// the given iterator.
func ReadStatusUpdateList(iterator *jsoniter.Iterator) []*StatusUpdate {
	list := []*StatusUpdate{}
	for iterator.ReadArray() {
		item := ReadStatusUpdate(iterator)
		list = append(list, item)
	}
	return list
}
