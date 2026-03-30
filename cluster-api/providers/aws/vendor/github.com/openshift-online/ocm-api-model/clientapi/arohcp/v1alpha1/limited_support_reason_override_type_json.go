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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalLimitedSupportReasonOverride writes a value of the 'limited_support_reason_override' type to the given writer.
func MarshalLimitedSupportReasonOverride(object *LimitedSupportReasonOverride, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLimitedSupportReasonOverride(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLimitedSupportReasonOverride writes a value of the 'limited_support_reason_override' type to the given stream.
func WriteLimitedSupportReasonOverride(object *LimitedSupportReasonOverride, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(LimitedSupportReasonOverrideLinkKind)
	} else {
		stream.WriteString(LimitedSupportReasonOverrideKind)
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
	object = ReadLimitedSupportReasonOverride(iterator)
	err = iterator.Error
	return
}

// ReadLimitedSupportReasonOverride reads a value of the 'limited_support_reason_override' type from the given iterator.
func ReadLimitedSupportReasonOverride(iterator *jsoniter.Iterator) *LimitedSupportReasonOverride {
	object := &LimitedSupportReasonOverride{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == LimitedSupportReasonOverrideLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
