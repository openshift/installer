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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAccessRequestPostRequest writes a value of the 'access_request_post_request' type to the given writer.
func MarshalAccessRequestPostRequest(object *AccessRequestPostRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccessRequestPostRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccessRequestPostRequest writes a value of the 'access_request_post_request' type to the given stream.
func WriteAccessRequestPostRequest(object *AccessRequestPostRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deadline")
		stream.WriteString(object.deadline)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("duration")
		stream.WriteString(object.duration)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_support_case_id")
		stream.WriteString(object.internalSupportCaseId)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("justification")
		stream.WriteString(object.justification)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support_case_id")
		stream.WriteString(object.supportCaseId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessRequestPostRequest reads a value of the 'access_request_post_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessRequestPostRequest(source interface{}) (object *AccessRequestPostRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccessRequestPostRequest(iterator)
	err = iterator.Error
	return
}

// ReadAccessRequestPostRequest reads a value of the 'access_request_post_request' type from the given iterator.
func ReadAccessRequestPostRequest(iterator *jsoniter.Iterator) *AccessRequestPostRequest {
	object := &AccessRequestPostRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.bitmap_ |= 1
		case "deadline":
			value := iterator.ReadString()
			object.deadline = value
			object.bitmap_ |= 2
		case "duration":
			value := iterator.ReadString()
			object.duration = value
			object.bitmap_ |= 4
		case "internal_support_case_id":
			value := iterator.ReadString()
			object.internalSupportCaseId = value
			object.bitmap_ |= 8
		case "justification":
			value := iterator.ReadString()
			object.justification = value
			object.bitmap_ |= 16
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.bitmap_ |= 32
		case "support_case_id":
			value := iterator.ReadString()
			object.supportCaseId = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
