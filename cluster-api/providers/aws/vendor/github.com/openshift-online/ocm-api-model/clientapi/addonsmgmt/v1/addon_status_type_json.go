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

// MarshalAddonStatus writes a value of the 'addon_status' type to the given writer.
func MarshalAddonStatus(object *AddonStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonStatus writes a value of the 'addon_status' type to the given stream.
func WriteAddonStatus(object *AddonStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddonStatusLinkKind)
	} else {
		stream.WriteString(AddonStatusKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon_id")
		stream.WriteString(object.addonId)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("correlation_id")
		stream.WriteString(object.correlationID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.statusConditions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_conditions")
		WriteAddonStatusConditionList(object.statusConditions, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonStatus reads a value of the 'addon_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonStatus(source interface{}) (object *AddonStatus, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonStatus(iterator)
	err = iterator.Error
	return
}

// ReadAddonStatus reads a value of the 'addon_status' type from the given iterator.
func ReadAddonStatus(iterator *jsoniter.Iterator) *AddonStatus {
	object := &AddonStatus{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddonStatusLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "addon_id":
			value := iterator.ReadString()
			object.addonId = value
			object.fieldSet_[3] = true
		case "correlation_id":
			value := iterator.ReadString()
			object.correlationID = value
			object.fieldSet_[4] = true
		case "status_conditions":
			value := ReadAddonStatusConditionList(iterator)
			object.statusConditions = value
			object.fieldSet_[5] = true
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
