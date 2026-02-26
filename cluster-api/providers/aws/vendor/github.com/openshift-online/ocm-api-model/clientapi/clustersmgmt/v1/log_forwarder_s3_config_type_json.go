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

// MarshalLogForwarderS3Config writes a value of the 'log_forwarder_S3_config' type to the given writer.
func MarshalLogForwarderS3Config(object *LogForwarderS3Config, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteLogForwarderS3Config(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteLogForwarderS3Config writes a value of the 'log_forwarder_S3_config' type to the given stream.
func WriteLogForwarderS3Config(object *LogForwarderS3Config, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bucket_name")
		stream.WriteString(object.bucketName)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("bucket_prefix")
		stream.WriteString(object.bucketPrefix)
	}
	stream.WriteObjectEnd()
}

// UnmarshalLogForwarderS3Config reads a value of the 'log_forwarder_S3_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalLogForwarderS3Config(source interface{}) (object *LogForwarderS3Config, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadLogForwarderS3Config(iterator)
	err = iterator.Error
	return
}

// ReadLogForwarderS3Config reads a value of the 'log_forwarder_S3_config' type from the given iterator.
func ReadLogForwarderS3Config(iterator *jsoniter.Iterator) *LogForwarderS3Config {
	object := &LogForwarderS3Config{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "bucket_name":
			value := iterator.ReadString()
			object.bucketName = value
			object.fieldSet_[0] = true
		case "bucket_prefix":
			value := iterator.ReadString()
			object.bucketPrefix = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
