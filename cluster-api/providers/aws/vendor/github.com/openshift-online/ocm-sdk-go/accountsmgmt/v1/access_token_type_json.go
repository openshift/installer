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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAccessToken writes a value of the 'access_token' type to the given writer.
func MarshalAccessToken(object *AccessToken, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAccessToken(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAccessToken writes a value of the 'access_token' type to the given stream.
func writeAccessToken(object *AccessToken, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.auths != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auths")
		if object.auths != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.auths))
			i := 0
			for key := range object.auths {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.auths[key]
				stream.WriteObjectField(key)
				writeAccessTokenAuth(item, stream)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessToken reads a value of the 'access_token' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessToken(source interface{}) (object *AccessToken, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAccessToken(iterator)
	err = iterator.Error
	return
}

// readAccessToken reads a value of the 'access_token' type from the given iterator.
func readAccessToken(iterator *jsoniter.Iterator) *AccessToken {
	object := &AccessToken{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "auths":
			value := map[string]*AccessTokenAuth{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := readAccessTokenAuth(iterator)
				value[key] = item
			}
			object.auths = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
