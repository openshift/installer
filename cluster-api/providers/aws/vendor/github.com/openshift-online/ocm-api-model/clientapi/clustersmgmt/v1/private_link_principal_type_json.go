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

// MarshalPrivateLinkPrincipal writes a value of the 'private_link_principal' type to the given writer.
func MarshalPrivateLinkPrincipal(object *PrivateLinkPrincipal, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WritePrivateLinkPrincipal(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WritePrivateLinkPrincipal writes a value of the 'private_link_principal' type to the given stream.
func WritePrivateLinkPrincipal(object *PrivateLinkPrincipal, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(PrivateLinkPrincipalLinkKind)
	} else {
		stream.WriteString(PrivateLinkPrincipalKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principal")
		stream.WriteString(object.principal)
	}
	stream.WriteObjectEnd()
}

// UnmarshalPrivateLinkPrincipal reads a value of the 'private_link_principal' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalPrivateLinkPrincipal(source interface{}) (object *PrivateLinkPrincipal, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadPrivateLinkPrincipal(iterator)
	err = iterator.Error
	return
}

// ReadPrivateLinkPrincipal reads a value of the 'private_link_principal' type from the given iterator.
func ReadPrivateLinkPrincipal(iterator *jsoniter.Iterator) *PrivateLinkPrincipal {
	object := &PrivateLinkPrincipal{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == PrivateLinkPrincipalLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "principal":
			value := iterator.ReadString()
			object.principal = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
