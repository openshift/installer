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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSelfAccessReviewResponse writes a value of the 'self_access_review_response' type to the given writer.
func MarshalSelfAccessReviewResponse(object *SelfAccessReviewResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSelfAccessReviewResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSelfAccessReviewResponse writes a value of the 'self_access_review_response' type to the given stream.
func writeSelfAccessReviewResponse(object *SelfAccessReviewResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("action")
		stream.WriteString(object.action)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed")
		stream.WriteBool(object.allowed)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUUID)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("is_ocm_internal")
		stream.WriteBool(object.isOCMInternal)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reason")
		stream.WriteString(object.reason)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSelfAccessReviewResponse reads a value of the 'self_access_review_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSelfAccessReviewResponse(source interface{}) (object *SelfAccessReviewResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSelfAccessReviewResponse(iterator)
	err = iterator.Error
	return
}

// readSelfAccessReviewResponse reads a value of the 'self_access_review_response' type from the given iterator.
func readSelfAccessReviewResponse(iterator *jsoniter.Iterator) *SelfAccessReviewResponse {
	object := &SelfAccessReviewResponse{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "action":
			value := iterator.ReadString()
			object.action = value
			object.bitmap_ |= 1
		case "allowed":
			value := iterator.ReadBool()
			object.allowed = value
			object.bitmap_ |= 2
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 4
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUUID = value
			object.bitmap_ |= 8
		case "is_ocm_internal":
			value := iterator.ReadBool()
			object.isOCMInternal = value
			object.bitmap_ |= 16
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.bitmap_ |= 32
		case "reason":
			value := iterator.ReadString()
			object.reason = value
			object.bitmap_ |= 64
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.bitmap_ |= 128
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.bitmap_ |= 256
		default:
			iterator.ReadAny()
		}
	}
	return object
}
