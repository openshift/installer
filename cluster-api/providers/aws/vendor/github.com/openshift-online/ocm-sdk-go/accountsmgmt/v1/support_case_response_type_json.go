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
	if object.bitmap_&1 != 0 {
		stream.WriteString(SupportCaseResponseLinkKind)
	} else {
		stream.WriteString(SupportCaseResponseKind)
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
		stream.WriteObjectField("uri")
		stream.WriteString(object.uri)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("case_number")
		stream.WriteString(object.caseNumber)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUuid)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(object.severity)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionId)
		count++
	}
	present_ = object.bitmap_&2048 != 0
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
	object := &SupportCaseResponse{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == SupportCaseResponseLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "uri":
			value := iterator.ReadString()
			object.uri = value
			object.bitmap_ |= 8
		case "case_number":
			value := iterator.ReadString()
			object.caseNumber = value
			object.bitmap_ |= 16
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.bitmap_ |= 32
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUuid = value
			object.bitmap_ |= 64
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 128
		case "severity":
			value := iterator.ReadString()
			object.severity = value
			object.bitmap_ |= 256
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 512
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionId = value
			object.bitmap_ |= 1024
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
