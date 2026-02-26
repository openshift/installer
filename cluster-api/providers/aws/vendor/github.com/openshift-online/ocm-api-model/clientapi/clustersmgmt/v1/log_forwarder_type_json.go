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

// MarshalLogForwarder writes a value of the 'log_forwarder' type to the given writer.
func MarshalLogForwarder(object *LogForwarder, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarder(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarder writes a value of the 'log_forwarder' type to the given stream.
func WriteLogForwarder(object *LogForwarder, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(LogForwarderLinkKind)
	} else {
		stream.WriteString(LogForwarderKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.s3 != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("s3")
		WriteLogForwarderS3Config(object.s3, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.applications != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("applications")
		WriteStringList(object.applications, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.cloudwatch != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloudwatch")
		WriteLogForwarderCloudWatchConfig(object.cloudwatch, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		WriteLogForwarderGroupList(object.groups, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteLogForwarderStatus(object.status, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarder reads a value of the 'log_forwarder' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarder(source interface{}) (object *LogForwarder, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarder(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarder reads a value of the 'log_forwarder' type from the given iterator.
func ReadLogForwarder(iterator *jsoniter.Iterator) *LogForwarder {
	object := &LogForwarder{
		fieldSet_: make([]bool, 9),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == LogForwarderLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "s3":
			value := ReadLogForwarderS3Config(iterator)
			object.s3 = value
			object.fieldSet_[3] = true
		case "applications":
			value := ReadStringList(iterator)
			object.applications = value
			object.fieldSet_[4] = true
		case "cloudwatch":
			value := ReadLogForwarderCloudWatchConfig(iterator)
			object.cloudwatch = value
			object.fieldSet_[5] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[6] = true
		case "groups":
			value := ReadLogForwarderGroupList(iterator)
			object.groups = value
			object.fieldSet_[7] = true
		case "status":
			value := ReadLogForwarderStatus(iterator)
			object.status = value
			object.fieldSet_[8] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
