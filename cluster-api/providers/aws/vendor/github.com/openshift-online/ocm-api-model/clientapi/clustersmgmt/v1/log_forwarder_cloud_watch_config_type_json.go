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

// MarshalLogForwarderCloudWatchConfig writes a value of the 'log_forwarder_cloud_watch_config' type to the given writer.
func MarshalLogForwarderCloudWatchConfig(object *LogForwarderCloudWatchConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderCloudWatchConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderCloudWatchConfig writes a value of the 'log_forwarder_cloud_watch_config' type to the given stream.
func WriteLogForwarderCloudWatchConfig(object *LogForwarderCloudWatchConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_distribution_role_arn")
		stream.WriteString(object.logDistributionRoleArn)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_group_name")
		stream.WriteString(object.logGroupName)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarderCloudWatchConfig reads a value of the 'log_forwarder_cloud_watch_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarderCloudWatchConfig(source interface{}) (object *LogForwarderCloudWatchConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarderCloudWatchConfig(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderCloudWatchConfig reads a value of the 'log_forwarder_cloud_watch_config' type from the given iterator.
func ReadLogForwarderCloudWatchConfig(iterator *jsoniter.Iterator) *LogForwarderCloudWatchConfig {
	object := &LogForwarderCloudWatchConfig{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "log_distribution_role_arn":
			value := iterator.ReadString()
			object.logDistributionRoleArn = value
			object.fieldSet_[0] = true
		case "log_group_name":
			value := iterator.ReadString()
			object.logGroupName = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
