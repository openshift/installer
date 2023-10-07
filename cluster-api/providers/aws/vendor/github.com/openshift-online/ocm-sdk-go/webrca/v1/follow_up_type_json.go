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

// MarshalFollowUp writes a value of the 'follow_up' type to the given writer.
func MarshalFollowUp(object *FollowUp, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeFollowUp(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeFollowUp writes a value of the 'follow_up' type to the given stream.
func writeFollowUp(object *FollowUp, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(FollowUpLinkKind)
	} else {
		stream.WriteString(FollowUpKind)
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
		stream.WriteObjectField("archived")
		stream.WriteBool(object.archived)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
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
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("done")
		stream.WriteBool(object.done)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("follow_up_type")
		stream.WriteString(object.followUpType)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.incident != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("incident")
		writeIncident(object.incident, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("owner")
		stream.WriteString(object.owner)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("priority")
		stream.WriteString(object.priority)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("title")
		stream.WriteString(object.title)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("url")
		stream.WriteString(object.url)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("worked_at")
		stream.WriteString((object.workedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalFollowUp reads a value of the 'follow_up' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalFollowUp(source interface{}) (object *FollowUp, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readFollowUp(iterator)
	err = iterator.Error
	return
}

// readFollowUp reads a value of the 'follow_up' type from the given iterator.
func readFollowUp(iterator *jsoniter.Iterator) *FollowUp {
	object := &FollowUp{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == FollowUpLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "archived":
			value := iterator.ReadBool()
			object.archived = value
			object.bitmap_ |= 8
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 16
		case "deleted_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deletedAt = value
			object.bitmap_ |= 32
		case "done":
			value := iterator.ReadBool()
			object.done = value
			object.bitmap_ |= 64
		case "follow_up_type":
			value := iterator.ReadString()
			object.followUpType = value
			object.bitmap_ |= 128
		case "incident":
			value := readIncident(iterator)
			object.incident = value
			object.bitmap_ |= 256
		case "owner":
			value := iterator.ReadString()
			object.owner = value
			object.bitmap_ |= 512
		case "priority":
			value := iterator.ReadString()
			object.priority = value
			object.bitmap_ |= 1024
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 2048
		case "title":
			value := iterator.ReadString()
			object.title = value
			object.bitmap_ |= 4096
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 8192
		case "url":
			value := iterator.ReadString()
			object.url = value
			object.bitmap_ |= 16384
		case "worked_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.workedAt = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
