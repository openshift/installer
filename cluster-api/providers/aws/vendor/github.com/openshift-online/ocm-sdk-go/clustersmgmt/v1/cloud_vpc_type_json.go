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

// MarshalCloudVPC writes a value of the 'cloud_VPC' type to the given writer.
func MarshalCloudVPC(object *CloudVPC, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeCloudVPC(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeCloudVPC writes a value of the 'cloud_VPC' type to the given stream.
func writeCloudVPC(object *CloudVPC, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.awsSubnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_subnets")
		writeSubnetworkList(object.awsSubnets, stream)
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
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.subnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnets")
		writeStringList(object.subnets, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalCloudVPC reads a value of the 'cloud_VPC' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCloudVPC(source interface{}) (object *CloudVPC, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readCloudVPC(iterator)
	err = iterator.Error
	return
}

// readCloudVPC reads a value of the 'cloud_VPC' type from the given iterator.
func readCloudVPC(iterator *jsoniter.Iterator) *CloudVPC {
	object := &CloudVPC{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "aws_subnets":
			value := readSubnetworkList(iterator)
			object.awsSubnets = value
			object.bitmap_ |= 1
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 2
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 4
		case "subnets":
			value := readStringList(iterator)
			object.subnets = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
