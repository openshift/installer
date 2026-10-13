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

// MarshalLogForwarderGroupList writes a list of values of the 'log_forwarder_group' type to
// the given writer.
func MarshalLogForwarderGroupList(list []*LogForwarderGroup, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderGroupList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderGroupList writes a list of value of the 'log_forwarder_group' type to
// the given stream.
func WriteLogForwarderGroupList(list []*LogForwarderGroup, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteLogForwarderGroup(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalLogForwarderGroupList reads a list of values of the 'log_forwarder_group' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalLogForwarderGroupList(source interface{}) (items []*LogForwarderGroup, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadLogForwarderGroupList(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderGroupList reads list of values of the ”log_forwarder_group' type from
// the given iterator.
func ReadLogForwarderGroupList(iterator *jsoniter.Iterator) []*LogForwarderGroup {
	list := []*LogForwarderGroup{}
	for iterator.ReadArray() {
		item := ReadLogForwarderGroup(iterator)
		list = append(list, item)
	}
	return list
}
