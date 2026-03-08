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

// MarshalIdentityProvider writes a value of the 'identity_provider' type to the given writer.
func MarshalIdentityProvider(object *IdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteIdentityProvider writes a value of the 'identity_provider' type to the given stream.
func WriteIdentityProvider(object *IdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(IdentityProviderLinkKind)
	} else {
		stream.WriteString(IdentityProviderKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.ldap != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ldap")
		WriteLDAPIdentityProvider(object.ldap, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("challenge")
		stream.WriteBool(object.challenge)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.github != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("github")
		WriteGithubIdentityProvider(object.github, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.gitlab != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gitlab")
		WriteGitlabIdentityProvider(object.gitlab, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.google != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("google")
		WriteGoogleIdentityProvider(object.google, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.htpasswd != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("htpasswd")
		WriteHTPasswdIdentityProvider(object.htpasswd, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("login")
		stream.WriteBool(object.login)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("mapping_method")
		stream.WriteString(string(object.mappingMethod))
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.openID != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("open_id")
		WriteOpenIDIdentityProvider(object.openID, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(string(object.type_))
	}
	stream.WriteObjectEnd()
}

// UnmarshalIdentityProvider reads a value of the 'identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalIdentityProvider(source interface{}) (object *IdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadIdentityProvider reads a value of the 'identity_provider' type from the given iterator.
func ReadIdentityProvider(iterator *jsoniter.Iterator) *IdentityProvider {
	object := &IdentityProvider{
		fieldSet_: make([]bool, 14),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == IdentityProviderLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "ldap":
			value := ReadLDAPIdentityProvider(iterator)
			object.ldap = value
			object.fieldSet_[3] = true
		case "challenge":
			value := iterator.ReadBool()
			object.challenge = value
			object.fieldSet_[4] = true
		case "github":
			value := ReadGithubIdentityProvider(iterator)
			object.github = value
			object.fieldSet_[5] = true
		case "gitlab":
			value := ReadGitlabIdentityProvider(iterator)
			object.gitlab = value
			object.fieldSet_[6] = true
		case "google":
			value := ReadGoogleIdentityProvider(iterator)
			object.google = value
			object.fieldSet_[7] = true
		case "htpasswd":
			value := ReadHTPasswdIdentityProvider(iterator)
			object.htpasswd = value
			object.fieldSet_[8] = true
		case "login":
			value := iterator.ReadBool()
			object.login = value
			object.fieldSet_[9] = true
		case "mapping_method":
			text := iterator.ReadString()
			value := IdentityProviderMappingMethod(text)
			object.mappingMethod = value
			object.fieldSet_[10] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[11] = true
		case "open_id":
			value := ReadOpenIDIdentityProvider(iterator)
			object.openID = value
			object.fieldSet_[12] = true
		case "type":
			text := iterator.ReadString()
			value := IdentityProviderType(text)
			object.type_ = value
			object.fieldSet_[13] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
