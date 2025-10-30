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

// MarshalResourceReviewRequest writes a value of the 'resource_review_request' type to the given writer.
func MarshalResourceReviewRequest(object *ResourceReviewRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteResourceReviewRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteResourceReviewRequest writes a value of the 'resource_review_request' type to the given stream.
func WriteResourceReviewRequest(object *ResourceReviewRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_username")
		stream.WriteString(object.accountUsername)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("action")
		stream.WriteString(object.action)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.excludeSubscriptionStatuses != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("exclude_subscription_statuses")
		WriteSubscriptionStatusList(object.excludeSubscriptionStatuses, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reduce_cluster_list")
		stream.WriteBool(object.reduceClusterList)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalResourceReviewRequest reads a value of the 'resource_review_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalResourceReviewRequest(source interface{}) (object *ResourceReviewRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadResourceReviewRequest(iterator)
	err = iterator.Error
	return
}

// ReadResourceReviewRequest reads a value of the 'resource_review_request' type from the given iterator.
func ReadResourceReviewRequest(iterator *jsoniter.Iterator) *ResourceReviewRequest {
	object := &ResourceReviewRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_username":
			value := iterator.ReadString()
			object.accountUsername = value
			object.bitmap_ |= 1
		case "action":
			value := iterator.ReadString()
			object.action = value
			object.bitmap_ |= 2
		case "exclude_subscription_statuses":
			value := ReadSubscriptionStatusList(iterator)
			object.excludeSubscriptionStatuses = value
			object.bitmap_ |= 4
		case "reduce_cluster_list":
			value := iterator.ReadBool()
			object.reduceClusterList = value
			object.bitmap_ |= 8
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
