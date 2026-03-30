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

// MarshalLogForwarderGroup writes a value of the 'log_forwarder_group' type to the given writer.
func MarshalLogForwarderGroup(object *LogForwarderGroup, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderGroup(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderGroup writes a value of the 'log_forwarder_group' type to the given stream.
func WriteLogForwarderGroup(object *LogForwarderGroup, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarderGroup reads a value of the 'log_forwarder_group' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarderGroup(source interface{}) (object *LogForwarderGroup, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarderGroup(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderGroup reads a value of the 'log_forwarder_group' type from the given iterator.
func ReadLogForwarderGroup(iterator *jsoniter.Iterator) *LogForwarderGroup {
	object := &LogForwarderGroup{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[0] = true
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
