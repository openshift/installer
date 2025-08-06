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

// MarshalRootVolume writes a value of the 'root_volume' type to the given writer.
func MarshalRootVolume(object *RootVolume, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteRootVolume(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteRootVolume writes a value of the 'root_volume' type to the given stream.
func WriteRootVolume(object *RootVolume, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAWSVolume(object.aws, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCPVolume(object.gcp, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRootVolume reads a value of the 'root_volume' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRootVolume(source interface{}) (object *RootVolume, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadRootVolume(iterator)
	err = iterator.Error
	return
}

// ReadRootVolume reads a value of the 'root_volume' type from the given iterator.
func ReadRootVolume(iterator *jsoniter.Iterator) *RootVolume {
	object := &RootVolume{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "aws":
			value := ReadAWSVolume(iterator)
			object.aws = value
			object.bitmap_ |= 1
		case "gcp":
			value := ReadGCPVolume(iterator)
			object.gcp = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
