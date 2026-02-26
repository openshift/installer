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

// MarshalLogForwarderGroupVersionList writes a list of values of the 'log_forwarder_group_version' type to
// the given writer.
func MarshalLogForwarderGroupVersionList(list []*LogForwarderGroupVersion, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderGroupVersionList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderGroupVersionList writes a list of value of the 'log_forwarder_group_version' type to
// the given stream.
func WriteLogForwarderGroupVersionList(list []*LogForwarderGroupVersion, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteLogForwarderGroupVersion(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalLogForwarderGroupVersionList reads a list of values of the 'log_forwarder_group_version' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalLogForwarderGroupVersionList(source interface{}) (items []*LogForwarderGroupVersion, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadLogForwarderGroupVersionList(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderGroupVersionList reads list of values of the ”log_forwarder_group_version' type from
// the given iterator.
func ReadLogForwarderGroupVersionList(iterator *jsoniter.Iterator) []*LogForwarderGroupVersion {
	list := []*LogForwarderGroupVersion{}
	for iterator.ReadArray() {
		item := ReadLogForwarderGroupVersion(iterator)
		list = append(list, item)
	}
	return list
}
