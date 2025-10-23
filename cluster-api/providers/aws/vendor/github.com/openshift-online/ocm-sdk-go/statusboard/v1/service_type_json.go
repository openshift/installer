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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalService writes a value of the 'service' type to the given writer.
func MarshalService(object *Service, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteService(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteService writes a value of the 'service' type to the given stream.
func WriteService(object *Service, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ServiceLinkKind)
	} else {
		stream.WriteString(ServiceKind)
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
	present_ = object.bitmap_&8 != 0 && object.application != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("application")
		WriteApplication(object.application, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("current_status")
		stream.WriteString(object.currentStatus)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("fullname")
		stream.WriteString(object.fullname)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_ping_at")
		stream.WriteString((object.lastPingAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("metadata")
		stream.WriteVal(object.metadata)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.owners != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("owners")
		WriteOwnerList(object.owners, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private")
		stream.WriteBool(object.private)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_endpoint")
		stream.WriteString(object.serviceEndpoint)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_type")
		stream.WriteString(object.statusType)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_updated_at")
		stream.WriteString((object.statusUpdatedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("token")
		stream.WriteString(object.token)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalService reads a value of the 'service' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalService(source interface{}) (object *Service, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadService(iterator)
	err = iterator.Error
	return
}

// ReadService reads a value of the 'service' type from the given iterator.
func ReadService(iterator *jsoniter.Iterator) *Service {
	object := &Service{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ServiceLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "application":
			value := ReadApplication(iterator)
			object.application = value
			object.bitmap_ |= 8
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 16
		case "current_status":
			value := iterator.ReadString()
			object.currentStatus = value
			object.bitmap_ |= 32
		case "fullname":
			value := iterator.ReadString()
			object.fullname = value
			object.bitmap_ |= 64
		case "last_ping_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastPingAt = value
			object.bitmap_ |= 128
		case "metadata":
			var value interface{}
			iterator.ReadVal(&value)
			object.metadata = value
			object.bitmap_ |= 256
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 512
		case "owners":
			value := ReadOwnerList(iterator)
			object.owners = value
			object.bitmap_ |= 1024
		case "private":
			value := iterator.ReadBool()
			object.private = value
			object.bitmap_ |= 2048
		case "service_endpoint":
			value := iterator.ReadString()
			object.serviceEndpoint = value
			object.bitmap_ |= 4096
		case "status_type":
			value := iterator.ReadString()
			object.statusType = value
			object.bitmap_ |= 8192
		case "status_updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.statusUpdatedAt = value
			object.bitmap_ |= 16384
		case "token":
			value := iterator.ReadString()
			object.token = value
			object.bitmap_ |= 32768
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 65536
		default:
			iterator.ReadAny()
		}
	}
	return object
}
