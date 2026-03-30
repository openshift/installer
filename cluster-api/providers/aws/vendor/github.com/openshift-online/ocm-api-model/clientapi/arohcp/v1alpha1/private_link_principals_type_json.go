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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalPrivateLinkPrincipals writes a value of the 'private_link_principals' type to the given writer.
func MarshalPrivateLinkPrincipals(object *PrivateLinkPrincipals, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WritePrivateLinkPrincipals(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WritePrivateLinkPrincipals writes a value of the 'private_link_principals' type to the given stream.
func WritePrivateLinkPrincipals(object *PrivateLinkPrincipals, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(PrivateLinkPrincipalsLinkKind)
	} else {
		stream.WriteString(PrivateLinkPrincipalsKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.principals != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principals")
		WritePrivateLinkPrincipalList(object.principals, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalPrivateLinkPrincipals reads a value of the 'private_link_principals' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalPrivateLinkPrincipals(source interface{}) (object *PrivateLinkPrincipals, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadPrivateLinkPrincipals(iterator)
	err = iterator.Error
	return
}

// ReadPrivateLinkPrincipals reads a value of the 'private_link_principals' type from the given iterator.
func ReadPrivateLinkPrincipals(iterator *jsoniter.Iterator) *PrivateLinkPrincipals {
	object := &PrivateLinkPrincipals{
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
			if value == PrivateLinkPrincipalsLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "principals":
			value := ReadPrivateLinkPrincipalList(iterator)
			object.principals = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
