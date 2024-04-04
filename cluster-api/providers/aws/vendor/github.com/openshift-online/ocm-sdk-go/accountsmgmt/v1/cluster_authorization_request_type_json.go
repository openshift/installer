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

// MarshalClusterAuthorizationRequest writes a value of the 'cluster_authorization_request' type to the given writer.
func MarshalClusterAuthorizationRequest(object *ClusterAuthorizationRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterAuthorizationRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterAuthorizationRequest writes a value of the 'cluster_authorization_request' type to the given stream.
func writeClusterAuthorizationRequest(object *ClusterAuthorizationRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byoc")
		stream.WriteBool(object.byoc)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_username")
		stream.WriteString(object.accountUsername)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_account_id")
		stream.WriteString(object.cloudAccountID)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider_id")
		stream.WriteString(object.cloudProviderID)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("disconnected")
		stream.WriteBool(object.disconnected)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_cluster_id")
		stream.WriteString(object.externalClusterID)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_id")
		stream.WriteString(object.productID)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_category")
		stream.WriteString(object.productCategory)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("quota_version")
		stream.WriteString(object.quotaVersion)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reserve")
		stream.WriteBool(object.reserve)
		count++
	}
	present_ = object.bitmap_&16384 != 0 && object.resources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resources")
		writeReservedResourceList(object.resources, stream)
		count++
	}
	present_ = object.bitmap_&32768 != 0
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
	object = readClusterAuthorizationRequest(iterator)
	err = iterator.Error
	return
}

// readClusterAuthorizationRequest reads a value of the 'cluster_authorization_request' type from the given iterator.
func readClusterAuthorizationRequest(iterator *jsoniter.Iterator) *ClusterAuthorizationRequest {
	object := &ClusterAuthorizationRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "byoc":
			value := iterator.ReadBool()
			object.byoc = value
			object.bitmap_ |= 1
		case "account_username":
			value := iterator.ReadString()
			object.accountUsername = value
			object.bitmap_ |= 2
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.bitmap_ |= 4
		case "cloud_account_id":
			value := iterator.ReadString()
			object.cloudAccountID = value
			object.bitmap_ |= 8
		case "cloud_provider_id":
			value := iterator.ReadString()
			object.cloudProviderID = value
			object.bitmap_ |= 16
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 32
		case "disconnected":
			value := iterator.ReadBool()
			object.disconnected = value
			object.bitmap_ |= 64
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.bitmap_ |= 128
		case "external_cluster_id":
			value := iterator.ReadString()
			object.externalClusterID = value
			object.bitmap_ |= 256
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.bitmap_ |= 512
		case "product_id":
			value := iterator.ReadString()
			object.productID = value
			object.bitmap_ |= 1024
		case "product_category":
			value := iterator.ReadString()
			object.productCategory = value
			object.bitmap_ |= 2048
		case "quota_version":
			value := iterator.ReadString()
			object.quotaVersion = value
			object.bitmap_ |= 4096
		case "reserve":
			value := iterator.ReadBool()
			object.reserve = value
			object.bitmap_ |= 8192
		case "resources":
			value := readReservedResourceList(iterator)
			object.resources = value
			object.bitmap_ |= 16384
		case "scope":
			value := iterator.ReadString()
			object.scope = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
