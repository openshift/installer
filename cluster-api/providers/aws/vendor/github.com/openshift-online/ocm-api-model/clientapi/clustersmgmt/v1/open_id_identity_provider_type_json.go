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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalOpenIDIdentityProvider writes a value of the 'open_ID_identity_provider' type to the given writer.
func MarshalOpenIDIdentityProvider(object *OpenIDIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteOpenIDIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOpenIDIdentityProvider writes a value of the 'open_ID_identity_provider' type to the given stream.
func WriteOpenIDIdentityProvider(object *OpenIDIdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ca")
		stream.WriteString(object.ca)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.claims != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("claims")
		WriteOpenIDClaims(object.claims, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_secret")
		stream.WriteString(object.clientSecret)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.extraAuthorizeParameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("extra_authorize_parameters")
		if object.extraAuthorizeParameters != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.extraAuthorizeParameters))
			i := 0
			for key := range object.extraAuthorizeParameters {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.extraAuthorizeParameters[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.extraScopes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("extra_scopes")
		WriteStringList(object.extraScopes, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("issuer")
		stream.WriteString(object.issuer)
	}
	stream.WriteObjectEnd()
}

// UnmarshalOpenIDIdentityProvider reads a value of the 'open_ID_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalOpenIDIdentityProvider(source interface{}) (object *OpenIDIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadOpenIDIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadOpenIDIdentityProvider reads a value of the 'open_ID_identity_provider' type from the given iterator.
func ReadOpenIDIdentityProvider(iterator *jsoniter.Iterator) *OpenIDIdentityProvider {
	object := &OpenIDIdentityProvider{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "ca":
			value := iterator.ReadString()
			object.ca = value
			object.fieldSet_[0] = true
		case "claims":
			value := ReadOpenIDClaims(iterator)
			object.claims = value
			object.fieldSet_[1] = true
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.fieldSet_[2] = true
		case "client_secret":
			value := iterator.ReadString()
			object.clientSecret = value
			object.fieldSet_[3] = true
		case "extra_authorize_parameters":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.extraAuthorizeParameters = value
			object.fieldSet_[4] = true
		case "extra_scopes":
			value := ReadStringList(iterator)
			object.extraScopes = value
			object.fieldSet_[5] = true
		case "issuer":
			value := iterator.ReadString()
			object.issuer = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
