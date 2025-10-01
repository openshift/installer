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

// MarshalExternalAuth writes a value of the 'external_auth' type to the given writer.
func MarshalExternalAuth(object *ExternalAuth, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteExternalAuth(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteExternalAuth writes a value of the 'external_auth' type to the given stream.
func WriteExternalAuth(object *ExternalAuth, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ExternalAuthLinkKind)
	} else {
		stream.WriteString(ExternalAuthKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.claim != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("claim")
		WriteExternalAuthClaim(object.claim, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.clients != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("clients")
		WriteExternalAuthClientConfigList(object.clients, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.issuer != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("issuer")
		WriteTokenIssuer(object.issuer, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalExternalAuth reads a value of the 'external_auth' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalExternalAuth(source interface{}) (object *ExternalAuth, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadExternalAuth(iterator)
	err = iterator.Error
	return
}

// ReadExternalAuth reads a value of the 'external_auth' type from the given iterator.
func ReadExternalAuth(iterator *jsoniter.Iterator) *ExternalAuth {
	object := &ExternalAuth{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ExternalAuthLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "claim":
			value := ReadExternalAuthClaim(iterator)
			object.claim = value
			object.bitmap_ |= 8
		case "clients":
			value := ReadExternalAuthClientConfigList(iterator)
			object.clients = value
			object.bitmap_ |= 16
		case "issuer":
			value := ReadTokenIssuer(iterator)
			object.issuer = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
