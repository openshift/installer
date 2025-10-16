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

// MarshalResourceReview writes a value of the 'resource_review' type to the given writer.
func MarshalResourceReview(object *ResourceReview, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteResourceReview(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteResourceReview writes a value of the 'resource_review' type to the given stream.
func WriteResourceReview(object *ResourceReview, stream *jsoniter.Stream) {
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
	present_ = object.bitmap_&4 != 0 && object.clusterIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_ids")
		WriteStringList(object.clusterIDs, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.clusterUUIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuids")
		WriteStringList(object.clusterUUIDs, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.organizationIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_ids")
		WriteStringList(object.organizationIDs, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.subscriptionIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_ids")
		WriteStringList(object.subscriptionIDs, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalResourceReview reads a value of the 'resource_review' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalResourceReview(source interface{}) (object *ResourceReview, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadResourceReview(iterator)
	err = iterator.Error
	return
}

// ReadResourceReview reads a value of the 'resource_review' type from the given iterator.
func ReadResourceReview(iterator *jsoniter.Iterator) *ResourceReview {
	object := &ResourceReview{}
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
		case "cluster_ids":
			value := ReadStringList(iterator)
			object.clusterIDs = value
			object.bitmap_ |= 4
		case "cluster_uuids":
			value := ReadStringList(iterator)
			object.clusterUUIDs = value
			object.bitmap_ |= 8
		case "organization_ids":
			value := ReadStringList(iterator)
			object.organizationIDs = value
			object.bitmap_ |= 16
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.bitmap_ |= 32
		case "subscription_ids":
			value := ReadStringList(iterator)
			object.subscriptionIDs = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
