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

// MarshalClusterRegistrationRequest writes a value of the 'cluster_registration_request' type to the given writer.
func MarshalClusterRegistrationRequest(object *ClusterRegistrationRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterRegistrationRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterRegistrationRequest writes a value of the 'cluster_registration_request' type to the given stream.
func WriteClusterRegistrationRequest(object *ClusterRegistrationRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("authorization_token")
		stream.WriteString(object.authorizationToken)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterRegistrationRequest reads a value of the 'cluster_registration_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterRegistrationRequest(source interface{}) (object *ClusterRegistrationRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterRegistrationRequest(iterator)
	err = iterator.Error
	return
}

// ReadClusterRegistrationRequest reads a value of the 'cluster_registration_request' type from the given iterator.
func ReadClusterRegistrationRequest(iterator *jsoniter.Iterator) *ClusterRegistrationRequest {
	object := &ClusterRegistrationRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "authorization_token":
			value := iterator.ReadString()
			object.authorizationToken = value
			object.bitmap_ |= 1
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
