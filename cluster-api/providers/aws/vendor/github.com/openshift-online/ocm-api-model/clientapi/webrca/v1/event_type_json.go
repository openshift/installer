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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/webrca/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalEvent writes a value of the 'event' type to the given writer.
func MarshalEvent(object *Event, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteEvent(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteEvent writes a value of the 'event' type to the given stream.
func WriteEvent(object *Event, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(EventLinkKind)
	} else {
		stream.WriteString(EventKind)
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
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.creator != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creator")
		WriteUser(object.creator, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deleted_at")
		stream.WriteString((object.deletedAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.escalation != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("escalation")
		WriteEscalation(object.escalation, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_type")
		stream.WriteString(object.eventType)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_reference_url")
		stream.WriteString(object.externalReferenceUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.followUp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("follow_up")
		WriteFollowUp(object.followUp, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.followUpChange != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("follow_up_change")
		WriteFollowUpChange(object.followUpChange, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.handoff != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("handoff")
		WriteHandoff(object.handoff, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.incident != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("incident")
		WriteIncident(object.incident, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("note")
		stream.WriteString(object.note)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14] && object.statusChange != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_change")
		WriteStatusChange(object.statusChange, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
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
	object = ReadEvent(iterator)
	err = iterator.Error
	return
}

// ReadEvent reads a value of the 'event' type from the given iterator.
func ReadEvent(iterator *jsoniter.Iterator) *Event {
	object := &Event{
		fieldSet_: make([]bool, 16),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == EventLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[3] = true
		case "creator":
			value := ReadUser(iterator)
			object.creator = value
			object.fieldSet_[4] = true
		case "deleted_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deletedAt = value
			object.fieldSet_[5] = true
		case "escalation":
			value := ReadEscalation(iterator)
			object.escalation = value
			object.fieldSet_[6] = true
		case "event_type":
			value := iterator.ReadString()
			object.eventType = value
			object.fieldSet_[7] = true
		case "external_reference_url":
			value := iterator.ReadString()
			object.externalReferenceUrl = value
			object.fieldSet_[8] = true
		case "follow_up":
			value := ReadFollowUp(iterator)
			object.followUp = value
			object.fieldSet_[9] = true
		case "follow_up_change":
			value := ReadFollowUpChange(iterator)
			object.followUpChange = value
			object.fieldSet_[10] = true
		case "handoff":
			value := ReadHandoff(iterator)
			object.handoff = value
			object.fieldSet_[11] = true
		case "incident":
			value := ReadIncident(iterator)
			object.incident = value
			object.fieldSet_[12] = true
		case "note":
			value := iterator.ReadString()
			object.note = value
			object.fieldSet_[13] = true
		case "status_change":
			value := ReadStatusChange(iterator)
			object.statusChange = value
			object.fieldSet_[14] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[15] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
