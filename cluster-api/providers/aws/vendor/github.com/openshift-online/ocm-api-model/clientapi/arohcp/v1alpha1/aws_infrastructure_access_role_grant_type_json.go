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

// MarshalAWSInfrastructureAccessRoleGrant writes a value of the 'AWS_infrastructure_access_role_grant' type to the given writer.
func MarshalAWSInfrastructureAccessRoleGrant(object *AWSInfrastructureAccessRoleGrant, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSInfrastructureAccessRoleGrant(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSInfrastructureAccessRoleGrant writes a value of the 'AWS_infrastructure_access_role_grant' type to the given stream.
func WriteAWSInfrastructureAccessRoleGrant(object *AWSInfrastructureAccessRoleGrant, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AWSInfrastructureAccessRoleGrantLinkKind)
	} else {
		stream.WriteString(AWSInfrastructureAccessRoleGrantKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console_url")
		stream.WriteString(object.consoleURL)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.role != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role")
		WriteAWSInfrastructureAccessRole(object.role, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state_description")
		stream.WriteString(object.stateDescription)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
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
	object = ReadAWSInfrastructureAccessRoleGrant(iterator)
	err = iterator.Error
	return
}

// ReadAWSInfrastructureAccessRoleGrant reads a value of the 'AWS_infrastructure_access_role_grant' type from the given iterator.
func ReadAWSInfrastructureAccessRoleGrant(iterator *jsoniter.Iterator) *AWSInfrastructureAccessRoleGrant {
	object := &AWSInfrastructureAccessRoleGrant{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSInfrastructureAccessRoleGrantLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "console_url":
			value := iterator.ReadString()
			object.consoleURL = value
			object.fieldSet_[3] = true
		case "role":
			value := ReadAWSInfrastructureAccessRole(iterator)
			object.role = value
			object.fieldSet_[4] = true
		case "state":
			text := iterator.ReadString()
			value := AWSInfrastructureAccessRoleGrantState(text)
			object.state = value
			object.fieldSet_[5] = true
		case "state_description":
			value := iterator.ReadString()
			object.stateDescription = value
			object.fieldSet_[6] = true
		case "user_arn":
			value := iterator.ReadString()
			object.userARN = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
