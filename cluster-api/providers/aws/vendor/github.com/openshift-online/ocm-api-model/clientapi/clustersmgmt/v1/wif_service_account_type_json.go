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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalWifServiceAccount writes a value of the 'wif_service_account' type to the given writer.
func MarshalWifServiceAccount(object *WifServiceAccount, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifServiceAccount(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifServiceAccount writes a value of the 'wif_service_account' type to the given stream.
func WriteWifServiceAccount(object *WifServiceAccount, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("access_method")
		stream.WriteString(string(object.accessMethod))
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.credentialRequest != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("credential_request")
		WriteWifCredentialRequest(object.credentialRequest, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("osd_role")
		stream.WriteString(object.osdRole)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.roles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("roles")
		WriteWifRoleList(object.roles, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account_id")
		stream.WriteString(object.serviceAccountId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifServiceAccount reads a value of the 'wif_service_account' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifServiceAccount(source interface{}) (object *WifServiceAccount, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifServiceAccount(iterator)
	err = iterator.Error
	return
}

// ReadWifServiceAccount reads a value of the 'wif_service_account' type from the given iterator.
func ReadWifServiceAccount(iterator *jsoniter.Iterator) *WifServiceAccount {
	object := &WifServiceAccount{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "access_method":
			text := iterator.ReadString()
			value := WifAccessMethod(text)
			object.accessMethod = value
			object.fieldSet_[0] = true
		case "credential_request":
			value := ReadWifCredentialRequest(iterator)
			object.credentialRequest = value
			object.fieldSet_[1] = true
		case "osd_role":
			value := iterator.ReadString()
			object.osdRole = value
			object.fieldSet_[2] = true
		case "roles":
			value := ReadWifRoleList(iterator)
			object.roles = value
			object.fieldSet_[3] = true
		case "service_account_id":
			value := iterator.ReadString()
			object.serviceAccountId = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
