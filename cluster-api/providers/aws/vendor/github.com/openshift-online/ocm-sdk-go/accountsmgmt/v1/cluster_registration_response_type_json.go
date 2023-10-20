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

// MarshalClusterRegistrationResponse writes a value of the 'cluster_registration_response' type to the given writer.
func MarshalClusterRegistrationResponse(object *ClusterRegistrationResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterRegistrationResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterRegistrationResponse writes a value of the 'cluster_registration_response' type to the given stream.
func writeClusterRegistrationResponse(object *ClusterRegistrationResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("authorization_token")
		stream.WriteString(object.authorizationToken)
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
	object = readClusterRegistrationResponse(iterator)
	err = iterator.Error
	return
}

// readClusterRegistrationResponse reads a value of the 'cluster_registration_response' type from the given iterator.
func readClusterRegistrationResponse(iterator *jsoniter.Iterator) *ClusterRegistrationResponse {
	object := &ClusterRegistrationResponse{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.bitmap_ |= 1
		case "authorization_token":
			value := iterator.ReadString()
			object.authorizationToken = value
			object.bitmap_ |= 2
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 4
		case "expires_at":
			value := iterator.ReadString()
			object.expiresAt = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
