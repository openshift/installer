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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalWifConfig writes a value of the 'wif_config' type to the given writer.
func MarshalWifConfig(object *WifConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifConfig writes a value of the 'wif_config' type to the given stream.
func WriteWifConfig(object *WifConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(WifConfigLinkKind)
	} else {
		stream.WriteString(WifConfigKind)
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
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteWifGcp(object.gcp, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		WriteOrganizationLink(object.organization, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.wifTemplates != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("wif_templates")
		WriteStringList(object.wifTemplates, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifConfig reads a value of the 'wif_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifConfig(source interface{}) (object *WifConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifConfig(iterator)
	err = iterator.Error
	return
}

// ReadWifConfig reads a value of the 'wif_config' type from the given iterator.
func ReadWifConfig(iterator *jsoniter.Iterator) *WifConfig {
	object := &WifConfig{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == WifConfigLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.fieldSet_[3] = true
		case "gcp":
			value := ReadWifGcp(iterator)
			object.gcp = value
			object.fieldSet_[4] = true
		case "organization":
			value := ReadOrganizationLink(iterator)
			object.organization = value
			object.fieldSet_[5] = true
		case "wif_templates":
			value := ReadStringList(iterator)
			object.wifTemplates = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
