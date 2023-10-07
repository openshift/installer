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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonStatus writes a value of the 'addon_status' type to the given writer.
func MarshalAddonStatus(object *AddonStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonStatus writes a value of the 'addon_status' type to the given stream.
func writeAddonStatus(object *AddonStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AddonStatusLinkKind)
	} else {
		stream.WriteString(AddonStatusKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon_id")
		stream.WriteString(object.addonId)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("correlation_id")
		stream.WriteString(object.correlationID)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.statusConditions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_conditions")
		writeAddonStatusConditionList(object.statusConditions, stream)
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
	object = readAddonStatus(iterator)
	err = iterator.Error
	return
}

// readAddonStatus reads a value of the 'addon_status' type from the given iterator.
func readAddonStatus(iterator *jsoniter.Iterator) *AddonStatus {
	object := &AddonStatus{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddonStatusLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "addon_id":
			value := iterator.ReadString()
			object.addonId = value
			object.bitmap_ |= 8
		case "correlation_id":
			value := iterator.ReadString()
			object.correlationID = value
			object.bitmap_ |= 16
		case "status_conditions":
			value := readAddonStatusConditionList(iterator)
			object.statusConditions = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
