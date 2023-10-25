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

// MarshalImageOverrides writes a value of the 'image_overrides' type to the given writer.
func MarshalImageOverrides(object *ImageOverrides, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeImageOverrides(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeImageOverrides writes a value of the 'image_overrides' type to the given stream.
func writeImageOverrides(object *ImageOverrides, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ImageOverridesLinkKind)
	} else {
		stream.WriteString(ImageOverridesKind)
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
	present_ = object.bitmap_&8 != 0 && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		writeAMIOverrideList(object.aws, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		writeGCPImageOverrideList(object.gcp, stream)
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
	object = readImageOverrides(iterator)
	err = iterator.Error
	return
}

// readImageOverrides reads a value of the 'image_overrides' type from the given iterator.
func readImageOverrides(iterator *jsoniter.Iterator) *ImageOverrides {
	object := &ImageOverrides{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ImageOverridesLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "aws":
			value := readAMIOverrideList(iterator)
			object.aws = value
			object.bitmap_ |= 8
		case "gcp":
			value := readGCPImageOverrideList(iterator)
			object.gcp = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
