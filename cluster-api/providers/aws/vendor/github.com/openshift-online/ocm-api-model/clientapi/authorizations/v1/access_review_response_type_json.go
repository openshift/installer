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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAccessReviewResponse writes a value of the 'access_review_response' type to the given writer.
func MarshalAccessReviewResponse(object *AccessReviewResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccessReviewResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccessReviewResponse writes a value of the 'access_review_response' type to the given stream.
func WriteAccessReviewResponse(object *AccessReviewResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_username")
		stream.WriteString(object.accountUsername)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("action")
		stream.WriteString(object.action)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed")
		stream.WriteBool(object.allowed)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUUID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("is_ocm_internal")
		stream.WriteBool(object.isOCMInternal)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reason")
		stream.WriteString(object.reason)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessReviewResponse reads a value of the 'access_review_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessReviewResponse(source interface{}) (object *AccessReviewResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccessReviewResponse(iterator)
	err = iterator.Error
	return
}

// ReadAccessReviewResponse reads a value of the 'access_review_response' type from the given iterator.
func ReadAccessReviewResponse(iterator *jsoniter.Iterator) *AccessReviewResponse {
	object := &AccessReviewResponse{
		fieldSet_: make([]bool, 10),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_username":
			value := iterator.ReadString()
			object.accountUsername = value
			object.fieldSet_[0] = true
		case "action":
			value := iterator.ReadString()
			object.action = value
			object.fieldSet_[1] = true
		case "allowed":
			value := iterator.ReadBool()
			object.allowed = value
			object.fieldSet_[2] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[3] = true
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUUID = value
			object.fieldSet_[4] = true
		case "is_ocm_internal":
			value := iterator.ReadBool()
			object.isOCMInternal = value
			object.fieldSet_[5] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.fieldSet_[6] = true
		case "reason":
			value := iterator.ReadString()
			object.reason = value
			object.fieldSet_[7] = true
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.fieldSet_[8] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.fieldSet_[9] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
