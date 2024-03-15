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

// MarshalExternalAuthClientConfig writes a value of the 'external_auth_client_config' type to the given writer.
func MarshalExternalAuthClientConfig(object *ExternalAuthClientConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeExternalAuthClientConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeExternalAuthClientConfig writes a value of the 'external_auth_client_config' type to the given stream.
func writeExternalAuthClientConfig(object *ExternalAuthClientConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.component != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("component")
		writeClientComponent(object.component, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.extraScopes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("extra_scopes")
		writeStringList(object.extraScopes, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret")
		stream.WriteString(object.secret)
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
	object = readExternalAuthClientConfig(iterator)
	err = iterator.Error
	return
}

// readExternalAuthClientConfig reads a value of the 'external_auth_client_config' type from the given iterator.
func readExternalAuthClientConfig(iterator *jsoniter.Iterator) *ExternalAuthClientConfig {
	object := &ExternalAuthClientConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 1
		case "component":
			value := readClientComponent(iterator)
			object.component = value
			object.bitmap_ |= 2
		case "extra_scopes":
			value := readStringList(iterator)
			object.extraScopes = value
			object.bitmap_ |= 4
		case "secret":
			value := iterator.ReadString()
			object.secret = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
