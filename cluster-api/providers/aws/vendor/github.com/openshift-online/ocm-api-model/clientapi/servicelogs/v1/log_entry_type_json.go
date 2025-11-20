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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicelogs/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalLogEntry writes a value of the 'log_entry' type to the given writer.
func MarshalLogEntry(object *LogEntry, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogEntry(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogEntry writes a value of the 'log_entry' type to the given stream.
func WriteLogEntry(object *LogEntry, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(LogEntryLinkKind)
	} else {
		stream.WriteString(LogEntryKind)
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
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_uuid")
		stream.WriteString(object.clusterUUID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_by")
		stream.WriteString(object.createdBy)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.docReferences != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("doc_references")
		WriteStringList(object.docReferences, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("event_stream_id")
		stream.WriteString(object.eventStreamID)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("internal_only")
		stream.WriteBool(object.internalOnly)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_type")
		stream.WriteString(string(object.logType))
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_name")
		stream.WriteString(object.serviceName)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(string(object.severity))
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("summary")
		stream.WriteString(object.summary)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("timestamp")
		stream.WriteString((object.timestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogEntry reads a value of the 'log_entry' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogEntry(source interface{}) (object *LogEntry, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogEntry(iterator)
	err = iterator.Error
	return
}

// ReadLogEntry reads a value of the 'log_entry' type from the given iterator.
func ReadLogEntry(iterator *jsoniter.Iterator) *LogEntry {
	object := &LogEntry{
		fieldSet_: make([]bool, 18),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == LogEntryLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[3] = true
		case "cluster_uuid":
			value := iterator.ReadString()
			object.clusterUUID = value
			object.fieldSet_[4] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[5] = true
		case "created_by":
			value := iterator.ReadString()
			object.createdBy = value
			object.fieldSet_[6] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[7] = true
		case "doc_references":
			value := ReadStringList(iterator)
			object.docReferences = value
			object.fieldSet_[8] = true
		case "event_stream_id":
			value := iterator.ReadString()
			object.eventStreamID = value
			object.fieldSet_[9] = true
		case "internal_only":
			value := iterator.ReadBool()
			object.internalOnly = value
			object.fieldSet_[10] = true
		case "log_type":
			text := iterator.ReadString()
			value := LogType(text)
			object.logType = value
			object.fieldSet_[11] = true
		case "service_name":
			value := iterator.ReadString()
			object.serviceName = value
			object.fieldSet_[12] = true
		case "severity":
			text := iterator.ReadString()
			value := Severity(text)
			object.severity = value
			object.fieldSet_[13] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.fieldSet_[14] = true
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.fieldSet_[15] = true
		case "timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.timestamp = value
			object.fieldSet_[16] = true
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.fieldSet_[17] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
