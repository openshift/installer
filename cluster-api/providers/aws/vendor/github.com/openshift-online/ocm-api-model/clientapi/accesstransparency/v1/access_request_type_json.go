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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accesstransparency/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
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
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AccessRequestLinkKind)
	} else {
		stream.WriteString(AccessRequestKind)
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
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deadline")
		stream.WriteString(object.deadline)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deadline_at")
		stream.WriteString((object.deadlineAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.decisions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("decisions")
		WriteDecisionList(object.decisions, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("duration")
		stream.WriteString(object.duration)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_support_case_id")
		stream.WriteString(object.internalSupportCaseId)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("justification")
		stream.WriteString(object.justification)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationId)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requested_by")
		stream.WriteString(object.requestedBy)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteAccessRequestStatus(object.status, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support_case_id")
		stream.WriteString(object.supportCaseId)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
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
	object := &AccessRequest{
		fieldSet_: make([]bool, 17),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AccessRequestLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.fieldSet_[3] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[4] = true
		case "deadline":
			value := iterator.ReadString()
			object.deadline = value
			object.fieldSet_[5] = true
		case "deadline_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deadlineAt = value
			object.fieldSet_[6] = true
		case "decisions":
			value := ReadDecisionList(iterator)
			object.decisions = value
			object.fieldSet_[7] = true
		case "duration":
			value := iterator.ReadString()
			object.duration = value
			object.fieldSet_[8] = true
		case "internal_support_case_id":
			value := iterator.ReadString()
			object.internalSupportCaseId = value
			object.fieldSet_[9] = true
		case "justification":
			value := iterator.ReadString()
			object.justification = value
			object.fieldSet_[10] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationId = value
			object.fieldSet_[11] = true
		case "requested_by":
			value := iterator.ReadString()
			object.requestedBy = value
			object.fieldSet_[12] = true
		case "status":
			value := ReadAccessRequestStatus(iterator)
			object.status = value
			object.fieldSet_[13] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.fieldSet_[14] = true
		case "support_case_id":
			value := iterator.ReadString()
			object.supportCaseId = value
			object.fieldSet_[15] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[16] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
