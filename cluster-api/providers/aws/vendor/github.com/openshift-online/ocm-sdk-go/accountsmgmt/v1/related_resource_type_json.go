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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalRelatedResource writes a value of the 'related_resource' type to the given writer.
func MarshalRelatedResource(object *RelatedResource, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeRelatedResource(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeRelatedResource writes a value of the 'related_resource' type to the given stream.
func writeRelatedResource(object *RelatedResource, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byoc")
		stream.WriteString(object.byoc)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone_type")
		stream.WriteString(object.availabilityZoneType)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		stream.WriteString(object.billingModel)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cost")
		stream.WriteInt(object.cost)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		stream.WriteString(object.product)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRelatedResource reads a value of the 'related_resource' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRelatedResource(source interface{}) (object *RelatedResource, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readRelatedResource(iterator)
	err = iterator.Error
	return
}

// readRelatedResource reads a value of the 'related_resource' type from the given iterator.
func readRelatedResource(iterator *jsoniter.Iterator) *RelatedResource {
	object := &RelatedResource{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "byoc":
			value := iterator.ReadString()
			object.byoc = value
			object.bitmap_ |= 1
		case "availability_zone_type":
			value := iterator.ReadString()
			object.availabilityZoneType = value
			object.bitmap_ |= 2
		case "billing_model":
			value := iterator.ReadString()
			object.billingModel = value
			object.bitmap_ |= 4
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.bitmap_ |= 8
		case "cost":
			value := iterator.ReadInt()
			object.cost = value
			object.bitmap_ |= 16
		case "product":
			value := iterator.ReadString()
			object.product = value
			object.bitmap_ |= 32
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.bitmap_ |= 64
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
