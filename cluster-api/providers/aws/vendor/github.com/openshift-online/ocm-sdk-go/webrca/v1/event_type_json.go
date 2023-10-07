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

// MarshalEvent writes a value of the 'event' type to the given writer.
func MarshalEvent(object *Event, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeEvent(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeEvent writes a value of the 'event' type to the given stream.
func writeEvent(object *Event, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(EventLinkKind)
	} else {
		stream.WriteString(EventKind)
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
	present_ = object.bitmap_&16 != 0 && object.creator != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creator")
		writeUser(object.creator, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deleted_at")
		stream.WriteString((object.deletedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.escalation != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("escalation")
		writeEscalation(object.escalation, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_type")
		stream.WriteString(object.eventType)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_reference_url")
		stream.WriteString(object.externalReferenceUrl)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.followUp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("follow_up")
		writeFollowUp(object.followUp, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.followUpChange != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("follow_up_change")
		writeFollowUpChange(object.followUpChange, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0 && object.handoff != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("handoff")
		writeHandoff(object.handoff, stream)
		count++
	}
	present_ = object.bitmap_&4096 != 0 && object.incident != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("incident")
		writeIncident(object.incident, stream)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("note")
		stream.WriteString(object.note)
		count++
	}
	present_ = object.bitmap_&16384 != 0 && object.statusChange != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_change")
		writeStatusChange(object.statusChange, stream)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalEvent reads a value of the 'event' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalEvent(source interface{}) (object *Event, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readEvent(iterator)
	err = iterator.Error
	return
}

// readEvent reads a value of the 'event' type from the given iterator.
func readEvent(iterator *jsoniter.Iterator) *Event {
	object := &Event{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == EventLinkKind {
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
		case "creator":
			value := readUser(iterator)
			object.creator = value
			object.bitmap_ |= 16
		case "deleted_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deletedAt = value
			object.bitmap_ |= 32
		case "escalation":
			value := readEscalation(iterator)
			object.escalation = value
			object.bitmap_ |= 64
		case "event_type":
			value := iterator.ReadString()
			object.eventType = value
			object.bitmap_ |= 128
		case "external_reference_url":
			value := iterator.ReadString()
			object.externalReferenceUrl = value
			object.bitmap_ |= 256
		case "follow_up":
			value := readFollowUp(iterator)
			object.followUp = value
			object.bitmap_ |= 512
		case "follow_up_change":
			value := readFollowUpChange(iterator)
			object.followUpChange = value
			object.bitmap_ |= 1024
		case "handoff":
			value := readHandoff(iterator)
			object.handoff = value
			object.bitmap_ |= 2048
		case "incident":
			value := readIncident(iterator)
			object.incident = value
			object.bitmap_ |= 4096
		case "note":
			value := iterator.ReadString()
			object.note = value
			object.bitmap_ |= 8192
		case "status_change":
			value := readStatusChange(iterator)
			object.statusChange = value
			object.bitmap_ |= 16384
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
