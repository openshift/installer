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

// MarshalNotificationDetailsRequest writes a value of the 'notification_details_request' type to the given writer.
func MarshalNotificationDetailsRequest(object *NotificationDetailsRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeNotificationDetailsRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeNotificationDetailsRequest writes a value of the 'notification_details_request' type to the given stream.
func writeNotificationDetailsRequest(object *NotificationDetailsRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bcc_address")
		stream.WriteString(object.bccAddress)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUUID)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("include_red_hat_associates")
		stream.WriteBool(object.includeRedHatAssociates)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_only")
		stream.WriteBool(object.internalOnly)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subject")
		stream.WriteString(object.subject)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNotificationDetailsRequest reads a value of the 'notification_details_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNotificationDetailsRequest(source interface{}) (object *NotificationDetailsRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readNotificationDetailsRequest(iterator)
	err = iterator.Error
	return
}

// readNotificationDetailsRequest reads a value of the 'notification_details_request' type from the given iterator.
func readNotificationDetailsRequest(iterator *jsoniter.Iterator) *NotificationDetailsRequest {
	object := &NotificationDetailsRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "bcc_address":
			value := iterator.ReadString()
			object.bccAddress = value
			object.bitmap_ |= 1
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 2
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUUID = value
			object.bitmap_ |= 4
		case "include_red_hat_associates":
			value := iterator.ReadBool()
			object.includeRedHatAssociates = value
			object.bitmap_ |= 8
		case "internal_only":
			value := iterator.ReadBool()
			object.internalOnly = value
			object.bitmap_ |= 16
		case "subject":
			value := iterator.ReadString()
			object.subject = value
			object.bitmap_ |= 32
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
