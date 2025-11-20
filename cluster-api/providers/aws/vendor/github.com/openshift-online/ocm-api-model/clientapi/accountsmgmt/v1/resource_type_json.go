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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalResource writes a value of the 'resource' type to the given writer.
func MarshalResource(object *Resource, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteResource(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteResource writes a value of the 'resource' type to the given stream.
func WriteResource(object *Resource, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ResourceLinkKind)
	} else {
		stream.WriteString(ResourceKind)
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
		stream.WriteObjectField("byoc")
		stream.WriteBool(object.byoc)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sku")
		stream.WriteString(object.sku)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed")
		stream.WriteInt(object.allowed)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone_type")
		stream.WriteString(object.availabilityZoneType)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalResource reads a value of the 'resource' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalResource(source interface{}) (object *Resource, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadResource(iterator)
	err = iterator.Error
	return
}

// ReadResource reads a value of the 'resource' type from the given iterator.
func ReadResource(iterator *jsoniter.Iterator) *Resource {
	object := &Resource{
		fieldSet_: make([]bool, 9),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ResourceLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "byoc":
			value := iterator.ReadBool()
			object.byoc = value
			object.fieldSet_[3] = true
		case "sku":
			value := iterator.ReadString()
			object.sku = value
			object.fieldSet_[4] = true
		case "allowed":
			value := iterator.ReadInt()
			object.allowed = value
			object.fieldSet_[5] = true
		case "availability_zone_type":
			value := iterator.ReadString()
			object.availabilityZoneType = value
			object.fieldSet_[6] = true
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.fieldSet_[7] = true
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.fieldSet_[8] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
