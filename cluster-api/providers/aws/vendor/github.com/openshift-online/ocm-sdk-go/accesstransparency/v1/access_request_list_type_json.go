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

// MarshalAccessRequestList writes a list of values of the 'access_request' type to
// the given writer.
func MarshalAccessRequestList(list []*AccessRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccessRequestList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccessRequestList writes a list of value of the 'access_request' type to
// the given stream.
func WriteAccessRequestList(list []*AccessRequest, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAccessRequest(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAccessRequestList reads a list of values of the 'access_request' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAccessRequestList(source interface{}) (items []*AccessRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAccessRequestList(iterator)
	err = iterator.Error
	return
}

// ReadAccessRequestList reads list of values of the ”access_request' type from
// the given iterator.
func ReadAccessRequestList(iterator *jsoniter.Iterator) []*AccessRequest {
	list := []*AccessRequest{}
	for iterator.ReadArray() {
		item := ReadAccessRequest(iterator)
		list = append(list, item)
	}
	return list
}
