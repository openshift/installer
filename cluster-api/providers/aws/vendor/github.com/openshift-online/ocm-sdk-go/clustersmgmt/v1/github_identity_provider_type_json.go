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

// MarshalGithubIdentityProvider writes a value of the 'github_identity_provider' type to the given writer.
func MarshalGithubIdentityProvider(object *GithubIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGithubIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGithubIdentityProvider writes a value of the 'github_identity_provider' type to the given stream.
func WriteGithubIdentityProvider(object *GithubIdentityProvider, stream *jsoniter.Stream) {
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
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_secret")
		stream.WriteString(object.clientSecret)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hostname")
		stream.WriteString(object.hostname)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.organizations != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organizations")
		WriteStringList(object.organizations, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.teams != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("teams")
		WriteStringList(object.teams, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGithubIdentityProvider reads a value of the 'github_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGithubIdentityProvider(source interface{}) (object *GithubIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGithubIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadGithubIdentityProvider reads a value of the 'github_identity_provider' type from the given iterator.
func ReadGithubIdentityProvider(iterator *jsoniter.Iterator) *GithubIdentityProvider {
	object := &GithubIdentityProvider{}
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
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.bitmap_ |= 2
		case "client_secret":
			value := iterator.ReadString()
			object.clientSecret = value
			object.bitmap_ |= 4
		case "hostname":
			value := iterator.ReadString()
			object.hostname = value
			object.bitmap_ |= 8
		case "organizations":
			value := ReadStringList(iterator)
			object.organizations = value
			object.bitmap_ |= 16
		case "teams":
			value := ReadStringList(iterator)
			object.teams = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
