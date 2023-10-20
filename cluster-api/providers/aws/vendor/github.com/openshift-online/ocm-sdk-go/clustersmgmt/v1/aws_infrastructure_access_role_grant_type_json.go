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

// MarshalAWSInfrastructureAccessRoleGrant writes a value of the 'AWS_infrastructure_access_role_grant' type to the given writer.
func MarshalAWSInfrastructureAccessRoleGrant(object *AWSInfrastructureAccessRoleGrant, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAWSInfrastructureAccessRoleGrant(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAWSInfrastructureAccessRoleGrant writes a value of the 'AWS_infrastructure_access_role_grant' type to the given stream.
func writeAWSInfrastructureAccessRoleGrant(object *AWSInfrastructureAccessRoleGrant, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AWSInfrastructureAccessRoleGrantLinkKind)
	} else {
		stream.WriteString(AWSInfrastructureAccessRoleGrantKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console_url")
		stream.WriteString(object.consoleURL)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.role != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role")
		writeAWSInfrastructureAccessRole(object.role, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state_description")
		stream.WriteString(object.stateDescription)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("user_arn")
		stream.WriteString(object.userARN)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSInfrastructureAccessRoleGrant reads a value of the 'AWS_infrastructure_access_role_grant' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSInfrastructureAccessRoleGrant(source interface{}) (object *AWSInfrastructureAccessRoleGrant, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAWSInfrastructureAccessRoleGrant(iterator)
	err = iterator.Error
	return
}

// readAWSInfrastructureAccessRoleGrant reads a value of the 'AWS_infrastructure_access_role_grant' type from the given iterator.
func readAWSInfrastructureAccessRoleGrant(iterator *jsoniter.Iterator) *AWSInfrastructureAccessRoleGrant {
	object := &AWSInfrastructureAccessRoleGrant{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSInfrastructureAccessRoleGrantLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "console_url":
			value := iterator.ReadString()
			object.consoleURL = value
			object.bitmap_ |= 8
		case "role":
			value := readAWSInfrastructureAccessRole(iterator)
			object.role = value
			object.bitmap_ |= 16
		case "state":
			text := iterator.ReadString()
			value := AWSInfrastructureAccessRoleGrantState(text)
			object.state = value
			object.bitmap_ |= 32
		case "state_description":
			value := iterator.ReadString()
			object.stateDescription = value
			object.bitmap_ |= 64
		case "user_arn":
			value := iterator.ReadString()
			object.userARN = value
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
