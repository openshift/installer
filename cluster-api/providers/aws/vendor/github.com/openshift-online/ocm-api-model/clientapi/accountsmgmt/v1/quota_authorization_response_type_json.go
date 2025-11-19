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

// MarshalQuotaAuthorizationResponse writes a value of the 'quota_authorization_response' type to the given writer.
func MarshalQuotaAuthorizationResponse(object *QuotaAuthorizationResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteQuotaAuthorizationResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteQuotaAuthorizationResponse writes a value of the 'quota_authorization_response' type to the given stream.
func WriteQuotaAuthorizationResponse(object *QuotaAuthorizationResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed")
		stream.WriteBool(object.allowed)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.excessResources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("excess_resources")
		WriteReservedResourceList(object.excessResources, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.subscription != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription")
		WriteSubscription(object.subscription, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalQuotaAuthorizationResponse reads a value of the 'quota_authorization_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalQuotaAuthorizationResponse(source interface{}) (object *QuotaAuthorizationResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadQuotaAuthorizationResponse(iterator)
	err = iterator.Error
	return
}

// ReadQuotaAuthorizationResponse reads a value of the 'quota_authorization_response' type from the given iterator.
func ReadQuotaAuthorizationResponse(iterator *jsoniter.Iterator) *QuotaAuthorizationResponse {
	object := &QuotaAuthorizationResponse{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "allowed":
			value := iterator.ReadBool()
			object.allowed = value
			object.fieldSet_[0] = true
		case "excess_resources":
			value := ReadReservedResourceList(iterator)
			object.excessResources = value
			object.fieldSet_[1] = true
		case "subscription":
			value := ReadSubscription(iterator)
			object.subscription = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
