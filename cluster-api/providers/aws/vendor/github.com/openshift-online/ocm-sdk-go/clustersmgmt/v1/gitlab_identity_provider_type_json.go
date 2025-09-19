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

// MarshalGitlabIdentityProvider writes a value of the 'gitlab_identity_provider' type to the given writer.
func MarshalGitlabIdentityProvider(object *GitlabIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGitlabIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGitlabIdentityProvider writes a value of the 'gitlab_identity_provider' type to the given stream.
func WriteGitlabIdentityProvider(object *GitlabIdentityProvider, stream *jsoniter.Stream) {
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
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_secret")
		stream.WriteString(object.clientSecret)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGitlabIdentityProvider reads a value of the 'gitlab_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGitlabIdentityProvider(source interface{}) (object *GitlabIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGitlabIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadGitlabIdentityProvider reads a value of the 'gitlab_identity_provider' type from the given iterator.
func ReadGitlabIdentityProvider(iterator *jsoniter.Iterator) *GitlabIdentityProvider {
	object := &GitlabIdentityProvider{}
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
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.bitmap_ |= 4
		case "client_secret":
			value := iterator.ReadString()
			object.clientSecret = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
