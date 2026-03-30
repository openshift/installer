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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalCloudVPC writes a value of the 'cloud_VPC' type to the given writer.
func MarshalCloudVPC(object *CloudVPC, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCloudVPC(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCloudVPC writes a value of the 'cloud_VPC' type to the given stream.
func WriteCloudVPC(object *CloudVPC, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.awsSecurityGroups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_security_groups")
		WriteSecurityGroupList(object.awsSecurityGroups, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.awsSubnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_subnets")
		WriteSubnetworkList(object.awsSubnets, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cidr_block")
		stream.WriteString(object.cidrBlock)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("red_hat_managed")
		stream.WriteBool(object.redHatManaged)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.subnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnets")
		WriteStringList(object.subnets, stream)
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
	object = ReadCloudVPC(iterator)
	err = iterator.Error
	return
}

// ReadCloudVPC reads a value of the 'cloud_VPC' type from the given iterator.
func ReadCloudVPC(iterator *jsoniter.Iterator) *CloudVPC {
	object := &CloudVPC{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "aws_security_groups":
			value := ReadSecurityGroupList(iterator)
			object.awsSecurityGroups = value
			object.fieldSet_[0] = true
		case "aws_subnets":
			value := ReadSubnetworkList(iterator)
			object.awsSubnets = value
			object.fieldSet_[1] = true
		case "cidr_block":
			value := iterator.ReadString()
			object.cidrBlock = value
			object.fieldSet_[2] = true
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[3] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[4] = true
		case "red_hat_managed":
			value := iterator.ReadBool()
			object.redHatManaged = value
			object.fieldSet_[5] = true
		case "subnets":
			value := ReadStringList(iterator)
			object.subnets = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
