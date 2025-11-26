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

// MarshalExternalAuthClientConfig writes a value of the 'external_auth_client_config' type to the given writer.
func MarshalExternalAuthClientConfig(object *ExternalAuthClientConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteExternalAuthClientConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteExternalAuthClientConfig writes a value of the 'external_auth_client_config' type to the given stream.
func WriteExternalAuthClientConfig(object *ExternalAuthClientConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.component != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("component")
		WriteClientComponent(object.component, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.extraScopes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("extra_scopes")
		WriteStringList(object.extraScopes, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret")
		stream.WriteString(object.secret)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(string(object.type_))
	}
	stream.WriteObjectEnd()
}

// UnmarshalExternalAuthClientConfig reads a value of the 'external_auth_client_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalExternalAuthClientConfig(source interface{}) (object *ExternalAuthClientConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadExternalAuthClientConfig(iterator)
	err = iterator.Error
	return
}

// ReadExternalAuthClientConfig reads a value of the 'external_auth_client_config' type from the given iterator.
func ReadExternalAuthClientConfig(iterator *jsoniter.Iterator) *ExternalAuthClientConfig {
	object := &ExternalAuthClientConfig{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[0] = true
		case "component":
			value := ReadClientComponent(iterator)
			object.component = value
			object.fieldSet_[1] = true
		case "extra_scopes":
			value := ReadStringList(iterator)
			object.extraScopes = value
			object.fieldSet_[2] = true
		case "secret":
			value := iterator.ReadString()
			object.secret = value
			object.fieldSet_[3] = true
		case "type":
			text := iterator.ReadString()
			value := ExternalAuthClientType(text)
			object.type_ = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
