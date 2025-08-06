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

// MarshalNotificationDetailsRequest writes a value of the 'notification_details_request' type to the given writer.
func MarshalNotificationDetailsRequest(object *NotificationDetailsRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNotificationDetailsRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNotificationDetailsRequest writes a value of the 'notification_details_request' type to the given stream.
func WriteNotificationDetailsRequest(object *NotificationDetailsRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bcc_address")
		stream.WriteString(object.bccAddress)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUUID)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("include_red_hat_associates")
		stream.WriteBool(object.includeRedHatAssociates)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_only")
		stream.WriteBool(object.internalOnly)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_type")
		stream.WriteString(object.logType)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subject")
		stream.WriteString(object.subject)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
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
	object = ReadNotificationDetailsRequest(iterator)
	err = iterator.Error
	return
}

// ReadNotificationDetailsRequest reads a value of the 'notification_details_request' type from the given iterator.
func ReadNotificationDetailsRequest(iterator *jsoniter.Iterator) *NotificationDetailsRequest {
	object := &NotificationDetailsRequest{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "bcc_address":
			value := iterator.ReadString()
			object.bccAddress = value
			object.fieldSet_[0] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[1] = true
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUUID = value
			object.fieldSet_[2] = true
		case "include_red_hat_associates":
			value := iterator.ReadBool()
			object.includeRedHatAssociates = value
			object.fieldSet_[3] = true
		case "internal_only":
			value := iterator.ReadBool()
			object.internalOnly = value
			object.fieldSet_[4] = true
		case "log_type":
			value := iterator.ReadString()
			object.logType = value
			object.fieldSet_[5] = true
		case "subject":
			value := iterator.ReadString()
			object.subject = value
			object.fieldSet_[6] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
