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

// MarshalLogForwarderApplicationList writes a list of values of the 'log_forwarder_application' type to
// the given writer.
func MarshalLogForwarderApplicationList(list []*LogForwarderApplication, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderApplicationList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderApplicationList writes a list of value of the 'log_forwarder_application' type to
// the given stream.
func WriteLogForwarderApplicationList(list []*LogForwarderApplication, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteLogForwarderApplication(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalLogForwarderApplicationList reads a list of values of the 'log_forwarder_application' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalLogForwarderApplicationList(source interface{}) (items []*LogForwarderApplication, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadLogForwarderApplicationList(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderApplicationList reads list of values of the ”log_forwarder_application' type from
// the given iterator.
func ReadLogForwarderApplicationList(iterator *jsoniter.Iterator) []*LogForwarderApplication {
	list := []*LogForwarderApplication{}
	for iterator.ReadArray() {
		item := ReadLogForwarderApplication(iterator)
		list = append(list, item)
	}
	return list
}
