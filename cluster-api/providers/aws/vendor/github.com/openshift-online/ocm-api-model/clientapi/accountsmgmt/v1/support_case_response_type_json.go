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

// MarshalSupportCaseResponse writes a value of the 'support_case_response' type to the given writer.
func MarshalSupportCaseResponse(object *SupportCaseResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSupportCaseResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSupportCaseResponse writes a value of the 'support_case_response' type to the given stream.
func WriteSupportCaseResponse(object *SupportCaseResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(SupportCaseResponseLinkKind)
	} else {
		stream.WriteString(SupportCaseResponseKind)
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
		stream.WriteObjectField("uri")
		stream.WriteString(object.uri)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("case_number")
		stream.WriteString(object.caseNumber)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUuid)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(object.severity)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("summary")
		stream.WriteString(object.summary)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSupportCaseResponse reads a value of the 'support_case_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSupportCaseResponse(source interface{}) (object *SupportCaseResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSupportCaseResponse(iterator)
	err = iterator.Error
	return
}

// ReadSupportCaseResponse reads a value of the 'support_case_response' type from the given iterator.
func ReadSupportCaseResponse(iterator *jsoniter.Iterator) *SupportCaseResponse {
	object := &SupportCaseResponse{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == SupportCaseResponseLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "uri":
			value := iterator.ReadString()
			object.uri = value
			object.fieldSet_[3] = true
		case "case_number":
			value := iterator.ReadString()
			object.caseNumber = value
			object.fieldSet_[4] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.fieldSet_[5] = true
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUuid = value
			object.fieldSet_[6] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[7] = true
		case "severity":
			value := iterator.ReadString()
			object.severity = value
			object.fieldSet_[8] = true
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.fieldSet_[9] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.fieldSet_[10] = true
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
