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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAWSVolume writes a value of the 'AWS_volume' type to the given writer.
func MarshalAWSVolume(object *AWSVolume, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSVolume(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSVolume writes a value of the 'AWS_volume' type to the given stream.
func WriteAWSVolume(object *AWSVolume, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("iops")
		stream.WriteInt(object.iops)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("size")
		stream.WriteInt(object.size)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSVolume reads a value of the 'AWS_volume' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSVolume(source interface{}) (object *AWSVolume, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSVolume(iterator)
	err = iterator.Error
	return
}

// ReadAWSVolume reads a value of the 'AWS_volume' type from the given iterator.
func ReadAWSVolume(iterator *jsoniter.Iterator) *AWSVolume {
	object := &AWSVolume{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "iops":
			value := iterator.ReadInt()
			object.iops = value
			object.bitmap_ |= 1
		case "size":
			value := iterator.ReadInt()
			object.size = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
