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

// MarshalWifCredentialRequest writes a value of the 'wif_credential_request' type to the given writer.
func MarshalWifCredentialRequest(object *WifCredentialRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifCredentialRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifCredentialRequest writes a value of the 'wif_credential_request' type to the given stream.
func WriteWifCredentialRequest(object *WifCredentialRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.secretRef != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_ref")
		WriteWifSecretRef(object.secretRef, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.serviceAccountNames != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account_names")
		WriteStringList(object.serviceAccountNames, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifCredentialRequest reads a value of the 'wif_credential_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifCredentialRequest(source interface{}) (object *WifCredentialRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifCredentialRequest(iterator)
	err = iterator.Error
	return
}

// ReadWifCredentialRequest reads a value of the 'wif_credential_request' type from the given iterator.
func ReadWifCredentialRequest(iterator *jsoniter.Iterator) *WifCredentialRequest {
	object := &WifCredentialRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "secret_ref":
			value := ReadWifSecretRef(iterator)
			object.secretRef = value
			object.bitmap_ |= 1
		case "service_account_names":
			value := ReadStringList(iterator)
			object.serviceAccountNames = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
