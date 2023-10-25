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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalVersionGate writes a value of the 'version_gate' type to the given writer.
func MarshalVersionGate(object *VersionGate, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeVersionGate(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeVersionGate writes a value of the 'version_gate' type to the given stream.
func writeVersionGate(object *VersionGate, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(VersionGateLinkKind)
	} else {
		stream.WriteString(VersionGateKind)
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
		stream.WriteObjectField("sts_only")
		stream.WriteBool(object.stsOnly)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("documentation_url")
		stream.WriteString(object.documentationURL)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("label")
		stream.WriteString(object.label)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.value)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version_raw_id_prefix")
		stream.WriteString(object.versionRawIDPrefix)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("warning_message")
		stream.WriteString(object.warningMessage)
	}
	stream.WriteObjectEnd()
}

// UnmarshalVersionGate reads a value of the 'version_gate' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalVersionGate(source interface{}) (object *VersionGate, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readVersionGate(iterator)
	err = iterator.Error
	return
}

// readVersionGate reads a value of the 'version_gate' type from the given iterator.
func readVersionGate(iterator *jsoniter.Iterator) *VersionGate {
	object := &VersionGate{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == VersionGateLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "sts_only":
			value := iterator.ReadBool()
			object.stsOnly = value
			object.bitmap_ |= 8
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 16
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 32
		case "documentation_url":
			value := iterator.ReadString()
			object.documentationURL = value
			object.bitmap_ |= 64
		case "label":
			value := iterator.ReadString()
			object.label = value
			object.bitmap_ |= 128
		case "value":
			value := iterator.ReadString()
			object.value = value
			object.bitmap_ |= 256
		case "version_raw_id_prefix":
			value := iterator.ReadString()
			object.versionRawIDPrefix = value
			object.bitmap_ |= 512
		case "warning_message":
			value := iterator.ReadString()
			object.warningMessage = value
			object.bitmap_ |= 1024
		default:
			iterator.ReadAny()
		}
	}
	return object
}
