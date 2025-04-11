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

// MarshalWifServiceAccount writes a value of the 'wif_service_account' type to the given writer.
func MarshalWifServiceAccount(object *WifServiceAccount, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeWifServiceAccount(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeWifServiceAccount writes a value of the 'wif_service_account' type to the given stream.
func writeWifServiceAccount(object *WifServiceAccount, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("access_method")
		stream.WriteString(string(object.accessMethod))
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.credentialRequest != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("credential_request")
		writeWifCredentialRequest(object.credentialRequest, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("osd_role")
		stream.WriteString(object.osdRole)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.roles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("roles")
		writeWifRoleList(object.roles, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
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
	object = readWifServiceAccount(iterator)
	err = iterator.Error
	return
}

// readWifServiceAccount reads a value of the 'wif_service_account' type from the given iterator.
func readWifServiceAccount(iterator *jsoniter.Iterator) *WifServiceAccount {
	object := &WifServiceAccount{}
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
			object.bitmap_ |= 1
		case "credential_request":
			value := readWifCredentialRequest(iterator)
			object.credentialRequest = value
			object.bitmap_ |= 2
		case "osd_role":
			value := iterator.ReadString()
			object.osdRole = value
			object.bitmap_ |= 4
		case "roles":
			value := readWifRoleList(iterator)
			object.roles = value
			object.bitmap_ |= 8
		case "service_account_id":
			value := iterator.ReadString()
			object.serviceAccountId = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
