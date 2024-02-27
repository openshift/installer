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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalProductTechnologyPreview writes a value of the 'product_technology_preview' type to the given writer.
func MarshalProductTechnologyPreview(object *ProductTechnologyPreview, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeProductTechnologyPreview(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeProductTechnologyPreview writes a value of the 'product_technology_preview' type to the given stream.
func writeProductTechnologyPreview(object *ProductTechnologyPreview, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ProductTechnologyPreviewLinkKind)
	} else {
		stream.WriteString(ProductTechnologyPreviewKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_text")
		stream.WriteString(object.additionalText)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("end_date")
		stream.WriteString((object.endDate).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
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
	object = readProductTechnologyPreview(iterator)
	err = iterator.Error
	return
}

// readProductTechnologyPreview reads a value of the 'product_technology_preview' type from the given iterator.
func readProductTechnologyPreview(iterator *jsoniter.Iterator) *ProductTechnologyPreview {
	object := &ProductTechnologyPreview{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ProductTechnologyPreviewLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "additional_text":
			value := iterator.ReadString()
			object.additionalText = value
			object.bitmap_ |= 8
		case "end_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endDate = value
			object.bitmap_ |= 16
		case "start_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.startDate = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
