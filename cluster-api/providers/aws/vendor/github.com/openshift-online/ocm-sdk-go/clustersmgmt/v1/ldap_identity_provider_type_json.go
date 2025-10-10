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

// MarshalLDAPIdentityProvider writes a value of the 'LDAP_identity_provider' type to the given writer.
func MarshalLDAPIdentityProvider(object *LDAPIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLDAPIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLDAPIdentityProvider writes a value of the 'LDAP_identity_provider' type to the given stream.
func WriteLDAPIdentityProvider(object *LDAPIdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ca")
		stream.WriteString(object.ca)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("url")
		stream.WriteString(object.url)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.attributes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("attributes")
		WriteLDAPAttributes(object.attributes, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bind_dn")
		stream.WriteString(object.bindDN)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bind_password")
		stream.WriteString(object.bindPassword)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("insecure")
		stream.WriteBool(object.insecure)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLDAPIdentityProvider reads a value of the 'LDAP_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLDAPIdentityProvider(source interface{}) (object *LDAPIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLDAPIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadLDAPIdentityProvider reads a value of the 'LDAP_identity_provider' type from the given iterator.
func ReadLDAPIdentityProvider(iterator *jsoniter.Iterator) *LDAPIdentityProvider {
	object := &LDAPIdentityProvider{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "ca":
			value := iterator.ReadString()
			object.ca = value
			object.bitmap_ |= 1
		case "url":
			value := iterator.ReadString()
			object.url = value
			object.bitmap_ |= 2
		case "attributes":
			value := ReadLDAPAttributes(iterator)
			object.attributes = value
			object.bitmap_ |= 4
		case "bind_dn":
			value := iterator.ReadString()
			object.bindDN = value
			object.bitmap_ |= 8
		case "bind_password":
			value := iterator.ReadString()
			object.bindPassword = value
			object.bitmap_ |= 16
		case "insecure":
			value := iterator.ReadBool()
			object.insecure = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
