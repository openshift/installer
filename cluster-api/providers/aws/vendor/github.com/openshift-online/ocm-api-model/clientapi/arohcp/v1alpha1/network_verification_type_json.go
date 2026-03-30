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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalNetworkVerification writes a value of the 'network_verification' type to the given writer.
func MarshalNetworkVerification(object *NetworkVerification, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNetworkVerification(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNetworkVerification writes a value of the 'network_verification' type to the given stream.
func WriteNetworkVerification(object *NetworkVerification, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.cloudProviderData != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider_data")
		WriteCloudProviderData(object.cloudProviderData, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.items != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("items")
		WriteSubnetNetworkVerificationList(object.items, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("platform")
		stream.WriteString(string(object.platform))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
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
	object = ReadNetworkVerification(iterator)
	err = iterator.Error
	return
}

// ReadNetworkVerification reads a value of the 'network_verification' type from the given iterator.
func ReadNetworkVerification(iterator *jsoniter.Iterator) *NetworkVerification {
	object := &NetworkVerification{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cloud_provider_data":
			value := ReadCloudProviderData(iterator)
			object.cloudProviderData = value
			object.fieldSet_[0] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.fieldSet_[1] = true
		case "items":
			value := ReadSubnetNetworkVerificationList(iterator)
			object.items = value
			object.fieldSet_[2] = true
		case "platform":
			text := iterator.ReadString()
			value := Platform(text)
			object.platform = value
			object.fieldSet_[3] = true
		case "total":
			value := iterator.ReadInt()
			object.total = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
