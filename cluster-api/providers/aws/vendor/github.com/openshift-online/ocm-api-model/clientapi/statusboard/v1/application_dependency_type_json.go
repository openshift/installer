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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalApplicationDependency writes a value of the 'application_dependency' type to the given writer.
func MarshalApplicationDependency(object *ApplicationDependency, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteApplicationDependency(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteApplicationDependency writes a value of the 'application_dependency' type to the given stream.
func WriteApplicationDependency(object *ApplicationDependency, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ApplicationDependencyLinkKind)
	} else {
		stream.WriteString(ApplicationDependencyKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.application != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("application")
		WriteApplication(object.application, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("metadata")
		stream.WriteVal(object.metadata)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.owners != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("owners")
		WriteOwnerList(object.owners, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.service != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service")
		WriteService(object.service, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
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
	object = ReadApplicationDependency(iterator)
	err = iterator.Error
	return
}

// ReadApplicationDependency reads a value of the 'application_dependency' type from the given iterator.
func ReadApplicationDependency(iterator *jsoniter.Iterator) *ApplicationDependency {
	object := &ApplicationDependency{
		fieldSet_: make([]bool, 11),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ApplicationDependencyLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "application":
			value := ReadApplication(iterator)
			object.application = value
			object.fieldSet_[3] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[4] = true
		case "metadata":
			var value interface{}
			iterator.ReadVal(&value)
			object.metadata = value
			object.fieldSet_[5] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[6] = true
		case "owners":
			value := ReadOwnerList(iterator)
			object.owners = value
			object.fieldSet_[7] = true
		case "service":
			value := ReadService(iterator)
			object.service = value
			object.fieldSet_[8] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[9] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[10] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
