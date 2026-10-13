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

// MarshalWifIdentityProvider writes a value of the 'wif_identity_provider' type to the given writer.
func MarshalWifIdentityProvider(object *WifIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifIdentityProvider writes a value of the 'wif_identity_provider' type to the given stream.
func WriteWifIdentityProvider(object *WifIdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.allowedAudiences != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed_audiences")
		WriteStringList(object.allowedAudiences, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_provider_id")
		stream.WriteString(object.identityProviderId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("issuer_url")
		stream.WriteString(object.issuerUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("jwks")
		stream.WriteString(object.jwks)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifIdentityProvider reads a value of the 'wif_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifIdentityProvider(source interface{}) (object *WifIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadWifIdentityProvider reads a value of the 'wif_identity_provider' type from the given iterator.
func ReadWifIdentityProvider(iterator *jsoniter.Iterator) *WifIdentityProvider {
	object := &WifIdentityProvider{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "allowed_audiences":
			value := ReadStringList(iterator)
			object.allowedAudiences = value
			object.fieldSet_[0] = true
		case "identity_provider_id":
			value := iterator.ReadString()
			object.identityProviderId = value
			object.fieldSet_[1] = true
		case "issuer_url":
			value := iterator.ReadString()
			object.issuerUrl = value
			object.fieldSet_[2] = true
		case "jwks":
			value := iterator.ReadString()
			object.jwks = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
