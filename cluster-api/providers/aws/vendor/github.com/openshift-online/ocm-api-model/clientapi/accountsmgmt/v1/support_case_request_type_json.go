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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalSupportCaseRequest writes a value of the 'support_case_request' type to the given writer.
func MarshalSupportCaseRequest(object *SupportCaseRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSupportCaseRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSupportCaseRequest writes a value of the 'support_case_request' type to the given stream.
func WriteSupportCaseRequest(object *SupportCaseRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(SupportCaseRequestLinkKind)
	} else {
		stream.WriteString(SupportCaseRequestKind)
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
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUuid)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_stream_id")
		stream.WriteString(object.eventStreamId)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(object.severity)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
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
	object = ReadSupportCaseRequest(iterator)
	err = iterator.Error
	return
}

// ReadSupportCaseRequest reads a value of the 'support_case_request' type from the given iterator.
func ReadSupportCaseRequest(iterator *jsoniter.Iterator) *SupportCaseRequest {
	object := &SupportCaseRequest{
		fieldSet_: make([]bool, 10),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == SupportCaseRequestLinkKind {
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
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUuid = value
			object.fieldSet_[4] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[5] = true
		case "event_stream_id":
			value := iterator.ReadString()
			object.eventStreamId = value
			object.fieldSet_[6] = true
		case "severity":
			value := iterator.ReadString()
			object.severity = value
			object.fieldSet_[7] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.fieldSet_[8] = true
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.fieldSet_[9] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
