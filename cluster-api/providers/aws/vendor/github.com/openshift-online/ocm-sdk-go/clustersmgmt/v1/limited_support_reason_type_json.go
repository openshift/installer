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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalLimitedSupportReason writes a value of the 'limited_support_reason' type to the given writer.
func MarshalLimitedSupportReason(object *LimitedSupportReason, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeLimitedSupportReason(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeLimitedSupportReason writes a value of the 'limited_support_reason' type to the given stream.
func writeLimitedSupportReason(object *LimitedSupportReason, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(LimitedSupportReasonLinkKind)
	} else {
		stream.WriteString(LimitedSupportReasonKind)
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
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("details")
		stream.WriteString(object.details)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("detection_type")
		stream.WriteString(string(object.detectionType))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("summary")
		stream.WriteString(object.summary)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.template != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("template")
		writeLimitedSupportReasonTemplate(object.template, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLimitedSupportReason reads a value of the 'limited_support_reason' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLimitedSupportReason(source interface{}) (object *LimitedSupportReason, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readLimitedSupportReason(iterator)
	err = iterator.Error
	return
}

// readLimitedSupportReason reads a value of the 'limited_support_reason' type from the given iterator.
func readLimitedSupportReason(iterator *jsoniter.Iterator) *LimitedSupportReason {
	object := &LimitedSupportReason{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == LimitedSupportReasonLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 8
		case "details":
			value := iterator.ReadString()
			object.details = value
			object.bitmap_ |= 16
		case "detection_type":
			text := iterator.ReadString()
			value := DetectionType(text)
			object.detectionType = value
			object.bitmap_ |= 32
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.bitmap_ |= 64
		case "template":
			value := readLimitedSupportReasonTemplate(iterator)
			object.template = value
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
