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

// MarshalImageMirror writes a value of the 'image_mirror' type to the given writer.
func MarshalImageMirror(object *ImageMirror, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteImageMirror(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteImageMirror writes a value of the 'image_mirror' type to the given stream.
func WriteImageMirror(object *ImageMirror, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ImageMirrorLinkKind)
	} else {
		stream.WriteString(ImageMirrorKind)
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
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.mirrors != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("mirrors")
		WriteStringList(object.mirrors, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("source")
		stream.WriteString(object.source)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalImageMirror reads a value of the 'image_mirror' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalImageMirror(source interface{}) (object *ImageMirror, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadImageMirror(iterator)
	err = iterator.Error
	return
}

// ReadImageMirror reads a value of the 'image_mirror' type from the given iterator.
func ReadImageMirror(iterator *jsoniter.Iterator) *ImageMirror {
	object := &ImageMirror{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ImageMirrorLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[3] = true
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.fieldSet_[4] = true
		case "mirrors":
			value := ReadStringList(iterator)
			object.mirrors = value
			object.fieldSet_[5] = true
		case "source":
			value := iterator.ReadString()
			object.source = value
			object.fieldSet_[6] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
