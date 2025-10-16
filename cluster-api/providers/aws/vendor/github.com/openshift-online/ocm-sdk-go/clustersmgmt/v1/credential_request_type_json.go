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

// MarshalCredentialRequest writes a value of the 'credential_request' type to the given writer.
func MarshalCredentialRequest(object *CredentialRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCredentialRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCredentialRequest writes a value of the 'credential_request' type to the given stream.
func WriteCredentialRequest(object *CredentialRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.policyPermissions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("policy_permissions")
		WriteStringList(object.policyPermissions, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account")
		stream.WriteString(object.serviceAccount)
	}
	stream.WriteObjectEnd()
}

// UnmarshalCredentialRequest reads a value of the 'credential_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCredentialRequest(source interface{}) (object *CredentialRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadCredentialRequest(iterator)
	err = iterator.Error
	return
}

// ReadCredentialRequest reads a value of the 'credential_request' type from the given iterator.
func ReadCredentialRequest(iterator *jsoniter.Iterator) *CredentialRequest {
	object := &CredentialRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 1
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.bitmap_ |= 2
		case "policy_permissions":
			value := ReadStringList(iterator)
			object.policyPermissions = value
			object.bitmap_ |= 4
		case "service_account":
			value := iterator.ReadString()
			object.serviceAccount = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
