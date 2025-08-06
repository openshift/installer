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

// MarshalWifConfigStatus writes a value of the 'wif_config_status' type to the given writer.
func MarshalWifConfigStatus(object *WifConfigStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifConfigStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifConfigStatus writes a value of the 'wif_config_status' type to the given stream.
func WriteWifConfigStatus(object *WifConfigStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("configured")
		stream.WriteBool(object.configured)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifConfigStatus reads a value of the 'wif_config_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifConfigStatus(source interface{}) (object *WifConfigStatus, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifConfigStatus(iterator)
	err = iterator.Error
	return
}

// ReadWifConfigStatus reads a value of the 'wif_config_status' type from the given iterator.
func ReadWifConfigStatus(iterator *jsoniter.Iterator) *WifConfigStatus {
	object := &WifConfigStatus{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "configured":
			value := iterator.ReadBool()
			object.configured = value
			object.fieldSet_[0] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
