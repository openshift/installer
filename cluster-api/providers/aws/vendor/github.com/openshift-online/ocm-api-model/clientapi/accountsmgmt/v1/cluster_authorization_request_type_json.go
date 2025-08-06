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

// MarshalClusterAuthorizationRequest writes a value of the 'cluster_authorization_request' type to the given writer.
func MarshalClusterAuthorizationRequest(object *ClusterAuthorizationRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterAuthorizationRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterAuthorizationRequest writes a value of the 'cluster_authorization_request' type to the given stream.
func WriteClusterAuthorizationRequest(object *ClusterAuthorizationRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byoc")
		stream.WriteBool(object.byoc)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_username")
		stream.WriteString(object.accountUsername)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_account_id")
		stream.WriteString(object.cloudAccountID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider_id")
		stream.WriteString(object.cloudProviderID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("disconnected")
		stream.WriteBool(object.disconnected)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_cluster_id")
		stream.WriteString(object.externalClusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_id")
		stream.WriteString(object.productID)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_category")
		stream.WriteString(object.productCategory)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("quota_version")
		stream.WriteString(object.quotaVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reserve")
		stream.WriteBool(object.reserve)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14] && object.resources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resources")
		WriteReservedResourceList(object.resources, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rh_region_id")
		stream.WriteString(object.rhRegionID)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("scope")
		stream.WriteString(object.scope)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterAuthorizationRequest reads a value of the 'cluster_authorization_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterAuthorizationRequest(source interface{}) (object *ClusterAuthorizationRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterAuthorizationRequest(iterator)
	err = iterator.Error
	return
}

// ReadClusterAuthorizationRequest reads a value of the 'cluster_authorization_request' type from the given iterator.
func ReadClusterAuthorizationRequest(iterator *jsoniter.Iterator) *ClusterAuthorizationRequest {
	object := &ClusterAuthorizationRequest{
		fieldSet_: make([]bool, 17),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "byoc":
			value := iterator.ReadBool()
			object.byoc = value
			object.fieldSet_[0] = true
		case "account_username":
			value := iterator.ReadString()
			object.accountUsername = value
			object.fieldSet_[1] = true
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.fieldSet_[2] = true
		case "cloud_account_id":
			value := iterator.ReadString()
			object.cloudAccountID = value
			object.fieldSet_[3] = true
		case "cloud_provider_id":
			value := iterator.ReadString()
			object.cloudProviderID = value
			object.fieldSet_[4] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[5] = true
		case "disconnected":
			value := iterator.ReadBool()
			object.disconnected = value
			object.fieldSet_[6] = true
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.fieldSet_[7] = true
		case "external_cluster_id":
			value := iterator.ReadString()
			object.externalClusterID = value
			object.fieldSet_[8] = true
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.fieldSet_[9] = true
		case "product_id":
			value := iterator.ReadString()
			object.productID = value
			object.fieldSet_[10] = true
		case "product_category":
			value := iterator.ReadString()
			object.productCategory = value
			object.fieldSet_[11] = true
		case "quota_version":
			value := iterator.ReadString()
			object.quotaVersion = value
			object.fieldSet_[12] = true
		case "reserve":
			value := iterator.ReadBool()
			object.reserve = value
			object.fieldSet_[13] = true
		case "resources":
			value := ReadReservedResourceList(iterator)
			object.resources = value
			object.fieldSet_[14] = true
		case "rh_region_id":
			value := iterator.ReadString()
			object.rhRegionID = value
			object.fieldSet_[15] = true
		case "scope":
			value := iterator.ReadString()
			object.scope = value
			object.fieldSet_[16] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
