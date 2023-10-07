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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalCloudResource writes a value of the 'cloud_resource' type to the given writer.
func MarshalCloudResource(object *CloudResource, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeCloudResource(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeCloudResource writes a value of the 'cloud_resource' type to the given stream.
func writeCloudResource(object *CloudResource, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(CloudResourceLinkKind)
	} else {
		stream.WriteString(CloudResourceKind)
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
		stream.WriteObjectField("active")
		stream.WriteBool(object.active)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("category")
		stream.WriteString(object.category)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("category_pretty")
		stream.WriteString(object.categoryPretty)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu_cores")
		stream.WriteInt(object.cpuCores)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("generic_name")
		stream.WriteString(object.genericName)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		stream.WriteInt(object.memory)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory_pretty")
		stream.WriteString(object.memoryPretty)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name_pretty")
		stream.WriteString(object.namePretty)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_type")
		stream.WriteString(object.resourceType)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("size_pretty")
		stream.WriteString(object.sizePretty)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalCloudResource reads a value of the 'cloud_resource' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCloudResource(source interface{}) (object *CloudResource, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readCloudResource(iterator)
	err = iterator.Error
	return
}

// readCloudResource reads a value of the 'cloud_resource' type from the given iterator.
func readCloudResource(iterator *jsoniter.Iterator) *CloudResource {
	object := &CloudResource{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == CloudResourceLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "active":
			value := iterator.ReadBool()
			object.active = value
			object.bitmap_ |= 8
		case "category":
			value := iterator.ReadString()
			object.category = value
			object.bitmap_ |= 16
		case "category_pretty":
			value := iterator.ReadString()
			object.categoryPretty = value
			object.bitmap_ |= 32
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.bitmap_ |= 64
		case "cpu_cores":
			value := iterator.ReadInt()
			object.cpuCores = value
			object.bitmap_ |= 128
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 256
		case "generic_name":
			value := iterator.ReadString()
			object.genericName = value
			object.bitmap_ |= 512
		case "memory":
			value := iterator.ReadInt()
			object.memory = value
			object.bitmap_ |= 1024
		case "memory_pretty":
			value := iterator.ReadString()
			object.memoryPretty = value
			object.bitmap_ |= 2048
		case "name_pretty":
			value := iterator.ReadString()
			object.namePretty = value
			object.bitmap_ |= 4096
		case "resource_type":
			value := iterator.ReadString()
			object.resourceType = value
			object.bitmap_ |= 8192
		case "size_pretty":
			value := iterator.ReadString()
			object.sizePretty = value
			object.bitmap_ |= 16384
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
