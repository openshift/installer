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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalHandoff writes a value of the 'handoff' type to the given writer.
func MarshalHandoff(object *Handoff, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeHandoff(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeHandoff writes a value of the 'handoff' type to the given stream.
func writeHandoff(object *Handoff, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(HandoffLinkKind)
	} else {
		stream.WriteString(HandoffKind)
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
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deleted_at")
		stream.WriteString((object.deletedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.handoffFrom != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("handoff_from")
		writeUser(object.handoffFrom, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.handoffTo != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("handoff_to")
		writeUser(object.handoffTo, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("handoff_type")
		stream.WriteString(object.handoffType)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalHandoff reads a value of the 'handoff' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalHandoff(source interface{}) (object *Handoff, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readHandoff(iterator)
	err = iterator.Error
	return
}

// readHandoff reads a value of the 'handoff' type from the given iterator.
func readHandoff(iterator *jsoniter.Iterator) *Handoff {
	object := &Handoff{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == HandoffLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 8
		case "deleted_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deletedAt = value
			object.bitmap_ |= 16
		case "handoff_from":
			value := readUser(iterator)
			object.handoffFrom = value
			object.bitmap_ |= 32
		case "handoff_to":
			value := readUser(iterator)
			object.handoffTo = value
			object.bitmap_ |= 64
		case "handoff_type":
			value := iterator.ReadString()
			object.handoffType = value
			object.bitmap_ |= 128
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 256
		default:
			iterator.ReadAny()
		}
	}
	return object
}
