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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSelfTermsReviewRequest writes a value of the 'self_terms_review_request' type to the given writer.
func MarshalSelfTermsReviewRequest(object *SelfTermsReviewRequest, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSelfTermsReviewRequest(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSelfTermsReviewRequest writes a value of the 'self_terms_review_request' type to the given stream.
func writeSelfTermsReviewRequest(object *SelfTermsReviewRequest, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_code")
		stream.WriteString(object.eventCode)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("site_code")
		stream.WriteString(object.siteCode)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSelfTermsReviewRequest reads a value of the 'self_terms_review_request' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSelfTermsReviewRequest(source interface{}) (object *SelfTermsReviewRequest, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSelfTermsReviewRequest(iterator)
	err = iterator.Error
	return
}

// readSelfTermsReviewRequest reads a value of the 'self_terms_review_request' type from the given iterator.
func readSelfTermsReviewRequest(iterator *jsoniter.Iterator) *SelfTermsReviewRequest {
	object := &SelfTermsReviewRequest{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "event_code":
			value := iterator.ReadString()
			object.eventCode = value
			object.bitmap_ |= 1
		case "site_code":
			value := iterator.ReadString()
			object.siteCode = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
