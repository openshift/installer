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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAccessRequest writes a value of the 'access_request' type to the given writer.
func MarshalAccessRequest(object *AccessRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccessRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccessRequest writes a value of the 'access_request' type to the given stream.
func WriteAccessRequest(object *AccessRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AccessRequestLinkKind)
	} else {
		stream.WriteString(AccessRequestKind)
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
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deadline")
		stream.WriteString(object.deadline)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deadline_at")
		stream.WriteString((object.deadlineAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.decisions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("decisions")
		WriteDecisionList(object.decisions, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("duration")
		stream.WriteString(object.duration)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_support_case_id")
		stream.WriteString(object.internalSupportCaseId)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("justification")
		stream.WriteString(object.justification)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationId)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requested_by")
		stream.WriteString(object.requestedBy)
		count++
	}
	present_ = object.bitmap_&8192 != 0 && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteAccessRequestStatus(object.status, stream)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support_case_id")
		stream.WriteString(object.supportCaseId)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessRequest reads a value of the 'access_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessRequest(source interface{}) (object *AccessRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccessRequest(iterator)
	err = iterator.Error
	return
}

// ReadAccessRequest reads a value of the 'access_request' type from the given iterator.
func ReadAccessRequest(iterator *jsoniter.Iterator) *AccessRequest {
	object := &AccessRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AccessRequestLinkKind {
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
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 16
		case "deadline":
			value := iterator.ReadString()
			object.deadline = value
			object.bitmap_ |= 32
		case "deadline_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deadlineAt = value
			object.bitmap_ |= 64
		case "decisions":
			value := ReadDecisionList(iterator)
			object.decisions = value
			object.bitmap_ |= 128
		case "duration":
			value := iterator.ReadString()
			object.duration = value
			object.bitmap_ |= 256
		case "internal_support_case_id":
			value := iterator.ReadString()
			object.internalSupportCaseId = value
			object.bitmap_ |= 512
		case "justification":
			value := iterator.ReadString()
			object.justification = value
			object.bitmap_ |= 1024
		case "organization_id":
			value := iterator.ReadString()
			object.organizationId = value
			object.bitmap_ |= 2048
		case "requested_by":
			value := iterator.ReadString()
			object.requestedBy = value
			object.bitmap_ |= 4096
		case "status":
			value := ReadAccessRequestStatus(iterator)
			object.status = value
			object.bitmap_ |= 8192
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.bitmap_ |= 16384
		case "support_case_id":
			value := iterator.ReadString()
			object.supportCaseId = value
			object.bitmap_ |= 32768
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 65536
		default:
			iterator.ReadAny()
		}
	}
	return object
}
