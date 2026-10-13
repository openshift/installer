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

// MarshalAWSSTSRole writes a value of the 'AWSSTS_role' type to the given writer.
func MarshalAWSSTSRole(object *AWSSTSRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSSTSRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSSTSRole writes a value of the 'AWSSTS_role' type to the given stream.
func WriteAWSSTSRole(object *AWSSTSRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hcpManagedPolicies")
		stream.WriteBool(object.hcpManagedPolicies)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("isAdmin")
		stream.WriteBool(object.isAdmin)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managedPolicies")
		stream.WriteBool(object.managedPolicies)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("arn")
		stream.WriteString(object.roleARN)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.roleType)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("roleVersion")
		stream.WriteString(object.roleVersion)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSSTSRole reads a value of the 'AWSSTS_role' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSSTSRole(source interface{}) (object *AWSSTSRole, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSSTSRole(iterator)
	err = iterator.Error
	return
}

// ReadAWSSTSRole reads a value of the 'AWSSTS_role' type from the given iterator.
func ReadAWSSTSRole(iterator *jsoniter.Iterator) *AWSSTSRole {
	object := &AWSSTSRole{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "hcpManagedPolicies":
			value := iterator.ReadBool()
			object.hcpManagedPolicies = value
			object.fieldSet_[0] = true
		case "isAdmin":
			value := iterator.ReadBool()
			object.isAdmin = value
			object.fieldSet_[1] = true
		case "managedPolicies":
			value := iterator.ReadBool()
			object.managedPolicies = value
			object.fieldSet_[2] = true
		case "arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.fieldSet_[3] = true
		case "type":
			value := iterator.ReadString()
			object.roleType = value
			object.fieldSet_[4] = true
		case "roleVersion":
			value := iterator.ReadString()
			object.roleVersion = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
