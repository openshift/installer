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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalIncident writes a value of the 'incident' type to the given writer.
func MarshalIncident(object *Incident, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeIncident(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeIncident writes a value of the 'incident' type to the given stream.
func writeIncident(object *Incident, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(IncidentLinkKind)
	} else {
		stream.WriteString(IncidentKind)
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
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creator_id")
		stream.WriteString(object.creatorId)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("deleted_at")
		stream.WriteString((object.deletedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.externalCoordination != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_coordination")
		writeStringList(object.externalCoordination, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("incident_id")
		stream.WriteString(object.incidentId)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("incident_type")
		stream.WriteString(object.incidentType)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_updated")
		stream.WriteString((object.lastUpdated).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("primary_team")
		stream.WriteString(object.primaryTeam)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("severity")
		stream.WriteString(object.severity)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("summary")
		stream.WriteString(object.summary)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("worked_at")
		stream.WriteString((object.workedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalIncident reads a value of the 'incident' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalIncident(source interface{}) (object *Incident, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readIncident(iterator)
	err = iterator.Error
	return
}

// readIncident reads a value of the 'incident' type from the given iterator.
func readIncident(iterator *jsoniter.Iterator) *Incident {
	object := &Incident{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == IncidentLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 8
		case "creator_id":
			value := iterator.ReadString()
			object.creatorId = value
			object.bitmap_ |= 16
		case "deleted_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.deletedAt = value
			object.bitmap_ |= 32
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 64
		case "external_coordination":
			value := readStringList(iterator)
			object.externalCoordination = value
			object.bitmap_ |= 128
		case "incident_id":
			value := iterator.ReadString()
			object.incidentId = value
			object.bitmap_ |= 256
		case "incident_type":
			value := iterator.ReadString()
			object.incidentType = value
			object.bitmap_ |= 512
		case "last_updated":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdated = value
			object.bitmap_ |= 1024
		case "primary_team":
			value := iterator.ReadString()
			object.primaryTeam = value
			object.bitmap_ |= 2048
		case "severity":
			value := iterator.ReadString()
			object.severity = value
			object.bitmap_ |= 4096
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 8192
		case "summary":
			value := iterator.ReadString()
			object.summary = value
			object.bitmap_ |= 16384
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 32768
		case "worked_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.workedAt = value
			object.bitmap_ |= 65536
		default:
			iterator.ReadAny()
		}
	}
	return object
}
