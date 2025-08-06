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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalInflightCheck writes a value of the 'inflight_check' type to the given writer.
func MarshalInflightCheck(object *InflightCheck, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteInflightCheck(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteInflightCheck writes a value of the 'inflight_check' type to the given stream.
func WriteInflightCheck(object *InflightCheck, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(InflightCheckLinkKind)
	} else {
		stream.WriteString(InflightCheckKind)
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
		stream.WriteObjectField("details")
		stream.WriteVal(object.details)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ended_at")
		stream.WriteString((object.endedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("restarts")
		stream.WriteInt(object.restarts)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("started_at")
		stream.WriteString((object.startedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
	}
	stream.WriteObjectEnd()
}

// UnmarshalInflightCheck reads a value of the 'inflight_check' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalInflightCheck(source interface{}) (object *InflightCheck, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadInflightCheck(iterator)
	err = iterator.Error
	return
}

// ReadInflightCheck reads a value of the 'inflight_check' type from the given iterator.
func ReadInflightCheck(iterator *jsoniter.Iterator) *InflightCheck {
	object := &InflightCheck{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == InflightCheckLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "details":
			var value interface{}
			iterator.ReadVal(&value)
			object.details = value
			object.bitmap_ |= 8
		case "ended_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endedAt = value
			object.bitmap_ |= 16
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 32
		case "restarts":
			value := iterator.ReadInt()
			object.restarts = value
			object.bitmap_ |= 64
		case "started_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.startedAt = value
			object.bitmap_ |= 128
		case "state":
			text := iterator.ReadString()
			value := InflightCheckState(text)
			object.state = value
			object.bitmap_ |= 256
		default:
			iterator.ReadAny()
		}
	}
	return object
}
