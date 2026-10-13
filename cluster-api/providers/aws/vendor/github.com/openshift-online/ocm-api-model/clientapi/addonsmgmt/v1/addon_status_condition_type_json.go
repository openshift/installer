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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddonStatusCondition writes a value of the 'addon_status_condition' type to the given writer.
func MarshalAddonStatusCondition(object *AddonStatusCondition, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonStatusCondition(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonStatusCondition writes a value of the 'addon_status_condition' type to the given stream.
func WriteAddonStatusCondition(object *AddonStatusCondition, stream *jsoniter.Stream) {
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
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reason")
		stream.WriteString(object.reason)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_type")
		stream.WriteString(string(object.statusType))
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_value")
		stream.WriteString(string(object.statusValue))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonStatusCondition reads a value of the 'addon_status_condition' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonStatusCondition(source interface{}) (object *AddonStatusCondition, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonStatusCondition(iterator)
	err = iterator.Error
	return
}

// ReadAddonStatusCondition reads a value of the 'addon_status_condition' type from the given iterator.
func ReadAddonStatusCondition(iterator *jsoniter.Iterator) *AddonStatusCondition {
	object := &AddonStatusCondition{
		fieldSet_: make([]bool, 4),
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
		case "reason":
			value := iterator.ReadString()
			object.reason = value
			object.fieldSet_[1] = true
		case "status_type":
			text := iterator.ReadString()
			value := AddonStatusConditionType(text)
			object.statusType = value
			object.fieldSet_[2] = true
		case "status_value":
			text := iterator.ReadString()
			value := AddonStatusConditionValue(text)
			object.statusValue = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
