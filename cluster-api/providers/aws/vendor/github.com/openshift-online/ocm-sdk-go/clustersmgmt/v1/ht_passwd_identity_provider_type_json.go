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

// MarshalHTPasswdIdentityProvider writes a value of the 'HT_passwd_identity_provider' type to the given writer.
func MarshalHTPasswdIdentityProvider(object *HTPasswdIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteHTPasswdIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteHTPasswdIdentityProvider writes a value of the 'HT_passwd_identity_provider' type to the given stream.
func WriteHTPasswdIdentityProvider(object *HTPasswdIdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("password")
		stream.WriteString(object.password)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.users != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("users")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteHTPasswdUserList(object.users.Items(), stream)
		stream.WriteObjectEnd()
	}
	stream.WriteObjectEnd()
}

// UnmarshalHTPasswdIdentityProvider reads a value of the 'HT_passwd_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalHTPasswdIdentityProvider(source interface{}) (object *HTPasswdIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadHTPasswdIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadHTPasswdIdentityProvider reads a value of the 'HT_passwd_identity_provider' type from the given iterator.
func ReadHTPasswdIdentityProvider(iterator *jsoniter.Iterator) *HTPasswdIdentityProvider {
	object := &HTPasswdIdentityProvider{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "password":
			value := iterator.ReadString()
			object.password = value
			object.bitmap_ |= 1
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.bitmap_ |= 2
		case "users":
			value := &HTPasswdUserList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == HTPasswdUserListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadHTPasswdUserList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.users = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
