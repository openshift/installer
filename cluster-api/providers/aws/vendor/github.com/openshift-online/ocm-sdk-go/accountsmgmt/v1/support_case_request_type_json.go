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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSupportCaseRequest writes a value of the 'support_case_request' type to the given writer.
func MarshalSupportCaseRequest(object *SupportCaseRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSupportCaseRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSupportCaseRequest writes a value of the 'support_case_request' type to the given stream.
func writeSupportCaseRequest(object *SupportCaseRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(SupportCaseRequestLinkKind)
	} else {
		stream.WriteString(SupportCaseRequestKind)
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
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUuid)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_stream_id")
		stream.WriteString(object.eventStreamId)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(object.severity)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("summary")
		stream.WriteString(object.summary)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSupportCaseRequest reads a value of the 'support_case_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSupportCaseRequest(source interface{}) (object *SupportCaseRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSupportCaseRequest(iterator)
	err = iterator.Error
	return
}

// readSupportCaseRequest reads a value of the 'support_case_request' type from the given iterator.
func readSupportCaseRequest(iterator *jsoniter.Iterator) *SupportCaseRequest {
	object := &SupportCaseRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == SupportCaseRequestLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.bitmap_ |= 8
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUuid = value
			object.bitmap_ |= 16
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 32
		case "event_stream_id":
			value := iterator.ReadString()
			object.eventStreamId = value
			object.bitmap_ |= 64
		case "severity":
			value := iterator.ReadString()
			object.severity = value
			object.bitmap_ |= 128
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.bitmap_ |= 256
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.bitmap_ |= 512
		default:
			iterator.ReadAny()
		}
	}
	return object
}
