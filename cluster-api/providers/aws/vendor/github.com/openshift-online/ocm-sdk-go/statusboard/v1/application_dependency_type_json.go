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

// MarshalApplicationDependency writes a value of the 'application_dependency' type to the given writer.
func MarshalApplicationDependency(object *ApplicationDependency, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeApplicationDependency(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeApplicationDependency writes a value of the 'application_dependency' type to the given stream.
func writeApplicationDependency(object *ApplicationDependency, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ApplicationDependencyLinkKind)
	} else {
		stream.WriteString(ApplicationDependencyKind)
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
		writeApplication(object.application, stream)
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
		stream.WriteObjectField("metadata")
		stream.WriteVal(object.metadata)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.owners != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("owners")
		writeOwnerList(object.owners, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.service != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service")
		writeService(object.service, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalApplicationDependency reads a value of the 'application_dependency' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalApplicationDependency(source interface{}) (object *ApplicationDependency, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readApplicationDependency(iterator)
	err = iterator.Error
	return
}

// readApplicationDependency reads a value of the 'application_dependency' type from the given iterator.
func readApplicationDependency(iterator *jsoniter.Iterator) *ApplicationDependency {
	object := &ApplicationDependency{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ApplicationDependencyLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "application":
			value := readApplication(iterator)
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
		case "metadata":
			var value interface{}
			iterator.ReadVal(&value)
			object.metadata = value
			object.bitmap_ |= 32
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 64
		case "owners":
			value := readOwnerList(iterator)
			object.owners = value
			object.bitmap_ |= 128
		case "service":
			value := readService(iterator)
			object.service = value
			object.bitmap_ |= 256
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.bitmap_ |= 512
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 1024
		default:
			iterator.ReadAny()
		}
	}
	return object
}
