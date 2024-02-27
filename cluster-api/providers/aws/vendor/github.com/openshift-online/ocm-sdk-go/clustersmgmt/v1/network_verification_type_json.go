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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalNetworkVerification writes a value of the 'network_verification' type to the given writer.
func MarshalNetworkVerification(object *NetworkVerification, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeNetworkVerification(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeNetworkVerification writes a value of the 'network_verification' type to the given stream.
func writeNetworkVerification(object *NetworkVerification, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.cloudProviderData != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider_data")
		writeCloudProviderData(object.cloudProviderData, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.items != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("items")
		writeSubnetNetworkVerificationList(object.items, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("platform")
		stream.WriteString(string(object.platform))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("total")
		stream.WriteInt(object.total)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNetworkVerification reads a value of the 'network_verification' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNetworkVerification(source interface{}) (object *NetworkVerification, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readNetworkVerification(iterator)
	err = iterator.Error
	return
}

// readNetworkVerification reads a value of the 'network_verification' type from the given iterator.
func readNetworkVerification(iterator *jsoniter.Iterator) *NetworkVerification {
	object := &NetworkVerification{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cloud_provider_data":
			value := readCloudProviderData(iterator)
			object.cloudProviderData = value
			object.bitmap_ |= 1
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.bitmap_ |= 2
		case "items":
			value := readSubnetNetworkVerificationList(iterator)
			object.items = value
			object.bitmap_ |= 4
		case "platform":
			text := iterator.ReadString()
			value := Platform(text)
			object.platform = value
			object.bitmap_ |= 8
		case "total":
			value := iterator.ReadInt()
			object.total = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
