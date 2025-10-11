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

// MarshalAdditionalCatalogSource writes a value of the 'additional_catalog_source' type to the given writer.
func MarshalAdditionalCatalogSource(object *AdditionalCatalogSource, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAdditionalCatalogSource(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAdditionalCatalogSource writes a value of the 'additional_catalog_source' type to the given stream.
func WriteAdditionalCatalogSource(object *AdditionalCatalogSource, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("image")
		stream.WriteString(object.image)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAdditionalCatalogSource reads a value of the 'additional_catalog_source' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAdditionalCatalogSource(source interface{}) (object *AdditionalCatalogSource, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAdditionalCatalogSource(iterator)
	err = iterator.Error
	return
}

// ReadAdditionalCatalogSource reads a value of the 'additional_catalog_source' type from the given iterator.
func ReadAdditionalCatalogSource(iterator *jsoniter.Iterator) *AdditionalCatalogSource {
	object := &AdditionalCatalogSource{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 1
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 2
		case "image":
			value := iterator.ReadString()
			object.image = value
			object.bitmap_ |= 4
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
