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

// MarshalGoogleIdentityProvider writes a value of the 'google_identity_provider' type to the given writer.
func MarshalGoogleIdentityProvider(object *GoogleIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGoogleIdentityProvider(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGoogleIdentityProvider writes a value of the 'google_identity_provider' type to the given stream.
func WriteGoogleIdentityProvider(object *GoogleIdentityProvider, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_secret")
		stream.WriteString(object.clientSecret)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hosted_domain")
		stream.WriteString(object.hostedDomain)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGoogleIdentityProvider reads a value of the 'google_identity_provider' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGoogleIdentityProvider(source interface{}) (object *GoogleIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGoogleIdentityProvider(iterator)
	err = iterator.Error
	return
}

// ReadGoogleIdentityProvider reads a value of the 'google_identity_provider' type from the given iterator.
func ReadGoogleIdentityProvider(iterator *jsoniter.Iterator) *GoogleIdentityProvider {
	object := &GoogleIdentityProvider{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.fieldSet_[0] = true
		case "client_secret":
			value := iterator.ReadString()
			object.clientSecret = value
			object.fieldSet_[1] = true
		case "hosted_domain":
			value := iterator.ReadString()
			object.hostedDomain = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
