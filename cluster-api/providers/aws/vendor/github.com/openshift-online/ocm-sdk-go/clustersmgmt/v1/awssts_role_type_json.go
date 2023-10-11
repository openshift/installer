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

// MarshalAWSSTSRole writes a value of the 'AWSSTS_role' type to the given writer.
func MarshalAWSSTSRole(object *AWSSTSRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAWSSTSRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAWSSTSRole writes a value of the 'AWSSTS_role' type to the given stream.
func writeAWSSTSRole(object *AWSSTSRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hcpManagedPolicies")
		stream.WriteBool(object.hcpManagedPolicies)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("isAdmin")
		stream.WriteBool(object.isAdmin)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managedPolicies")
		stream.WriteBool(object.managedPolicies)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("arn")
		stream.WriteString(object.roleARN)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.roleType)
		count++
	}
	present_ = object.bitmap_&32 != 0
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
	object = readAWSSTSRole(iterator)
	err = iterator.Error
	return
}

// readAWSSTSRole reads a value of the 'AWSSTS_role' type from the given iterator.
func readAWSSTSRole(iterator *jsoniter.Iterator) *AWSSTSRole {
	object := &AWSSTSRole{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "hcpManagedPolicies":
			value := iterator.ReadBool()
			object.hcpManagedPolicies = value
			object.bitmap_ |= 1
		case "isAdmin":
			value := iterator.ReadBool()
			object.isAdmin = value
			object.bitmap_ |= 2
		case "managedPolicies":
			value := iterator.ReadBool()
			object.managedPolicies = value
			object.bitmap_ |= 4
		case "arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.bitmap_ |= 8
		case "type":
			value := iterator.ReadString()
			object.roleType = value
			object.bitmap_ |= 16
		case "roleVersion":
			value := iterator.ReadString()
			object.roleVersion = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
