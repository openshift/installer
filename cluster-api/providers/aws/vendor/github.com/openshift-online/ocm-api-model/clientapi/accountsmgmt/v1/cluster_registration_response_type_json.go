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

// MarshalClusterRegistrationResponse writes a value of the 'cluster_registration_response' type to the given writer.
func MarshalClusterRegistrationResponse(object *ClusterRegistrationResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterRegistrationResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterRegistrationResponse writes a value of the 'cluster_registration_response' type to the given stream.
func WriteClusterRegistrationResponse(object *ClusterRegistrationResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("authorization_token")
		stream.WriteString(object.authorizationToken)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("expires_at")
		stream.WriteString(object.expiresAt)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterRegistrationResponse reads a value of the 'cluster_registration_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterRegistrationResponse(source interface{}) (object *ClusterRegistrationResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterRegistrationResponse(iterator)
	err = iterator.Error
	return
}

// ReadClusterRegistrationResponse reads a value of the 'cluster_registration_response' type from the given iterator.
func ReadClusterRegistrationResponse(iterator *jsoniter.Iterator) *ClusterRegistrationResponse {
	object := &ClusterRegistrationResponse{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.fieldSet_[0] = true
		case "authorization_token":
			value := iterator.ReadString()
			object.authorizationToken = value
			object.fieldSet_[1] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[2] = true
		case "expires_at":
			value := iterator.ReadString()
			object.expiresAt = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
