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

// MarshalQuotaAuthorizationRequest writes a value of the 'quota_authorization_request' type to the given writer.
func MarshalQuotaAuthorizationRequest(object *QuotaAuthorizationRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteQuotaAuthorizationRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteQuotaAuthorizationRequest writes a value of the 'quota_authorization_request' type to the given stream.
func WriteQuotaAuthorizationRequest(object *QuotaAuthorizationRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_username")
		stream.WriteString(object.accountUsername)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_id")
		stream.WriteString(object.productID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_category")
		stream.WriteString(object.productCategory)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("quota_version")
		stream.WriteString(object.quotaVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reserve")
		stream.WriteBool(object.reserve)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.resources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resources")
		WriteReservedResourceList(object.resources, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalQuotaAuthorizationRequest reads a value of the 'quota_authorization_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalQuotaAuthorizationRequest(source interface{}) (object *QuotaAuthorizationRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadQuotaAuthorizationRequest(iterator)
	err = iterator.Error
	return
}

// ReadQuotaAuthorizationRequest reads a value of the 'quota_authorization_request' type from the given iterator.
func ReadQuotaAuthorizationRequest(iterator *jsoniter.Iterator) *QuotaAuthorizationRequest {
	object := &QuotaAuthorizationRequest{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_username":
			value := iterator.ReadString()
			object.accountUsername = value
			object.fieldSet_[0] = true
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.fieldSet_[1] = true
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.fieldSet_[2] = true
		case "product_id":
			value := iterator.ReadString()
			object.productID = value
			object.fieldSet_[3] = true
		case "product_category":
			value := iterator.ReadString()
			object.productCategory = value
			object.fieldSet_[4] = true
		case "quota_version":
			value := iterator.ReadString()
			object.quotaVersion = value
			object.fieldSet_[5] = true
		case "reserve":
			value := iterator.ReadBool()
			object.reserve = value
			object.fieldSet_[6] = true
		case "resources":
			value := ReadReservedResourceList(iterator)
			object.resources = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
