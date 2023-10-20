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

// MarshalSubscriptionNotify writes a value of the 'subscription_notify' type to the given writer.
func MarshalSubscriptionNotify(object *SubscriptionNotify, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSubscriptionNotify(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSubscriptionNotify writes a value of the 'subscription_notify' type to the given stream.
func writeSubscriptionNotify(object *SubscriptionNotify, stream *jsoniter.Stream) {
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
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("template_name")
		stream.WriteString(object.templateName)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.templateParameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("template_parameters")
		writeTemplateParameterList(object.templateParameters, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSubscriptionNotify reads a value of the 'subscription_notify' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSubscriptionNotify(source interface{}) (object *SubscriptionNotify, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSubscriptionNotify(iterator)
	err = iterator.Error
	return
}

// readSubscriptionNotify reads a value of the 'subscription_notify' type from the given iterator.
func readSubscriptionNotify(iterator *jsoniter.Iterator) *SubscriptionNotify {
	object := &SubscriptionNotify{}
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
		case "template_name":
			value := iterator.ReadString()
			object.templateName = value
			object.bitmap_ |= 128
		case "template_parameters":
			value := readTemplateParameterList(iterator)
			object.templateParameters = value
			object.bitmap_ |= 256
		default:
			iterator.ReadAny()
		}
	}
	return object
}
