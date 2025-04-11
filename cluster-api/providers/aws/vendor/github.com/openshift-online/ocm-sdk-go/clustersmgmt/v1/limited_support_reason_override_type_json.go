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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalLimitedSupportReasonOverride writes a value of the 'limited_support_reason_override' type to the given writer.
func MarshalLimitedSupportReasonOverride(object *LimitedSupportReasonOverride, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeLimitedSupportReasonOverride(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeLimitedSupportReasonOverride writes a value of the 'limited_support_reason_override' type to the given stream.
func writeLimitedSupportReasonOverride(object *LimitedSupportReasonOverride, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(LimitedSupportReasonOverrideLinkKind)
	} else {
		stream.WriteString(LimitedSupportReasonOverrideKind)
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
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLimitedSupportReasonOverride reads a value of the 'limited_support_reason_override' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLimitedSupportReasonOverride(source interface{}) (object *LimitedSupportReasonOverride, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readLimitedSupportReasonOverride(iterator)
	err = iterator.Error
	return
}

// readLimitedSupportReasonOverride reads a value of the 'limited_support_reason_override' type from the given iterator.
func readLimitedSupportReasonOverride(iterator *jsoniter.Iterator) *LimitedSupportReasonOverride {
	object := &LimitedSupportReasonOverride{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == LimitedSupportReasonOverrideLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
