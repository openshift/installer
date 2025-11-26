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

// MarshalFeatureReviewResponse writes a value of the 'feature_review_response' type to the given writer.
func MarshalFeatureReviewResponse(object *FeatureReviewResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteFeatureReviewResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteFeatureReviewResponse writes a value of the 'feature_review_response' type to the given stream.
func WriteFeatureReviewResponse(object *FeatureReviewResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("feature_id")
		stream.WriteString(object.featureID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalFeatureReviewResponse reads a value of the 'feature_review_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalFeatureReviewResponse(source interface{}) (object *FeatureReviewResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadFeatureReviewResponse(iterator)
	err = iterator.Error
	return
}

// ReadFeatureReviewResponse reads a value of the 'feature_review_response' type from the given iterator.
func ReadFeatureReviewResponse(iterator *jsoniter.Iterator) *FeatureReviewResponse {
	object := &FeatureReviewResponse{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[0] = true
		case "feature_id":
			value := iterator.ReadString()
			object.featureID = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
