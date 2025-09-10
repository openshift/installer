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

// MarshalWifGcp writes a value of the 'wif_gcp' type to the given writer.
func MarshalWifGcp(object *WifGcp, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifGcp(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifGcp writes a value of the 'wif_gcp' type to the given stream.
func WriteWifGcp(object *WifGcp, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("impersonator_email")
		stream.WriteString(object.impersonatorEmail)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_id")
		stream.WriteString(object.projectId)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_number")
		stream.WriteString(object.projectNumber)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_prefix")
		stream.WriteString(object.rolePrefix)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.serviceAccounts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_accounts")
		WriteWifServiceAccountList(object.serviceAccounts, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.support != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support")
		WriteWifSupport(object.support, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.workloadIdentityPool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("workload_identity_pool")
		WriteWifPool(object.workloadIdentityPool, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifGcp reads a value of the 'wif_gcp' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifGcp(source interface{}) (object *WifGcp, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifGcp(iterator)
	err = iterator.Error
	return
}

// ReadWifGcp reads a value of the 'wif_gcp' type from the given iterator.
func ReadWifGcp(iterator *jsoniter.Iterator) *WifGcp {
	object := &WifGcp{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "impersonator_email":
			value := iterator.ReadString()
			object.impersonatorEmail = value
			object.bitmap_ |= 1
		case "project_id":
			value := iterator.ReadString()
			object.projectId = value
			object.bitmap_ |= 2
		case "project_number":
			value := iterator.ReadString()
			object.projectNumber = value
			object.bitmap_ |= 4
		case "role_prefix":
			value := iterator.ReadString()
			object.rolePrefix = value
			object.bitmap_ |= 8
		case "service_accounts":
			value := ReadWifServiceAccountList(iterator)
			object.serviceAccounts = value
			object.bitmap_ |= 16
		case "support":
			value := ReadWifSupport(iterator)
			object.support = value
			object.bitmap_ |= 32
		case "workload_identity_pool":
			value := ReadWifPool(iterator)
			object.workloadIdentityPool = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
