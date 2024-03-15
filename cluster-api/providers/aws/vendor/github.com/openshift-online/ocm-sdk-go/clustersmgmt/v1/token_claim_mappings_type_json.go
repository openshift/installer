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

// MarshalTokenClaimMappings writes a value of the 'token_claim_mappings' type to the given writer.
func MarshalTokenClaimMappings(object *TokenClaimMappings, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeTokenClaimMappings(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeTokenClaimMappings writes a value of the 'token_claim_mappings' type to the given stream.
func writeTokenClaimMappings(object *TokenClaimMappings, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		writeGroupsClaim(object.groups, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.userName != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("user_name")
		writeUsernameClaim(object.userName, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalTokenClaimMappings reads a value of the 'token_claim_mappings' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalTokenClaimMappings(source interface{}) (object *TokenClaimMappings, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readTokenClaimMappings(iterator)
	err = iterator.Error
	return
}

// readTokenClaimMappings reads a value of the 'token_claim_mappings' type from the given iterator.
func readTokenClaimMappings(iterator *jsoniter.Iterator) *TokenClaimMappings {
	object := &TokenClaimMappings{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "groups":
			value := readGroupsClaim(iterator)
			object.groups = value
			object.bitmap_ |= 1
		case "user_name":
			value := readUsernameClaim(iterator)
			object.userName = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
