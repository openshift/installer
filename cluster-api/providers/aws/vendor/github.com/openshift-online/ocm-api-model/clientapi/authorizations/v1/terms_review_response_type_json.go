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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalTermsReviewResponse writes a value of the 'terms_review_response' type to the given writer.
func MarshalTermsReviewResponse(object *TermsReviewResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteTermsReviewResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteTermsReviewResponse writes a value of the 'terms_review_response' type to the given stream.
func WriteTermsReviewResponse(object *TermsReviewResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountId)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("redirect_url")
		stream.WriteString(object.redirectUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("terms_available")
		stream.WriteBool(object.termsAvailable)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("terms_required")
		stream.WriteBool(object.termsRequired)
	}
	stream.WriteObjectEnd()
}

// UnmarshalTermsReviewResponse reads a value of the 'terms_review_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalTermsReviewResponse(source interface{}) (object *TermsReviewResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadTermsReviewResponse(iterator)
	err = iterator.Error
	return
}

// ReadTermsReviewResponse reads a value of the 'terms_review_response' type from the given iterator.
func ReadTermsReviewResponse(iterator *jsoniter.Iterator) *TermsReviewResponse {
	object := &TermsReviewResponse{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "account_id":
			value := iterator.ReadString()
			object.accountId = value
			object.fieldSet_[0] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.fieldSet_[1] = true
		case "redirect_url":
			value := iterator.ReadString()
			object.redirectUrl = value
			object.fieldSet_[2] = true
		case "terms_available":
			value := iterator.ReadBool()
			object.termsAvailable = value
			object.fieldSet_[3] = true
		case "terms_required":
			value := iterator.ReadBool()
			object.termsRequired = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
