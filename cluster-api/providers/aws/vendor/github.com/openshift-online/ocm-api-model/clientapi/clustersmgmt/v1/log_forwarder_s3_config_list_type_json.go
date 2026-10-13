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

// MarshalLogForwarderS3ConfigList writes a list of values of the 'log_forwarder_S3_config' type to
// the given writer.
func MarshalLogForwarderS3ConfigList(list []*LogForwarderS3Config, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderS3ConfigList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderS3ConfigList writes a list of value of the 'log_forwarder_S3_config' type to
// the given stream.
func WriteLogForwarderS3ConfigList(list []*LogForwarderS3Config, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteLogForwarderS3Config(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalLogForwarderS3ConfigList reads a list of values of the 'log_forwarder_S3_config' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalLogForwarderS3ConfigList(source interface{}) (items []*LogForwarderS3Config, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadLogForwarderS3ConfigList(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderS3ConfigList reads list of values of the ”log_forwarder_S3_config' type from
// the given iterator.
func ReadLogForwarderS3ConfigList(iterator *jsoniter.Iterator) []*LogForwarderS3Config {
	list := []*LogForwarderS3Config{}
	for iterator.ReadArray() {
		item := ReadLogForwarderS3Config(iterator)
		list = append(list, item)
	}
	return list
}
