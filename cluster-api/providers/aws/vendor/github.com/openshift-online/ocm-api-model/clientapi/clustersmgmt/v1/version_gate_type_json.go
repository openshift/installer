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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalVersionGate writes a value of the 'version_gate' type to the given writer.
func MarshalVersionGate(object *VersionGate, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteVersionGate(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteVersionGate writes a value of the 'version_gate' type to the given stream.
func WriteVersionGate(object *VersionGate, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(VersionGateLinkKind)
	} else {
		stream.WriteString(VersionGateKind)
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
		stream.WriteObjectField("sts_only")
		stream.WriteBool(object.stsOnly)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_condition")
		stream.WriteString(object.clusterCondition)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("documentation_url")
		stream.WriteString(object.documentationURL)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("label")
		stream.WriteString(object.label)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.value)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version_raw_id_prefix")
		stream.WriteString(object.versionRawIDPrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
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
	object = ReadVersionGate(iterator)
	err = iterator.Error
	return
}

// ReadVersionGate reads a value of the 'version_gate' type from the given iterator.
func ReadVersionGate(iterator *jsoniter.Iterator) *VersionGate {
	object := &VersionGate{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == VersionGateLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "sts_only":
			value := iterator.ReadBool()
			object.stsOnly = value
			object.fieldSet_[3] = true
		case "cluster_condition":
			value := iterator.ReadString()
			object.clusterCondition = value
			object.fieldSet_[4] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[5] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[6] = true
		case "documentation_url":
			value := iterator.ReadString()
			object.documentationURL = value
			object.fieldSet_[7] = true
		case "label":
			value := iterator.ReadString()
			object.label = value
			object.fieldSet_[8] = true
		case "value":
			value := iterator.ReadString()
			object.value = value
			object.fieldSet_[9] = true
		case "version_raw_id_prefix":
			value := iterator.ReadString()
			object.versionRawIDPrefix = value
			object.fieldSet_[10] = true
		case "warning_message":
			value := iterator.ReadString()
			object.warningMessage = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
