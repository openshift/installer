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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAWSSTSPolicy writes a value of the 'AWSSTS_policy' type to the given writer.
func MarshalAWSSTSPolicy(object *AWSSTSPolicy, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSSTSPolicy(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSSTSPolicy writes a value of the 'AWSSTS_policy' type to the given stream.
func WriteAWSSTSPolicy(object *AWSSTSPolicy, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("arn")
		stream.WriteString(object.arn)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("details")
		stream.WriteString(object.details)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSSTSPolicy reads a value of the 'AWSSTS_policy' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSSTSPolicy(source interface{}) (object *AWSSTSPolicy, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSSTSPolicy(iterator)
	err = iterator.Error
	return
}

// ReadAWSSTSPolicy reads a value of the 'AWSSTS_policy' type from the given iterator.
func ReadAWSSTSPolicy(iterator *jsoniter.Iterator) *AWSSTSPolicy {
	object := &AWSSTSPolicy{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "arn":
			value := iterator.ReadString()
			object.arn = value
			object.bitmap_ |= 1
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 2
		case "details":
			value := iterator.ReadString()
			object.details = value
			object.bitmap_ |= 4
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
