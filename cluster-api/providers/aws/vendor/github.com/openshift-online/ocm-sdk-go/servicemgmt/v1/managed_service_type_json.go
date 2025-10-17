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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalManagedService writes a value of the 'managed_service' type to the given writer.
func MarshalManagedService(object *ManagedService, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteManagedService(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteManagedService writes a value of the 'managed_service' type to the given stream.
func WriteManagedService(object *ManagedService, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ManagedServiceLinkKind)
	} else {
		stream.WriteString(ManagedServiceKind)
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
	present_ = object.bitmap_&8 != 0 && object.addon != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon")
		WriteStatefulObject(object.addon, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.cluster != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster")
		WriteCluster(object.cluster, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("expired_at")
		stream.WriteString((object.expiredAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.parameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parameters")
		WriteServiceParameterList(object.parameters, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.resources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resources")
		WriteStatefulObjectList(object.resources, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service")
		stream.WriteString(object.service)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_state")
		stream.WriteString(object.serviceState)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalManagedService reads a value of the 'managed_service' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalManagedService(source interface{}) (object *ManagedService, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadManagedService(iterator)
	err = iterator.Error
	return
}

// ReadManagedService reads a value of the 'managed_service' type from the given iterator.
func ReadManagedService(iterator *jsoniter.Iterator) *ManagedService {
	object := &ManagedService{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ManagedServiceLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "addon":
			value := ReadStatefulObject(iterator)
			object.addon = value
			object.bitmap_ |= 8
		case "cluster":
			value := ReadCluster(iterator)
			object.cluster = value
			object.bitmap_ |= 16
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 32
		case "expired_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expiredAt = value
			object.bitmap_ |= 64
		case "parameters":
			value := ReadServiceParameterList(iterator)
			object.parameters = value
			object.bitmap_ |= 128
		case "resources":
			value := ReadStatefulObjectList(iterator)
			object.resources = value
			object.bitmap_ |= 256
		case "service":
			value := iterator.ReadString()
			object.service = value
			object.bitmap_ |= 512
		case "service_state":
			value := iterator.ReadString()
			object.serviceState = value
			object.bitmap_ |= 1024
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
