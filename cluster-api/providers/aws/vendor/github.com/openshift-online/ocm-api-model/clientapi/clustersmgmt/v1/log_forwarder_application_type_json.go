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

// MarshalLogForwarderApplication writes a value of the 'log_forwarder_application' type to the given writer.
func MarshalLogForwarderApplication(object *LogForwarderApplication, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderApplication(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderApplication writes a value of the 'log_forwarder_application' type to the given stream.
func WriteLogForwarderApplication(object *LogForwarderApplication, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarderApplication reads a value of the 'log_forwarder_application' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarderApplication(source interface{}) (object *LogForwarderApplication, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarderApplication(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderApplication reads a value of the 'log_forwarder_application' type from the given iterator.
func ReadLogForwarderApplication(iterator *jsoniter.Iterator) *LogForwarderApplication {
	object := &LogForwarderApplication{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[0] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
