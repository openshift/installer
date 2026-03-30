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

// MarshalLogForwarderStatus writes a value of the 'log_forwarder_status' type to the given writer.
func MarshalLogForwarderStatus(object *LogForwarderStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderStatus writes a value of the 'log_forwarder_status' type to the given stream.
func WriteLogForwarderStatus(object *LogForwarderStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("message")
		stream.WriteString(object.message)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.resolvedApplications != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resolved_applications")
		WriteStringList(object.resolvedApplications, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(object.state)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarderStatus reads a value of the 'log_forwarder_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarderStatus(source interface{}) (object *LogForwarderStatus, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarderStatus(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderStatus reads a value of the 'log_forwarder_status' type from the given iterator.
func ReadLogForwarderStatus(iterator *jsoniter.Iterator) *LogForwarderStatus {
	object := &LogForwarderStatus{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "message":
			value := iterator.ReadString()
			object.message = value
			object.fieldSet_[0] = true
		case "resolved_applications":
			value := ReadStringList(iterator)
			object.resolvedApplications = value
			object.fieldSet_[1] = true
		case "state":
			value := iterator.ReadString()
			object.state = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
