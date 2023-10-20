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

// MarshalTermsReviewResponseList writes a list of values of the 'terms_review_response' type to
// the given writer.
func MarshalTermsReviewResponseList(list []*TermsReviewResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeTermsReviewResponseList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeTermsReviewResponseList writes a list of value of the 'terms_review_response' type to
// the given stream.
func writeTermsReviewResponseList(list []*TermsReviewResponse, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeTermsReviewResponse(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalTermsReviewResponseList reads a list of values of the 'terms_review_response' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalTermsReviewResponseList(source interface{}) (items []*TermsReviewResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readTermsReviewResponseList(iterator)
	err = iterator.Error
	return
}

// readTermsReviewResponseList reads list of values of the ”terms_review_response' type from
// the given iterator.
func readTermsReviewResponseList(iterator *jsoniter.Iterator) []*TermsReviewResponse {
	list := []*TermsReviewResponse{}
	for iterator.ReadArray() {
		item := readTermsReviewResponse(iterator)
		list = append(list, item)
	}
	return list
}
