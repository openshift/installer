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

// MarshalImageOverrides writes a value of the 'image_overrides' type to the given writer.
func MarshalImageOverrides(object *ImageOverrides, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteImageOverrides(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteImageOverrides writes a value of the 'image_overrides' type to the given stream.
func WriteImageOverrides(object *ImageOverrides, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ImageOverridesLinkKind)
	} else {
		stream.WriteString(ImageOverridesKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAMIOverrideList(object.aws, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCPImageOverrideList(object.gcp, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalImageOverrides reads a value of the 'image_overrides' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalImageOverrides(source interface{}) (object *ImageOverrides, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadImageOverrides(iterator)
	err = iterator.Error
	return
}

// ReadImageOverrides reads a value of the 'image_overrides' type from the given iterator.
func ReadImageOverrides(iterator *jsoniter.Iterator) *ImageOverrides {
	object := &ImageOverrides{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ImageOverridesLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "aws":
			value := ReadAMIOverrideList(iterator)
			object.aws = value
			object.fieldSet_[3] = true
		case "gcp":
			value := ReadGCPImageOverrideList(iterator)
			object.gcp = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
