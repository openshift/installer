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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalProductTechnologyPreview writes a value of the 'product_technology_preview' type to the given writer.
func MarshalProductTechnologyPreview(object *ProductTechnologyPreview, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteProductTechnologyPreview(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteProductTechnologyPreview writes a value of the 'product_technology_preview' type to the given stream.
func WriteProductTechnologyPreview(object *ProductTechnologyPreview, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ProductTechnologyPreviewLinkKind)
	} else {
		stream.WriteString(ProductTechnologyPreviewKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_text")
		stream.WriteString(object.additionalText)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("end_date")
		stream.WriteString((object.endDate).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("start_date")
		stream.WriteString((object.startDate).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalProductTechnologyPreview reads a value of the 'product_technology_preview' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalProductTechnologyPreview(source interface{}) (object *ProductTechnologyPreview, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadProductTechnologyPreview(iterator)
	err = iterator.Error
	return
}

// ReadProductTechnologyPreview reads a value of the 'product_technology_preview' type from the given iterator.
func ReadProductTechnologyPreview(iterator *jsoniter.Iterator) *ProductTechnologyPreview {
	object := &ProductTechnologyPreview{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ProductTechnologyPreviewLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "additional_text":
			value := iterator.ReadString()
			object.additionalText = value
			object.fieldSet_[3] = true
		case "end_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endDate = value
			object.fieldSet_[4] = true
		case "start_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.startDate = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
