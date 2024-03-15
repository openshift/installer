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

// MarshalExternalAuthClaim writes a value of the 'external_auth_claim' type to the given writer.
func MarshalExternalAuthClaim(object *ExternalAuthClaim, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeExternalAuthClaim(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeExternalAuthClaim writes a value of the 'external_auth_claim' type to the given stream.
func writeExternalAuthClaim(object *ExternalAuthClaim, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.mappings != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("mappings")
		writeTokenClaimMappings(object.mappings, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.validationRules != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("validation_rules")
		writeTokenClaimValidationRuleList(object.validationRules, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalExternalAuthClaim reads a value of the 'external_auth_claim' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalExternalAuthClaim(source interface{}) (object *ExternalAuthClaim, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readExternalAuthClaim(iterator)
	err = iterator.Error
	return
}

// readExternalAuthClaim reads a value of the 'external_auth_claim' type from the given iterator.
func readExternalAuthClaim(iterator *jsoniter.Iterator) *ExternalAuthClaim {
	object := &ExternalAuthClaim{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "mappings":
			value := readTokenClaimMappings(iterator)
			object.mappings = value
			object.bitmap_ |= 1
		case "validation_rules":
			value := readTokenClaimValidationRuleList(iterator)
			object.validationRules = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
