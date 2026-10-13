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
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ca")
		stream.WriteString(object.ca)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("url")
		stream.WriteString(object.url)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.attributes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("attributes")
		WriteLDAPAttributes(object.attributes, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bind_dn")
		stream.WriteString(object.bindDN)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bind_password")
		stream.WriteString(object.bindPassword)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
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
	object := &LDAPIdentityProvider{
		fieldSet_: make([]bool, 6),
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
		case "url":
			value := iterator.ReadString()
			object.url = value
			object.fieldSet_[1] = true
		case "attributes":
			value := ReadLDAPAttributes(iterator)
			object.attributes = value
			object.fieldSet_[2] = true
		case "bind_dn":
			value := iterator.ReadString()
			object.bindDN = value
			object.fieldSet_[3] = true
		case "bind_password":
			value := iterator.ReadString()
			object.bindPassword = value
			object.fieldSet_[4] = true
		case "insecure":
			value := iterator.ReadBool()
			object.insecure = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
