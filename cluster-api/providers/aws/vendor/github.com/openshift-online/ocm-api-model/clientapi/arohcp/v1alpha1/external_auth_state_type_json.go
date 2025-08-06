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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalExternalAuthState writes a value of the 'external_auth_state' type to the given writer.
func MarshalExternalAuthState(object *ExternalAuthState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteExternalAuthState(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteExternalAuthState writes a value of the 'external_auth_state' type to the given stream.
func WriteExternalAuthState(object *ExternalAuthState, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_updated_timestamp")
		stream.WriteString((object.lastUpdatedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.value)
	}
	stream.WriteObjectEnd()
}

// UnmarshalExternalAuthState reads a value of the 'external_auth_state' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalExternalAuthState(source interface{}) (object *ExternalAuthState, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadExternalAuthState(iterator)
	err = iterator.Error
	return
}

// ReadExternalAuthState reads a value of the 'external_auth_state' type from the given iterator.
func ReadExternalAuthState(iterator *jsoniter.Iterator) *ExternalAuthState {
	object := &ExternalAuthState{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "last_updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdatedTimestamp = value
			object.fieldSet_[0] = true
		case "value":
			value := iterator.ReadString()
			object.value = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
